package handlers_test

import (
	"bytes"
	"io/ioutil"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gofiber/fiber/v2"
	"github.com/jamesyang124/ethereum-auth/handlers"
)

var _ = Describe(".\\HealthCheck", func() {
	app := fiber.New()
	app.Get("/api/ethereum-auth/health", handlers.HealthCheckHandler)

	It("should respond 200 for health api without input payload", func() {
		resp, _ := app.Test(httptest.NewRequest("GET", "/api/ethereum-auth/health", nil))
		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		Expect(resp.StatusCode).To(Equal(200))
		Expect(string(bodyBytes)).To(Equal("OK"))
	})

	It("should respond 200 for health api with input payload", func() {
		payload := bytes.NewBuffer([]byte(`{"id": 1}`))

		resp, _ := app.Test(httptest.NewRequest("GET", "/api/ethereum-auth/health", payload))
		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		Expect(resp.StatusCode).To(Equal(200))
		Expect(string(bodyBytes)).To(Equal("OK"))
	})
})
