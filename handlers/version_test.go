package handlers_test

import (
	"io/ioutil"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gofiber/fiber/v2"
	"viveportengineering/DoubleA/ethereum-auth/handlers"
)

var _ = Describe(".\\Version", func() {
	app := fiber.New()
	appVersion := "experiment"
	app.Get("/version", handlers.VersionHandler(appVersion))

	It("should respond 200 for version api without input payload", func() {
		resp, _ := app.Test(httptest.NewRequest("GET", "/version", nil))
		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		Expect(resp.StatusCode).To(Equal(200))
		Expect(string(bodyBytes)).To(Equal(appVersion))
	})
})
