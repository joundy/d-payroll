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

func TestCreateAndApproveOvertime(t *testing.T) {
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
		Username: "employee-overtime",
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

	// Create an attendance check-in for the employee (required before overtime)
	// This simulates that the employee has checked in and checked out
	req, err := testApp.makeAuthenticatedRequest("POST", "/attendances/checkin", nil, employeeToken)
	require.NoError(t, err, "Failed to create check-in request")
	resp, err := testApp.App.Test(req)
	require.NoError(t, err, "Failed to test check-in request")
	require.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected check-in to succeed")

	// Create a checkout to avoid the "overtime submit before checkout" error
	req, err = testApp.makeAuthenticatedRequest("POST", "/attendances/checkout", nil, employeeToken)
	require.NoError(t, err, "Failed to create check-out request")
	resp, err = testApp.App.Test(req)
	require.NoError(t, err, "Failed to test check-out request")
	require.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected check-out to succeed")

	// Test successful overtime creation
	t.Run("Successful Overtime Creation", func(t *testing.T) {
		// Prepare overtime request
		overtimeRequest := dto.CreateOvertimeBodyDto{
			Description:   "Working on urgent project",
			OvertimeAt:    time.Now(),
			DurationMilis: 7200000, // 2 hours in milliseconds
		}

		// Convert to JSON
		requestBody, err := json.Marshal(overtimeRequest)
		require.NoError(t, err, "Failed to marshal overtime request")

		// Create a new request for overtime creation
		req, err := testApp.makeAuthenticatedRequest("POST", "/overtimes", requestBody, employeeToken)
		require.NoError(t, err, "Failed to create overtime request")
		req.Header.Set("Content-Type", "application/json")

		// Perform the request
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test overtime request")

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

		// Verify the overtime data
		var overtimeData map[string]interface{}
		overtimeDataJson, err := json.Marshal(response.Data)
		require.NoError(t, err, "Failed to marshal overtime data")
		err = json.Unmarshal(overtimeDataJson, &overtimeData)
		require.NoError(t, err, "Failed to unmarshal overtime data")

		// Check overtime properties
		assert.Equal(t, "Working on urgent project", overtimeData["description"], "Description should match")
		assert.Equal(t, float64(7200000), overtimeData["duration_milis"], "Duration should match")
		assert.Equal(t, false, overtimeData["is_approved"], "Overtime should not be approved initially")
		assert.NotNil(t, overtimeData["id"], "Overtime ID should not be nil")

		// Store overtime ID for approval test
		overtimeId := overtimeData["id"].(float64)

		// Test overtime approval by admin
		t.Run("Successful Overtime Approval", func(t *testing.T) {
			// Create a new request for overtime approval
			approvalPath := fmt.Sprintf("/overtimes/%d/approve", int(overtimeId))
			req, err := testApp.makeAuthenticatedRequest("POST", approvalPath, nil, testApp.AdminToken)
			require.NoError(t, err, "Failed to create approval request")

			// Perform the request
			resp, err := testApp.App.Test(req)
			require.NoError(t, err, "Failed to test approval request")

			// Check the response status code
			assert.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected status code to be 200 OK")

			// Parse response body
			var response entity.HttpResponse
			body, _ := io.ReadAll(resp.Body)
			err = json.Unmarshal(body, &response)
			require.NoError(t, err, "Failed to parse response body")

			// Check response content
			assert.True(t, response.Success, "Expected success to be true")

			// Test already approved error
			t.Run("Already Approved Error", func(t *testing.T) {
				// Create a new request for overtime approval again
				req, err := testApp.makeAuthenticatedRequest("POST", approvalPath, nil, testApp.AdminToken)
				require.NoError(t, err, "Failed to create approval request")

				// Perform the request
				resp, err := testApp.App.Test(req)
				require.NoError(t, err, "Failed to test approval request")

				// Check the response status code - should be conflict
				assert.Equal(t, fiber.StatusConflict, resp.StatusCode, "Expected status code to be 409 Conflict")

				// Parse response body
				var response entity.HttpResponse
				body, _ := io.ReadAll(resp.Body)
				err = json.Unmarshal(body, &response)
				require.NoError(t, err, "Failed to parse response body")

				// Check response content
				assert.False(t, response.Success, "Expected success to be false")
				assert.Equal(t, "Overtime already approved", response.Message, "Expected already approved message")
			})
		})
	})
}

