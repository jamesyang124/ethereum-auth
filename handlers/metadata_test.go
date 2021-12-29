package handlers_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gofiber/fiber/v2"
	"viveportengineering/DoubleA/ethereum-auth/handlers"
)

var _ = Describe(".\\Metadata", func() {
	app := fiber.New()
	templateText := `I am signin with this %s`
	app.Get("/api/ethereum-auth/v1/metadata", handlers.MetadataHandler(templateText))

	It("should respond 200 for metadata api without input payload", func() {
		resp, _ := app.Test(httptest.NewRequest("GET", "/api/ethereum-auth/v1/metadata", nil))
		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		var result map[string]interface{}
		json.Unmarshal(bodyBytes, &result)

		Expect(resp.StatusCode).To(Equal(200))
		Expect(result).To(Equal(map[string]interface{}{"signin-text-template": templateText}))
	})
})
