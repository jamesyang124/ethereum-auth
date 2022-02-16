package handlers_test

import (
	"bytes"
	"context"
	sys_error "errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"os"
	"time"

	"github.com/go-redis/redismock/v8"
	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/walkerus/go-wiremock"
	"viveportengineering/DoubleA/ethereum-auth/errors"
	"viveportengineering/DoubleA/ethereum-auth/handlers"
)

var _ = Describe(".\\Login", func() {
	redisTTL := "10"
	template := "I am sign-in with this one-time 6-digit nonce: %s"
	downstreamAuthUri := "http://localhost:8080/auth"
	l := log.New(os.Stdout, "15:04:05 | ", 0)
	// discard app log
	l.SetOutput(ioutil.Discard)

	var ctx = context.TODO()
	db, mock := redismock.NewClientMock()

	app := fiber.New(fiber.Config{
		ErrorHandler: errors.ErrorResponseHandler,
	})
	app.Post("/api/ethereum-auth/v1/login", handlers.LoginHandler(ctx, db, l, redisTTL, template, downstreamAuthUri))

	It("should respond 200 for login api with input payload", func() {
		wiremockClient := wiremock.NewClient("http://0.0.0.0:8080")
		defer wiremockClient.Reset()

		wiremockClient.StubFor(wiremock.Post(wiremock.URLPathEqualTo("/auth")).
			WillReturn(
				`{"code": 200, "detail": "detail"}`,
				map[string]string{"Content-Type": "application/json", "Connection": "Close"},
				200,
			))

		sig := "0xd5557bce14b5b70d8af657f08abd4d757c7ecca1923820f08833c07d4a022a937a59151840b18bf4b44e9ded2457c1e8a2fa0f549b535c9668f68fdbce0edd151c"
		nonce := "197007"
		paddr := "0x77b8e619b9e0Fb95C6c57A9fCb46Bd3D993F5238"

		payload := fmt.Sprintf(
			`{"paddr": "%s", "sig": "%s", "extra": {}}`,
			paddr,
			sig,
		)

		req := httptest.NewRequest("POST", "/api/ethereum-auth/v1/login", bytes.NewBuffer([]byte(payload)))
		req.Header.Set("Content-Type", "application/json")

		duration, _ := time.ParseDuration(redisTTL + "s")
		redisKey := fmt.Sprintf("ethereum-auth-%s", paddr)
		mock.ClearExpect()
		mock.ExpectGet(redisKey).SetVal(nonce)
		mock.Regexp().ExpectSetEX(redisKey, `\d{6}`, duration).SetVal("OK")

		resp, _ := app.Test(req)
		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		Expect(mock.ExpectationsWereMet()).Should(Succeed())
		Expect(resp.StatusCode).To(Equal(200))
		Expect(string(bodyBytes)).To(Equal(`{"code":200,"detail":"detail"}`))
	})

	It("should respond 400 for login api without input payload", func() {
		req := httptest.NewRequest("POST", "/api/ethereum-auth/v1/login", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		Expect(resp.StatusCode).To(Equal(400))
		Expect(string(bodyBytes)).Should(MatchRegexp(`parsing login reuqest input failed invalid character.*`))
	})

	It("should respond 400 with code 4001801 when nonce is expired", func() {
		sig := "0xd5557bce14b5b70d8af657f08abd4d757c7ecca1923820f08833c07d4a022a937a59151840b18bf4b44e9ded2457c1e8a2fa0f549b535c9668f68fdbce0edd151c"
		paddr := "0x77b8e619b9e0Fb95C6c57A9fCb46Bd3D993F5238"

		payload := fmt.Sprintf(
			`{"paddr": "%s", "sig": "%s", "extra": {}}`,
			paddr,
			sig,
		)

		req := httptest.NewRequest("POST", "/api/ethereum-auth/v1/login", bytes.NewBuffer([]byte(payload)))
		req.Header.Set("Content-Type", "application/json")

		redisKey := fmt.Sprintf("ethereum-auth-%s", paddr)
		mock.ClearExpect()
		mock.ExpectGet(redisKey).RedisNil()

		resp, _ := app.Test(req)
		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		Expect(resp.StatusCode).To(Equal(400))
		Expect(string(bodyBytes)).To(Equal(`{"code":4001801,"message":"nonce is expired"}`))
	})

	It("should respond 400 with code 4001802 when input signature is invalid", func() {
		sig := "0xsome-invalid-sig"
		nonce := "197007"
		paddr := "0x77b8e619b9e0Fb95C6c57A9fCb46Bd3D993F5238"

		payload := fmt.Sprintf(
			`{"paddr": "%s", "sig": "%s", "extra": {}}`,
			paddr,
			sig,
		)

		req := httptest.NewRequest("POST", "/api/ethereum-auth/v1/login", bytes.NewBuffer([]byte(payload)))
		req.Header.Set("Content-Type", "application/json")

		duration, _ := time.ParseDuration(redisTTL + "s")
		redisKey := fmt.Sprintf("ethereum-auth-%s", paddr)
		mock.ClearExpect()
		mock.ExpectGet(redisKey).SetVal(nonce)
		mock.Regexp().ExpectSetEX(redisKey, `\d{6}`, duration).SetVal("OK")

		resp, _ := app.Test(req)
		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		Expect(resp.StatusCode).To(Equal(400))
		Expect(string(bodyBytes)).To(Equal(`{"code":4001802,"message":"invalid input signature"}`))
	})

	It("should respond 424 when downstream api respond 500 internal server error", func() {
		wiremockClient := wiremock.NewClient("http://0.0.0.0:8080")
		defer wiremockClient.Reset()

		wiremockClient.StubFor(wiremock.Post(wiremock.URLPathEqualTo("/auth")).
			WillReturn(
				`internal server error`,
				map[string]string{"Content-Type": "text/html", "Connection": "Close"},
				500,
			))

		sig := "0xd5557bce14b5b70d8af657f08abd4d757c7ecca1923820f08833c07d4a022a937a59151840b18bf4b44e9ded2457c1e8a2fa0f549b535c9668f68fdbce0edd151c"
		nonce := "197007"
		paddr := "0x77b8e619b9e0Fb95C6c57A9fCb46Bd3D993F5238"

		payload := fmt.Sprintf(
			`{"paddr": "%s", "sig": "%s", "extra": {}}`,
			paddr,
			sig,
		)

		req := httptest.NewRequest("POST", "/api/ethereum-auth/v1/login", bytes.NewBuffer([]byte(payload)))
		req.Header.Set("Content-Type", "application/json")

		duration, _ := time.ParseDuration(redisTTL + "s")
		redisKey := fmt.Sprintf("ethereum-auth-%s", paddr)
		mock.ClearExpect()
		mock.ExpectGet(redisKey).SetVal(nonce)
		mock.Regexp().ExpectSetEX(redisKey, `\d{6}`, duration).SetVal("OK")

		resp, _ := app.Test(req)
		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		Expect(mock.ExpectationsWereMet()).Should(Succeed())
		Expect(resp.StatusCode).To(Equal(424))
		Expect(string(bodyBytes)).To(Equal("failed dependency"))
	})

	It("should respond 424 when downstream api respond 424 failed dependency", func() {
		wiremockClient := wiremock.NewClient("http://0.0.0.0:8080")
		defer wiremockClient.Reset()

		wiremockClient.StubFor(wiremock.Post(wiremock.URLPathEqualTo("/auth")).
			WillReturn(
				`failed dependency`,
				map[string]string{"Content-Type": "text/html", "Connection": "Close"},
				424,
			))

		sig := "0xd5557bce14b5b70d8af657f08abd4d757c7ecca1923820f08833c07d4a022a937a59151840b18bf4b44e9ded2457c1e8a2fa0f549b535c9668f68fdbce0edd151c"
		nonce := "197007"
		paddr := "0x77b8e619b9e0Fb95C6c57A9fCb46Bd3D993F5238"

		payload := fmt.Sprintf(
			`{"paddr": "%s", "sig": "%s", "extra": {}}`,
			paddr,
			sig,
		)

		req := httptest.NewRequest("POST", "/api/ethereum-auth/v1/login", bytes.NewBuffer([]byte(payload)))
		req.Header.Set("Content-Type", "application/json")

		duration, _ := time.ParseDuration(redisTTL + "s")
		redisKey := fmt.Sprintf("ethereum-auth-%s", paddr)
		mock.ClearExpect()
		mock.ExpectGet(redisKey).SetVal(nonce)
		mock.Regexp().ExpectSetEX(redisKey, `\d{6}`, duration).SetVal("OK")

		resp, _ := app.Test(req)
		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		Expect(mock.ExpectationsWereMet()).Should(Succeed())
		Expect(resp.StatusCode).To(Equal(424))
		Expect(string(bodyBytes)).To(Equal("failed dependency"))
	})

	It("should respond 424 when downstream api respond invalid json response", func() {
		wiremockClient := wiremock.NewClient("http://0.0.0.0:8080")
		defer wiremockClient.Reset()

		wiremockClient.StubFor(wiremock.Post(wiremock.URLPathEqualTo("/auth")).
			WillReturn(
				`invalid json format`,
				map[string]string{"Content-Type": "text/html", "Connection": "Close"},
				200,
			))

		sig := "0xd5557bce14b5b70d8af657f08abd4d757c7ecca1923820f08833c07d4a022a937a59151840b18bf4b44e9ded2457c1e8a2fa0f549b535c9668f68fdbce0edd151c"
		nonce := "197007"
		paddr := "0x77b8e619b9e0Fb95C6c57A9fCb46Bd3D993F5238"

		payload := fmt.Sprintf(
			`{"paddr": "%s", "sig": "%s", "extra": {}}`,
			paddr,
			sig,
		)

		req := httptest.NewRequest("POST", "/api/ethereum-auth/v1/login", bytes.NewBuffer([]byte(payload)))
		req.Header.Set("Content-Type", "application/json")

		duration, _ := time.ParseDuration(redisTTL + "s")
		redisKey := fmt.Sprintf("ethereum-auth-%s", paddr)
		mock.ClearExpect()
		mock.ExpectGet(redisKey).SetVal(nonce)
		mock.Regexp().ExpectSetEX(redisKey, `\d{6}`, duration).SetVal("OK")

		resp, _ := app.Test(req)
		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		Expect(mock.ExpectationsWereMet()).Should(Succeed())
		Expect(resp.StatusCode).To(Equal(424))
		Expect(string(bodyBytes)).To(Equal("failed dependency"))
	})

	It("should respond 500 when redis get nonce failed", func() {
		sig := "0xd5557bce14b5b70d8af657f08abd4d757c7ecca1923820f08833c07d4a022a937a59151840b18bf4b44e9ded2457c1e8a2fa0f549b535c9668f68fdbce0edd151c"
		paddr := "0x77b8e619b9e0Fb95C6c57A9fCb46Bd3D993F5238"

		payload := fmt.Sprintf(
			`{"paddr": "%s", "sig": "%s", "extra": {}}`,
			paddr,
			sig,
		)

		req := httptest.NewRequest("POST", "/api/ethereum-auth/v1/login", bytes.NewBuffer([]byte(payload)))
		req.Header.Set("Content-Type", "application/json")

		redisKey := fmt.Sprintf("ethereum-auth-%s", paddr)
		mock.ClearExpect()
		mock.ExpectGet(redisKey).SetErr(sys_error.New("system level error"))

		resp, _ := app.Test(req)
		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		Expect(resp.StatusCode).To(Equal(500))
		Expect(string(bodyBytes)).To(Equal("internal server error"))
	})
})
