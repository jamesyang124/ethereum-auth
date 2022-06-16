package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gofiber/fiber/v2"

	"github.com/go-redis/redis/v8"

	"viveportengineering/DoubleA/ethereum-auth/errors"
	"viveportengineering/DoubleA/ethereum-auth/utils"

	"log"
	"time"
)

const SIGNATURE_RI_MAGIC_NUM = 27

type LoginRequest struct {
	Extra         map[string]interface{} `json:"extra"`
	Signature     string                 `json:"sig"`
	PublicAddress string                 `json:"paddr"`
}

type DownstreamAuthRequest struct {
	Extra         map[string]interface{} `json:"extra"`
	PublicAddress string                 `json:"paddr"`
}

func responseErrorLogging(code int, errs []error, l *log.Logger) {
	l.Printf("unexpected error with status code - %d\n", code)

	for i := 0; i < len(errs); i++ {
		l.Print("error msg: ", errs[i].Error())
	}
}

func downstreamAuthRequest(url string, extra map[string]interface{}, paddr string) (int, []byte, []error) {
	agent := fiber.Post(url)
	agent.JSON(DownstreamAuthRequest{
		Extra:         extra,
		PublicAddress: paddr,
	})

	return agent.Bytes()
}

// @Summary      check signed message and authenticate user
// @Tags         login
// @Param        extra  body  interface{}  false  "auth info for downstream auth system could carry by this field, as json format"
// @Param        paddr  body  string  true  "ethereum digital wallet public address"
// @Accept       json
// @Produce      json
// @Success      200  {object} interface{} "proxy downstream authenticated response json"
// @Router       /api/ethereum-auth/v1/login [post]
func LoginHandler(ctx context.Context, rdb *redis.Client,
	l *log.Logger, redisTTL string,
	signInTextTemplate string, downstreamAuthUri string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		lr := new(LoginRequest)

		if err := c.BodyParser(lr); err != nil {
			errorText := fmt.Sprintf("parsing login reuqest input failed %s\n", err.Error())
			l.Print(errorText)
			return fiber.NewError(fiber.StatusBadRequest, errorText)
		}
		l.Printf("%+v\n", lr)

		var b strings.Builder
		b.Reset()
		b.WriteString("ethereum-auth-")
		b.WriteString(lr.PublicAddress)
		key := b.String()
		l.Printf("redis key %s", key)

		// fetch public address bound nonce, if address not existed respond error
		// so user should follow /api/ethereum-auth/v1/nonce api as first step
		nonce, err := rdb.Get(ctx, key).Result()
		if err == redis.Nil {
			l.Printf("not found key %s from redis  - %s", key, err.Error())
			return errors.NONCE_EXPIRED_ERROR
		}

		if err != nil {
			l.Printf("get key %s from redis failed - %s", key, err.Error())
			return errors.INTERNAL_SERVER_ERROR
		}

		l.Printf("%s nonce is %s\n", lr.PublicAddress, nonce)

		// compose raw message
		raw_message := fmt.Sprintf(signInTextTemplate, nonce)

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
			l.Printf("input signature decode error: %s\n", err.Error())
			return errors.INVALID_SIGNATURE_ERROR
		}

		// 32 32 1
		//  r  s v
		// reverse recovery identifier v to 0 or 1 (signed with 27 or 28)
		if signature[64] >= SIGNATURE_RI_MAGIC_NUM { // https://htcsense.jira.com/browse/OPS-31163
			signature[64] -= SIGNATURE_RI_MAGIC_NUM
		}
		sigPublicKeyBytes, err := crypto.Ecrecover(hash.Bytes(), signature)
		if err != nil {
			l.Printf("crypto recovering error: %s\n", err.Error())
			return errors.INVALID_SIGNATURE_ERROR
		}

		ecdsaPublicKey, err := crypto.UnmarshalPubkey(sigPublicKeyBytes)
		if err != nil {
			l.Printf("crypto recovering error: %s\n", err.Error())
			return errors.INVALID_SIGNATURE_ERROR
		}
		l.Printf("secp256k1 public key - %s\n", ecdsaPublicKey)

		recoveredAddress := crypto.PubkeyToAddress(*ecdsaPublicKey).Hex()
		l.Printf("recovered public address - %s\n", recoveredAddress)

		// verify public adddress as authentication
		if !strings.EqualFold(lr.PublicAddress, recoveredAddress) {
			l.Printf("Verify signature and paddr failed, input: [%s], recovered: [%s]",
				lr.PublicAddress,
				recoveredAddress,
			)
			return errors.INVALID_SIGNATURE_ERROR
		}

		// update nonce to avoid replay attack
		nonce = strconv.Itoa(utils.RandNonce())
		duration, _ := time.ParseDuration(redisTTL + "s")
		err = rdb.SetEX(ctx, key, nonce, duration).Err()

		if err != nil {
			l.Printf(
				"authenticated but fail to update key %s and nonce %s from redis - %s",
				key,
				nonce,
				err.Error(),
			)
			return errors.INTERNAL_SERVER_ERROR
		}

		// bind downstream auth system
		code, body, errs := downstreamAuthRequest(downstreamAuthUri, lr.Extra, lr.PublicAddress)

		if errs != nil {
			l.Printf("unexpected http client request error for downstream uri [%s]", downstreamAuthUri)
			responseErrorLogging(code, errs, l)
			return errors.INTERNAL_SERVER_ERROR
		}

		var resp interface{}
		err = json.Unmarshal(body, &resp)
		if err != nil {
			l.Printf(
				"downstream uri [%s] respond code [%d] body [%s] but as invalid JSON response [%v], error: %s",
				downstreamAuthUri,
				code,
				body,
				resp,
				err.Error(),
			)
			return errors.SERVICE_FAILED_DEPENDENCY_ERROR
		}

		return c.Type("application/json").Status(code).JSON(resp)
	}
}
