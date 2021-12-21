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
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe(".\\Login", func() {
	redisTTL := "10"
	template := "I am sign-in with this one-time 6-digit nonce: %s"
	downstreamAuthUri := "https://csdev.htcwowdev.com/SS/api/social-authentication/v1/ethereum/auth"
	l := log.New(os.Stdout, "15:04:05 | ", 0)
	// discard app log
	l.SetOutput(ioutil.Discard)

	var ctx = context.TODO()
	db, mock := redismock.NewClientMock()

	app := fiber.New()
	app.Post("/auth/login", handlers.LoginHandler(ctx, db, l, redisTTL, template, downstreamAuthUri))

	It("should respond 200 for login api with input payload", func() {
		// TODO: extra and sig mock
		payload := `{"cid": "cid", "nid": "nid", "paddr": "paddr", "sig": "sig", "extra": "extra"}`

		req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer([]byte(payload)))
		req.Header.Set("Content-Type", "application/json")

		duration, _ := time.ParseDuration(redisTTL + "s")
		mock.ClearExpect()
		mock.ExpectGet("cid-nid-paddr").RedisNil()
		mock.Regexp().ExpectSetEX("cid-nid-paddr", `\d{6}`, duration).SetVal("OK")

		fixture := `{"status":{"message": "Your message", "code": 200}}`
		responder := httpmock.NewStringResponder(200, fixture)
		httpmock.RegisterResponder("POST", downstreamAuthUri, responder)

		resp, _ := app.Test(req)
		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		Expect(mock.ExpectationsWereMet()).Should(Succeed())
		Expect(resp.StatusCode).To(Equal(200))
		Expect(string(bodyBytes)).Should(MatchRegexp(`\d{6}`))
	})

	It("should respond 400 for login api without input payload", func() {
		req := httptest.NewRequest("POST", "/auth/login", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		Expect(resp.StatusCode).To(Equal(400))
		Expect(string(bodyBytes)).Should(MatchRegexp(`parsing login reuqest input failed invalid character.*`))
	})
})
