package handlers_test

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gofiber/fiber/v2"
	"viveportengineering/DoubleA/ethereum-auth/handlers"
)

var _ = Describe(".\\Metadata", func() {
	godotenv.Load("../.env")
	app := fiber.New()
	templateText := `I am signin with this %s`
	redisTTL := os.Getenv("REDIS_CACHE_TTL_SECONDS")
	ttl, err := strconv.Atoi(redisTTL)

	app.Get("/api/ethereum-auth/v1/metadata", handlers.MetadataHandler(templateText, ttl))

	Context("", func() {
		BeforeEach(func() {
			Expect(err).Should(BeNil())
		})

		It("should respond 200 when without input payload", func() {
			resp, _ := app.Test(httptest.NewRequest("GET", "/api/ethereum-auth/v1/metadata", nil))
			bodyBytes, _ := ioutil.ReadAll(resp.Body)

			var result map[string]interface{}
			json.Unmarshal(bodyBytes, &result)

			Expect(resp.StatusCode).To(Equal(200))
			Expect(result).To(Not(BeNil()))
			Expect(result["signin-text-template"]).To(Equal(templateText))
			Expect(int(result["ttl-seconds"].(float64))).To(Equal(ttl))
		})
	})
})
