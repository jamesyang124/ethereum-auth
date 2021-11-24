package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gofiber/fiber/v2"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"log"
	"time"
)

const SIGNATURE_RI_MAGIC_NUM = 27

type NonceRequest struct {
	ChainId       string `json:"cid"`
	NetworkId     string `json:"nid"`
	PublicAddress string `json:"paddr"`
}

type LoginRequest struct {
	Meta          map[string]interface{} `json:"meta"`
	ChainId       string                 `json:"cid"`
	NetworkId     string                 `json:"nid"`
	Signature     string                 `json:"sig"`
	PublicAddress string                 `json:"paddr"`
}

type DownstreamAuthRequest struct {
	Meta          map[string]interface{} `json:"meta"`
	ChainId       string                 `json:"cid"`
	NetworkId     string                 `json:"nid"`
	PublicAddress string                 `json:"paddr"`
}

func randNonce() int {
	return 100000 + rand.Intn(int(time.Now().UnixNano()%1000000))
}

func downstreamAuthRequest(url string, meta map[string]interface{}, cid string, nid string, paddr string) (int, []byte, []error) {
	agent := fiber.Post(url)
	agent.JSON(DownstreamAuthRequest{
		Meta:          meta,
		ChainId:       cid,
		NetworkId:     nid,
		PublicAddress: paddr,
	})

	return agent.Bytes()
}

func responseErrorLogging(code int, errs []error, l *log.Logger) {
	l.Printf("unexpected error with status code - %d\n", code)

	for i := 0; i < len(errs); i++ {
		l.Printf("error msg: ", errs[i].Error())
	}
}

func loadNonEmptyEnv(key string, l *log.Logger) string {
	v := os.Getenv(key)
	if v == "" {
		l.Fatal("missing " + key + " env")
	} else {
		l.Printf("%s=%s", key, v)
	}

	return v
}