func TestOvertimeValidationErrors(t *testing.T) {
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
		Username: "employee-overtime-validation",
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

	// Test missing required fields
	t.Run("Missing Required Fields", func(t *testing.T) {
		// Prepare invalid overtime request (missing description)
		invalidRequest := map[string]interface{}{
			"overtime_at":    time.Now(),
			"duration_milis": 7200000,
			// Missing description
		}

		// Convert to JSON
		requestBody, err := json.Marshal(invalidRequest)
		require.NoError(t, err, "Failed to marshal overtime request")

		// Create a new request
		req, err := testApp.makeAuthenticatedRequest("POST", "/overtimes", requestBody, employeeToken)
		require.NoError(t, err, "Failed to create overtime request")
		req.Header.Set("Content-Type", "application/json")

		// Perform the request
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test overtime request")

		// Check the response status code - should be bad request or unprocessable entity
		assert.True(t, resp.StatusCode == fiber.StatusBadRequest || resp.StatusCode == fiber.StatusUnprocessableEntity,
			"Expected status code to indicate validation error")
	})

	// Test invalid duration
	t.Run("Invalid Duration", func(t *testing.T) {
		// Prepare invalid overtime request (duration <= 0)
		invalidRequest := dto.CreateOvertimeBodyDto{
			Description:   "Working on urgent project",
			OvertimeAt:    time.Now(),
			DurationMilis: 0, // Invalid duration
		}

		// Convert to JSON
		requestBody, err := json.Marshal(invalidRequest)
		require.NoError(t, err, "Failed to marshal overtime request")

		// Create a new request
		req, err := testApp.makeAuthenticatedRequest("POST", "/overtimes", requestBody, employeeToken)
		require.NoError(t, err, "Failed to create overtime request")
		req.Header.Set("Content-Type", "application/json")

		// Perform the request
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test overtime request")

		// Check the response status code - should indicate validation error
		assert.True(t, resp.StatusCode == fiber.StatusBadRequest || resp.StatusCode == fiber.StatusUnprocessableEntity,
			"Expected status code to indicate validation error")
	})

	// Test overtime submit before checkout
	t.Run("Overtime Submit Before Checkout", func(t *testing.T) {
		// Create an attendance check-in for the employee
		req, err := testApp.makeAuthenticatedRequest("POST", "/attendances/checkin", nil, employeeToken)
		require.NoError(t, err, "Failed to create check-in request")
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test check-in request")
		require.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected check-in to succeed")

		// Prepare overtime request
		overtimeRequest := dto.CreateOvertimeBodyDto{
			Description:   "Working on urgent project",
			OvertimeAt:    time.Now(),
			DurationMilis: 7200000,
		}

		// Convert to JSON
		requestBody, err := json.Marshal(overtimeRequest)
		require.NoError(t, err, "Failed to marshal overtime request")

		// Create a new request for overtime creation before checkout
		req, err = testApp.makeAuthenticatedRequest("POST", "/overtimes", requestBody, employeeToken)
		require.NoError(t, err, "Failed to create overtime request")
		req.Header.Set("Content-Type", "application/json")

		// Perform the request
		resp, err = testApp.App.Test(req)
		require.NoError(t, err, "Failed to test overtime request")

		// Check the response status code - should be unprocessable entity
		assert.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode, "Expected status code to be 422 Unprocessable Entity")

		// Parse response body
		var response entity.HttpResponse
		body, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(body, &response)
		require.NoError(t, err, "Failed to parse response body")

		// Check response content
		assert.False(t, response.Success, "Expected success to be false")
		assert.Equal(t, "Overtime submit before checkout", response.Message, "Expected overtime submit before checkout message")
	})
}

