package integration

import (
	"d-payroll/controller/http/dto"
	"d-payroll/entity"
	"d-payroll/utils"
	"encoding/json"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	// Mock time.Now to be Monday at 9am
	originalTimeNow := utils.TimeNow
	defer func() { utils.TimeNow = originalTimeNow }()
	utils.TimeNow = func() time.Time {
		return time.Date(2025, 6, 16, 9, 0, 0, 0, time.Local) // Monday, June 16, 2025 at 9:00 AM
	}

	// Set up the test app
	testApp, err := SetupTestApp(t)
	if err != nil {
		t.Fatalf("Failed to set up test app: %v", err)
	}
	defer testApp.TeardownTestApp()

	// Prepare request body for creating a new employee user
	salary := 4000000
	createUserRequest := dto.CreateUserBodyDto{
		Username: "employee1",
		Password: "password123",
		Role:     string(entity.UserRoleEmployee),
		UserInfo: &dto.CreateUserInfoBodyDto{
			MonthlySalary: &salary,
		},
	}

	// Convert to JSON
	requestBody, err := json.Marshal(createUserRequest)
	require.NoError(t, err, "Failed to marshal request body")

	// Create a new request with admin token for authorization
	req, err := testApp.makeAuthenticatedRequest("POST", "/users", requestBody, testApp.AdminToken)
	require.NoError(t, err, "Failed to create request")
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := testApp.App.Test(req)
	require.NoError(t, err, "Failed to test request")

	// Check the response status code
	assert.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected status code to be 201 Created")

	// Parse response body
	var response entity.HttpResponse
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &response)
	require.NoError(t, err, "Failed to parse response body")

	fmt.Println(response)

	// Check response content
	assert.True(t, response.Success, "Expected success to be true")
	assert.NotNil(t, response.Data, "Expected data to not be nil")

	// Verify the created user data
	var userData map[string]interface{}
	userDataJson, err := json.Marshal(response.Data)
	require.NoError(t, err, "Failed to marshal user data")
	err = json.Unmarshal(userDataJson, &userData)
	require.NoError(t, err, "Failed to unmarshal user data")

	// Check user properties
	assert.Equal(t, "employee1", userData["username"], "Username should match")
	assert.Equal(t, "EMPLOYEE", userData["role"], "Role should match")
	assert.NotNil(t, userData["id"], "User ID should not be nil")
	userInfo, ok := userData["user_info"].(map[string]interface{})
	require.True(t, ok, "user_info should be a map")
	assert.Equal(t, float64(4000000), userInfo["monthly_salary"], "Monthly salary should match")
}
