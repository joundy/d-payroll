package integration

import (
	"d-payroll/entity"
	"d-payroll/utils"
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAttendanceCheckinCheckout(t *testing.T) {
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

	// Create a test employee user
	salary := 4500000
	testUser := &entity.User{
		Username: "employee-attendance",
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

	// Generate token for the employee
	employeeToken, err := utils.GenerateToken(testApp.Config.Auth.JwtSecret, &entity.AuthTokenPayload{
		ID:   *createdUser.Id,
		Role: createdUser.Role,
	})
	require.NoError(t, err, "Failed to generate employee token")

	// Test successful check-in
	t.Run("Successful Check-in", func(t *testing.T) {
		// Create a new request for check-in
		req, err := testApp.makeAuthenticatedRequest("POST", "/attendances/checkin", nil, employeeToken)
		require.NoError(t, err, "Failed to create check-in request")

		// Perform the request
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test check-in request")

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

		// Verify the attendance data
		var attendanceData map[string]interface{}
		attendanceDataJson, err := json.Marshal(response.Data)
		require.NoError(t, err, "Failed to marshal attendance data")
		err = json.Unmarshal(attendanceDataJson, &attendanceData)
		require.NoError(t, err, "Failed to unmarshal attendance data")

		// Check attendance properties
		assert.Equal(t, "CHECKIN", attendanceData["type"], "Attendance type should be CHECKIN")
		assert.NotNil(t, attendanceData["id"], "Attendance ID should not be nil")
		assert.NotNil(t, attendanceData["created_at"], "Created at should not be nil")
	})

	// Test already checked in error
	t.Run("Already Checked In Error", func(t *testing.T) {
		// Create a new request for check-in again
		req, err := testApp.makeAuthenticatedRequest("POST", "/attendances/checkin", nil, employeeToken)
		require.NoError(t, err, "Failed to create check-in request")

		// Perform the request
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test check-in request")

		// Check the response status code - should be conflict
		assert.Equal(t, fiber.StatusConflict, resp.StatusCode, "Expected status code to be 409 Conflict")

		// Parse response body
		var response entity.HttpResponse
		body, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(body, &response)
		require.NoError(t, err, "Failed to parse response body")

		// Check response content
		assert.False(t, response.Success, "Expected success to be false")
		assert.Equal(t, "User already checked in", response.Message, "Expected already checked in message")
	})

	// Test successful check-out
	t.Run("Successful Check-out", func(t *testing.T) {
		// Create a new request for check-out
		req, err := testApp.makeAuthenticatedRequest("POST", "/attendances/checkout", nil, employeeToken)
		require.NoError(t, err, "Failed to create check-out request")

		// Perform the request
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test check-out request")

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

		// Verify the attendance data
		var attendanceData map[string]interface{}
		attendanceDataJson, err := json.Marshal(response.Data)
		require.NoError(t, err, "Failed to marshal attendance data")
		err = json.Unmarshal(attendanceDataJson, &attendanceData)
		require.NoError(t, err, "Failed to unmarshal attendance data")

		// Check attendance properties
		assert.Equal(t, "CHECKOUT", attendanceData["type"], "Attendance type should be CHECKOUT")
		assert.NotNil(t, attendanceData["id"], "Attendance ID should not be nil")
		assert.NotNil(t, attendanceData["created_at"], "Created at should not be nil")
	})

	// Test already checked out error
	t.Run("Already Checked Out Error", func(t *testing.T) {
		// Create a new request for check-out again
		req, err := testApp.makeAuthenticatedRequest("POST", "/attendances/checkout", nil, employeeToken)
		require.NoError(t, err, "Failed to create check-out request")

		// Perform the request
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test check-out request")

		// Check the response status code - should be conflict
		assert.Equal(t, fiber.StatusConflict, resp.StatusCode, "Expected status code to be 409 Conflict")

		// Parse response body
		var response entity.HttpResponse
		body, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(body, &response)
		require.NoError(t, err, "Failed to parse response body")

		// Check response content
		assert.False(t, response.Success, "Expected success to be false")
		assert.Equal(t, "User already checked out", response.Message, "Expected already checked out message")
	})
}

func TestAttendanceCheckoutWithoutCheckin(t *testing.T) {
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

	// Create a test employee user
	salary := 4500000
	testUser := &entity.User{
		Username: "employee-checkout-error",
		Password: "password123",
		Role:     entity.UserRoleEmployee,
		UserInfo: &entity.UserInfo{
			MonthlySalary: &salary,
		},
	}

	// Create the user using the service directly
	createdUser, err := testApp.UserService.CreateUser(testApp.ctx, testUser)
	require.NoError(t, err, "Failed to create test user")

	// Generate token for the employee
	employeeToken, err := utils.GenerateToken(testApp.Config.Auth.JwtSecret, &entity.AuthTokenPayload{
		ID:   *createdUser.Id,
		Role: createdUser.Role,
	})
	require.NoError(t, err, "Failed to generate employee token")

	// Test checkout without checkin error
	// Create a new request for check-out without check-in
	req, err := testApp.makeAuthenticatedRequest("POST", "/attendances/checkout", nil, employeeToken)
	require.NoError(t, err, "Failed to create check-out request")

	// Perform the request
	resp, err := testApp.App.Test(req)
	require.NoError(t, err, "Failed to test check-out request")

	// Check the response status code - should be unprocessable entity
	assert.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode, "Expected status code to be 422 Unprocessable Entity")

	// Parse response body
	var response entity.HttpResponse
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &response)
	require.NoError(t, err, "Failed to parse response body")

	// Check response content
	assert.False(t, response.Success, "Expected success to be false")
	assert.Equal(t, "User cannot checked out because it is not checked in", response.Message, "Expected cannot check out message")
}

