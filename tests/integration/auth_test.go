package integration

import (
	"d-payroll/controller/http/dto"
	"d-payroll/entity"
	"encoding/json"
	"io"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	// Set up the test app
	testApp, err := SetupTestApp(t)
	if err != nil {
		t.Fatalf("Failed to set up test app: %v", err)
	}
	defer testApp.TeardownTestApp()

	// First, create a test user to login with
	salary := 4500000
	testUser := &entity.User{
		Username: "testlogin",
		Password: "password123",
		Role:     entity.UserRoleEmployee,
		UserInfo: &entity.UserInfo{
			MonthlySalary: &salary,
		},
	}

	// Create the user using the service directly
	createdUser, err := testApp.UserService.CreateUser(testApp.ctx, testUser)
	require.NoError(t, err, "Failed to create test user")
	require.NotNil(t, createdUser, "Created user should not be nil")
	require.NotNil(t, createdUser.Id, "Created user ID should not be nil")

	// Prepare login request
	loginRequest := dto.LoginBodyDto{
		Username: "testlogin",
		Password: "password123",
	}

	// Convert to JSON
	requestBody, err := json.Marshal(loginRequest)
	require.NoError(t, err, "Failed to marshal login request")

	// Create a new request
	req, err := testApp.makeAuthenticatedRequest("POST", "/login", requestBody, "")
	require.NoError(t, err, "Failed to create request")
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := testApp.App.Test(req)
	require.NoError(t, err, "Failed to test request")

	// Check the response status code
	assert.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected status code to be 200 OK")

	// Parse response body
	var response entity.HttpResponse
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &response)
	require.NoError(t, err, "Failed to parse response body")

	// Check response content
	assert.True(t, response.Success, "Expected success to be true")
	assert.NotNil(t, response.Data, "Expected data to not be nil")

	// Verify the token is returned
	var loginResponse map[string]interface{}
	loginDataJson, err := json.Marshal(response.Data)
	require.NoError(t, err, "Failed to marshal login data")
	err = json.Unmarshal(loginDataJson, &loginResponse)
	require.NoError(t, err, "Failed to unmarshal login data")

	// Check token is present
	token, exists := loginResponse["token"]
	assert.True(t, exists, "Token should exist in response")
	assert.NotEmpty(t, token, "Token should not be empty")
}

func TestLoginInvalidCredentials(t *testing.T) {
	// Set up the test app
	testApp, err := SetupTestApp(t)
	if err != nil {
		t.Fatalf("Failed to set up test app: %v", err)
	}
	defer testApp.TeardownTestApp()

	// Prepare login request with invalid credentials
	loginRequest := dto.LoginBodyDto{
		Username: "admin",
		Password: "wrongpassword",
	}

	// Convert to JSON
	requestBody, err := json.Marshal(loginRequest)
	require.NoError(t, err, "Failed to marshal login request")

	// Create a new request
	req, err := testApp.makeAuthenticatedRequest("POST", "/login", requestBody, "")
	require.NoError(t, err, "Failed to create request")
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := testApp.App.Test(req)
	require.NoError(t, err, "Failed to test request")

	// Check the response status code - should be unauthorized
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode, "Expected status code to be 401 Unauthorized")

	// Parse response body
	var response entity.HttpResponse
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &response)
	require.NoError(t, err, "Failed to parse response body")

	// Check response content
	assert.False(t, response.Success, "Expected success to be false")
	assert.Equal(t, "Invalid credentials", response.Message, "Expected invalid credentials message")
}
