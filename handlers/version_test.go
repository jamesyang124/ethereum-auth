package handlers

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestVersionRoute(t *testing.T) {
	app := fiber.New()
	appVersion := "experiment"
	app.Get("/version", VersionHandler(appVersion))

	tests := []struct {
		name         string
		response     interface{}
		expectedCode int
	}{
		{name: "get-version-api-respond-200", response: appVersion, expectedCode: 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/version", nil)

			resp, _ := app.Test(req)
			bodyBytes, _ := ioutil.ReadAll(resp.Body)

			assert.Equalf(t, tt.response, string(bodyBytes), tt.name)
			assert.Equalf(t, tt.expectedCode, resp.StatusCode, tt.name)
		})
	}
}
