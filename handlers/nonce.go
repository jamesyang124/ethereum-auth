package handlers

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/jamesyang124/ethereum-auth/utils"
	"log"
	"strconv"
	"strings"
	"time"
)

type NonceRequest struct {
	ChainId       string `json:"cid"`
	NetworkId     string `json:"nid"`
	PublicAddress string `json:"paddr"`
}

// @Summary      generate nonce and cached with TTL for specific chain, network, and public address
// @Tags         nonce
// @Param        cid    body  string  true  "ethereum chain id"
// @Param        nid    body  string  true  "ethereum network id"
// @Param        paddr  body  string  true  "ethereum digital wallet public address"
// @Accept       json
// @Produce      text/html
// @Success      200  {string} string "6 digit random nonce ex: 123453"
// @Router       /auth/nonce [post]
func NonceHandler(ctx context.Context, rdb *redis.Client,
	l *log.Logger, redisTTL string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
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
			nonce = strconv.Itoa(utils.RandNonce())
			duration, _ := time.ParseDuration(redisTTL + "s")
			err = rdb.SetEX(ctx, key, nonce, duration).Err()
		}

		if err != nil {
			l.Printf("get/set cid-nid-paddr key from redis failed - %s", err.Error())
			return fiber.NewError(fiber.StatusInternalServerError)
		}

		l.Printf("key %s nonce is %s\n", key, nonce)
		return c.SendString(nonce)
	}
}