func main() {
	// setup logger
	l := log.New(os.Stdout, "15:04:05 | ", 0)

	// load envs
	runTimeEnv := os.Getenv("RUNTIME_ENV")
	if runTimeEnv == "local" {
		godotenv.Load(".env")
	}

	downstreamAuthUri := loadNonEmptyEnv("DOWNSTREAM_AUTH_URI", l)
	redisHost := loadNonEmptyEnv("REDIS_CACHE_HOST", l)
	redisPort := loadNonEmptyEnv("REDIS_CACHE_PORT", l)
	redisTTL := loadNonEmptyEnv("REDIS_CACHE_TTL_SECONDS", l)

	// setup redis client
	var ctx = context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: "",
	})

	// init and inject middleware before route/handler registration
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	// TODO: add version, health check api

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	// respond nonce and persist to redis with TTL
	app.Post("/auth/nonce", func(c *fiber.Ctx) error {
		ar := new(NonceRequest)

		if err := c.BodyParser(ar); err != nil {
			l.Printf("parsing nonce reuqest input failed %s\n", err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		l.Printf("%+v\n", ar)

		var b strings.Builder
		b.Reset()
		b.WriteString(ar.ChainId)
		b.WriteString("-")
		b.WriteString(ar.NetworkId)
		b.WriteString("-")
		b.WriteString(ar.PublicAddress)
		key := b.String()
		l.Printf("redis key %s", key)

		// cache it to mitigate redis w+ ops
		nonce, err := rdb.Get(ctx, key).Result()
		if err == redis.Nil {
			// alway create new nonce and upate to its redis key
			nonce = strconv.Itoa(randNonce())
			duration, _ := time.ParseDuration(redisTTL + "s")
			err = rdb.SetEX(ctx, key, nonce, duration).Err()
		}

		if err != nil {
			l.Printf("get/set cid-nid-paddr key from redis failed - %s", err.Error())
			return fiber.NewError(fiber.StatusInternalServerError)
		}

		l.Printf("key %s nonce is %s\n", key, nonce)
		return c.SendString(nonce)
	})

	// check signed message and authenticate user
	app.Post("/auth/login", func(c *fiber.Ctx) error {
		lr := new(LoginRequest)

		if err := c.BodyParser(lr); err != nil {
			l.Printf("parsing login reuqest input failed - %s\n", err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		l.Printf("%+v\n", lr)

		var b strings.Builder
		b.Reset()
		b.WriteString(lr.ChainId)
		b.WriteString("-")
		b.WriteString(lr.NetworkId)
		b.WriteString("-")
		b.WriteString(lr.PublicAddress)
		key := b.String()
		l.Printf("redis key %s", key)

		// fetch public address bound nonce, if address not existed respond error
		// so user should follow /auth/nonce api as first step
		nonce, err := rdb.Get(ctx, key).Result()
		if err == redis.Nil {
			l.Printf("not found key %s from redis  - %s", key, err.Error())
			return fiber.NewError(fiber.StatusBadRequest, "nonce already expired")
		}

		if err != nil {
			l.Printf("get key %s from redis failed - %s", key, err.Error())
			return fiber.NewError(fiber.StatusBadRequest, "nonce already expired")
		}

		l.Printf("%s nonce is %s\n", lr.PublicAddress, nonce)

		// compose raw message
		b.Reset()
		b.WriteString("I am signing with this one-time 6-digit nonce: ")
		b.WriteString(nonce)
		raw_message := b.String()

		b.Reset()
		b.WriteString("\x19Ethereum Signed Message:\n")
		b.WriteString(strconv.Itoa(len(raw_message)))
		b.WriteString(raw_message)
		msg := b.String()

		l.Printf("msg: %s\n", msg)

		// raw message to keccak 256 hashed hex string
		data := []byte(msg)
		hash := crypto.Keccak256Hash(data)

		// input signature to []byte
		signature, err := hexutil.Decode(lr.Signature)
		if err != nil {
			l.Printf("input signature decode error:\n")
			l.Printf("%s\n", err.Error())

			return fiber.NewError(fiber.StatusBadRequest, "invalid input signature")
		}

		// 32 32 1
		//  r  s v
		// reverse recovery identifier v to 0 or 1 (signed with 27 or 28)
		signature[64] -= SIGNATURE_RI_MAGIC_NUM
		sigPublicKeyBytes, err := crypto.Ecrecover(hash.Bytes(), signature)
		if err != nil {
			l.Printf("crypto recovering error:\n")
			log.Fatal(err)
		}

		ecdsaPublicKey, err := crypto.UnmarshalPubkey(sigPublicKeyBytes)
		if err != nil {
			l.Printf("crypto recovering error:\n")
			log.Fatal(err)
		}
		l.Printf("secp256k1 public key - %s\n", ecdsaPublicKey)

		recoveredAddress := crypto.PubkeyToAddress(*ecdsaPublicKey).Hex()
		l.Printf("recovered public address - %s\n", recoveredAddress)

		// verify public adddress as authentication
		if strings.ToLower(lr.PublicAddress) != strings.ToLower(recoveredAddress) {
			return fiber.NewError(fiber.StatusBadRequest, "Verify signature failed")
		}

		// update nonce to avoid replay attack
		nonce = strconv.Itoa(randNonce())
		duration, _ := time.ParseDuration(redisTTL + "s")
		err = rdb.SetEX(ctx, key, nonce, duration).Err()

		if err != nil {
			l.Printf("authenticated but fail to update key %s and nonce %s from redis - %s", key, nonce, err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, "Request Service failed, please try again later")
		}

		// bind downstream auth system
		code, body, errs := downstreamAuthRequest(downstreamAuthUri, lr.Meta, lr.ChainId, lr.NetworkId, lr.PublicAddress)

		if errs != nil {
			responseErrorLogging(code, errs, l)
			return fiber.NewError(fiber.StatusBadRequest, errs[0].Error())
		}

		// we assume downstreamAuthRequest should respond json in any case
		var resp interface{}
		json.Unmarshal(body, &resp)
		return c.Type("application/json").Status(code).JSON(resp.(map[string]interface{}))
	})

	app.Listen(":3030")
}
