package integration

import (
	"d-payroll/entity"
	"encoding/json"
	"io"
	nethttp "net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthEndpoint(t *testing.T) {
	// Set up the test app
	testApp, err := SetupTestApp(t)
	if err != nil {
		t.Fatalf("Failed to set up test app: %v", err)
	}
	defer testApp.TeardownTestApp()

	// Create a new request
	req, _ := nethttp.NewRequest("GET", "/_health", nil)

	// Perform the request
	resp, err := testApp.App.Test(req)
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}

	// Check the response status code
	assert.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected status code to be 200 OK")

	// Parse response body
	var response entity.HttpResponse
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	// Validate response fields
	assert.True(t, response.Success)
	assert.Equal(t, "Ok", response.Message)
}
