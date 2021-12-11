package handlers_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"os"
	"time"

	"github.com/go-redis/redismock/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/jamesyang124/ethereum-auth/handlers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe(".\\Nonce", func() {
	redisTTL := "10"
	l := log.New(os.Stdout, "15:04:05 | ", 0)
	// discard app log
	l.SetOutput(ioutil.Discard)

	var ctx = context.TODO()
	db, mock := redismock.NewClientMock()

	app := fiber.New()
	app.Post("/auth/nonce", handlers.NonceHandler(ctx, db, l, redisTTL))

	It("should respond 200 for nonce api with input payload", func() {
		payload := `{"cid": "cid", "nid": "nid", "paddr": "paddr"}`

		req := httptest.NewRequest("POST", "/auth/nonce", bytes.NewBuffer([]byte(payload)))
		req.Header.Set("Content-Type", "application/json")

		duration, _ := time.ParseDuration(redisTTL + "s")
		mock.ClearExpect()
		mock.ExpectGet("cid-nid-paddr").RedisNil()
		mock.Regexp().ExpectSetEX("cid-nid-paddr", `\d{6}`, duration).SetVal("OK")

		resp, _ := app.Test(req)
		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		Expect(mock.ExpectationsWereMet()).Should(Succeed())
		Expect(resp.StatusCode).To(Equal(200))
		Expect(string(bodyBytes)).Should(MatchRegexp(`\d{6}`))
	})

	It("should respond 400 for nonce api without input payload", func() {
		req := httptest.NewRequest("POST", "/auth/nonce", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		Expect(resp.StatusCode).To(Equal(400))
		Expect(string(bodyBytes)).Should(MatchRegexp(`parsing nonce reuqest input failed invalid character.*`))
	})

	It("should respond 500 for nonce api if redis client server error", func() {
		payload := `{"cid": "cid", "nid": "nid", "paddr": "paddr"}`

		req := httptest.NewRequest("POST", "/auth/nonce", bytes.NewBuffer([]byte(payload)))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		Expect(resp.StatusCode).To(Equal(500))
		Expect(string(bodyBytes)).To(Equal("Internal Server Error"))
	})
})