func TestOvertimeOnWeekend(t *testing.T) {
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
		Username: "employee-weekend-overtime",
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

	// Test successful overtime creation on weekend (no need for check-in/check-out)
	// Prepare overtime request
	overtimeRequest := dto.CreateOvertimeBodyDto{
		Description:   "Weekend project work",
		OvertimeAt:    utils.TimeNow(),
		DurationMilis: 7200000, // 2 hours in milliseconds
	}

	// Convert to JSON
	requestBody, err := json.Marshal(overtimeRequest)
	require.NoError(t, err, "Failed to marshal overtime request")

	// Create a new request for overtime creation
	req, err := testApp.makeAuthenticatedRequest("POST", "/overtimes", requestBody, employeeToken)
	require.NoError(t, err, "Failed to create overtime request")
	req.Header.Set("Content-Type", "application/json")

	// Perform the request
	resp, err := testApp.App.Test(req)
	require.NoError(t, err, "Failed to test overtime request")

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

	// Verify the overtime data
	var overtimeData map[string]interface{}
	overtimeDataJson, err := json.Marshal(response.Data)
	require.NoError(t, err, "Failed to marshal overtime data")
	err = json.Unmarshal(overtimeDataJson, &overtimeData)
	require.NoError(t, err, "Failed to unmarshal overtime data")

	// Check overtime properties
	assert.Equal(t, "Weekend project work", overtimeData["description"], "Description should match")
	assert.Equal(t, float64(7200000), overtimeData["duration_milis"], "Duration should match")
	assert.Equal(t, false, overtimeData["is_approved"], "Overtime should not be approved initially")
	assert.NotNil(t, overtimeData["id"], "Overtime ID should not be nil")
}

func TestOvertimeAuthorizationErrors(t *testing.T) {
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
		req, err := testApp.makeAuthenticatedRequest("POST", "/overtimes", nil, "")
		require.NoError(t, err, "Failed to create overtime request")

		// Perform the request
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test overtime request")

		// Check the response status code - should be unauthorized
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode, "Expected status code to be 401 Unauthorized")
	})

	// Test unauthorized role (admin trying to create overtime)
	t.Run("Admin Creating Overtime", func(t *testing.T) {
		// Prepare overtime request
		overtimeRequest := dto.CreateOvertimeBodyDto{
			Description:   "Working on urgent project",
			OvertimeAt:    time.Now(),
			DurationMilis: 7200000,
		}

		// Convert to JSON
		requestBody, err := json.Marshal(overtimeRequest)
		require.NoError(t, err, "Failed to marshal overtime request")

		// Create a new request with admin token
		req, err := testApp.makeAuthenticatedRequest("POST", "/overtimes", requestBody, testApp.AdminToken)
		require.NoError(t, err, "Failed to create overtime request")
		req.Header.Set("Content-Type", "application/json")

		// Perform the request
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test overtime request")

		// Check the response status code - should be forbidden
		assert.Equal(t, fiber.StatusForbidden, resp.StatusCode, "Expected status code to be 403 Forbidden")
	})

	// Test unauthorized role (employee trying to approve overtime)
	t.Run("Employee Approving Overtime", func(t *testing.T) {
		// Create a test employee user
		salary := 4500000
		testUser := &entity.User{
			Username: "employee-approve-attempt",
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

		// Create a new request for overtime approval with employee token
		req, err := testApp.makeAuthenticatedRequest("POST", "/overtimes/1/approve", nil, employeeToken)
		require.NoError(t, err, "Failed to create approval request")

		// Perform the request
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test approval request")

		// Check the response status code - should be forbidden
		assert.Equal(t, fiber.StatusForbidden, resp.StatusCode, "Expected status code to be 403 Forbidden")
	})
}
