package handlers

import (
	"bytes"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler(t *testing.T) {
	app := fiber.New()
	app.Get("/health", HealthCheckHandler)

	tests := []struct {
		name         string
		payload      interface{}
		response     interface{}
		expectedCode int
	}{
		{name: "get-health-check-api-respond-200", payload: nil, response: "OK", expectedCode: 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var payload bytes.Buffer

			if tt.payload != nil {
				payload = *bytes.NewBuffer([]byte(tt.payload.(string)))
			}

			req := httptest.NewRequest("GET", "/health", &payload)

			resp, _ := app.Test(req)
			bodyBytes, _ := ioutil.ReadAll(resp.Body)

			assert.Equalf(t, tt.response, string(bodyBytes), tt.name)
			assert.Equalf(t, tt.expectedCode, resp.StatusCode, tt.name)
		})
	}
}