func TestAttendanceWeekendError(t *testing.T) {
	// Mock time.Now to be Saturday at 9am
	originalTimeNow := utils.TimeNow
	defer func() { utils.TimeNow = originalTimeNow }()
	utils.TimeNow = func() time.Time {
		return time.Date(2025, 6, 14, 9, 0, 0, 0, time.Local) // Saturday, June 14, 2025 at 9:00 AM
	}

	// Set up the test app
	testApp, err := SetupTestApp(t)
	if err != nil {
		t.Fatalf("Failed to set up test app: %v", err)
	}
	defer testApp.TeardownTestApp()

	// Create a test employee user
	salary := 4500000
	testUser := &entity.User{
		Username: "employee-weekend",
		Password: "password123",
		Role:     entity.UserRoleEmployee,
		UserInfo: &entity.UserInfo{
			MonthlySalary: &salary,
		},
	}

	// Create the user using the service directly
	createdUser, err := testApp.UserService.CreateUser(testApp.ctx, testUser)
	require.NoError(t, err, "Failed to create test user")

	// Generate token for the employee
	employeeToken, err := utils.GenerateToken(testApp.Config.Auth.JwtSecret, &entity.AuthTokenPayload{
		ID:   *createdUser.Id,
		Role: createdUser.Role,
	})
	require.NoError(t, err, "Failed to generate employee token")

	// Create a new request for check-in on weekend
	req, err := testApp.makeAuthenticatedRequest("POST", "/attendances/checkin", nil, employeeToken)
	require.NoError(t, err, "Failed to create check-in request")

	// Perform the request
	resp, err := testApp.App.Test(req)
	require.NoError(t, err, "Failed to test check-in request")

	// Check the response status code - should be unprocessable entity
	assert.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode, "Expected status code to be 422 Unprocessable Entity")
	
	// Parse response body
	var response entity.HttpResponse
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &response)
	require.NoError(t, err, "Failed to parse response body")

	// Check response content
	assert.False(t, response.Success, "Expected success to be false")
	assert.Equal(t, "User cannot checked in on weekend", response.Message, "Expected weekend error message")
}

func TestAttendanceUnauthorized(t *testing.T) {
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

	// Test unauthorized access (no token)
	t.Run("No Token", func(t *testing.T) {
		// Create a new request without token
		req, err := testApp.makeAuthenticatedRequest("POST", "/attendances/checkin", nil, "")
		require.NoError(t, err, "Failed to create check-in request")

		// Perform the request
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test check-in request")

		// Check the response status code - should be unauthorized
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode, "Expected status code to be 401 Unauthorized")
	})

	// Test unauthorized access (admin role trying to check-in)
	t.Run("Admin Role", func(t *testing.T) {
		// Create a new request with admin token
		req, err := testApp.makeAuthenticatedRequest("POST", "/attendances/checkin", nil, testApp.AdminToken)
		require.NoError(t, err, "Failed to create check-in request")

		// Perform the request
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test check-in request")

		// Check the response status code - should be forbidden
		assert.Equal(t, fiber.StatusForbidden, resp.StatusCode, "Expected status code to be 403 Forbidden")
	})
}
