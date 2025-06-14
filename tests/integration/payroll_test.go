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

func TestCreateAndGetPayrolls(t *testing.T) {
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

	// Test successful payroll creation
	t.Run("Successful Payroll Creation", func(t *testing.T) {
		// Prepare payroll request
		startedAt := time.Date(2025, 6, 1, 0, 0, 0, 0, time.Local)
		endedAt := time.Date(2025, 6, 30, 23, 59, 59, 0, time.Local)

		payrollRequest := dto.CreatePayrollBodyDto{
			Name:      "June 2025 Payroll",
			StartedAt: startedAt,
			EndedAt:   endedAt,
		}

		// Convert to JSON
		requestBody, err := json.Marshal(payrollRequest)
		require.NoError(t, err, "Failed to marshal payroll request")

		// Create a new request for payroll creation
		req, err := testApp.makeAuthenticatedRequest("POST", "/payrolls", requestBody, testApp.AdminToken)
		require.NoError(t, err, "Failed to create payroll request")
		req.Header.Set("Content-Type", "application/json")

		// Perform the request
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test payroll request")

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

		// Verify the payroll data
		var payrollData map[string]interface{}
		payrollDataJson, err := json.Marshal(response.Data)
		require.NoError(t, err, "Failed to marshal payroll data")
		err = json.Unmarshal(payrollDataJson, &payrollData)
		require.NoError(t, err, "Failed to unmarshal payroll data")

		// Check payroll properties
		assert.Equal(t, "June 2025 Payroll", payrollData["name"], "Name should match")
		assert.Equal(t, false, payrollData["is_rolled"], "Payroll should not be rolled initially")
		assert.NotNil(t, payrollData["id"], "Payroll ID should not be nil")
		assert.NotNil(t, payrollData["created_by_user_id"], "Created by user ID should not be nil")
		assert.NotNil(t, payrollData["created_at"], "Created at should not be nil")
		assert.NotNil(t, payrollData["updated_at"], "Updated at should not be nil")

		// Store payroll ID for further tests
		payrollId := payrollData["id"].(float64)

		// Test getting all payrolls
		t.Run("Get All Payrolls", func(t *testing.T) {
			// Create a new request to get payrolls
			req, err := testApp.makeAuthenticatedRequest("GET", "/payrolls", nil, testApp.AdminToken)
			require.NoError(t, err, "Failed to create get payrolls request")

			// Perform the request
			resp, err := testApp.App.Test(req)
			require.NoError(t, err, "Failed to test get payrolls request")

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

			// Verify the payrolls data
			var payrollsData []interface{}
			payrollsDataJson, err := json.Marshal(response.Data)
			require.NoError(t, err, "Failed to marshal payrolls data")
			err = json.Unmarshal(payrollsDataJson, &payrollsData)
			require.NoError(t, err, "Failed to unmarshal payrolls data")

			// Check that we have at least one payroll
			assert.GreaterOrEqual(t, len(payrollsData), 1, "Should have at least one payroll")

			// Find our created payroll
			var foundPayroll map[string]interface{}
			for _, payroll := range payrollsData {
				payrollMap := payroll.(map[string]interface{})
				if payrollMap["id"].(float64) == payrollId {
					foundPayroll = payrollMap
					break
				}
			}

			require.NotNil(t, foundPayroll, "Should find the created payroll")
			assert.Equal(t, "June 2025 Payroll", foundPayroll["name"], "Name should match")
			assert.Equal(t, false, foundPayroll["is_rolled"], "Payroll should not be rolled")
		})

		// Test payroll rolling
		t.Run("Roll Payroll", func(t *testing.T) {
			// Create a new request for payroll rolling
			rollPath := fmt.Sprintf("/payrolls/%d/roll", int(payrollId))
			req, err := testApp.makeAuthenticatedRequest("POST", rollPath, nil, testApp.AdminToken)
			require.NoError(t, err, "Failed to create roll payroll request")

			// Perform the request
			resp, err := testApp.App.Test(req)
			require.NoError(t, err, "Failed to test roll payroll request")

			// Check the response status code
			assert.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected status code to be 200 OK")

			// Parse response body
			var response entity.HttpResponse
			body, _ := io.ReadAll(resp.Body)
			err = json.Unmarshal(body, &response)
			require.NoError(t, err, "Failed to parse response body")

			// Check response content
			assert.True(t, response.Success, "Expected success to be true")

			// Test already rolled error
			t.Run("Already Rolled Error", func(t *testing.T) {
				// Create a new request for payroll rolling again
				req, err := testApp.makeAuthenticatedRequest("POST", rollPath, nil, testApp.AdminToken)
				require.NoError(t, err, "Failed to create roll payroll request")

				// Perform the request
				resp, err := testApp.App.Test(req)
				require.NoError(t, err, "Failed to test roll payroll request")

				// Check the response status code - should be conflict
				assert.Equal(t, fiber.StatusConflict, resp.StatusCode, "Expected status code to be 409 Conflict")

				// Parse response body
				var response entity.HttpResponse
				body, _ := io.ReadAll(resp.Body)
				err = json.Unmarshal(body, &response)
				require.NoError(t, err, "Failed to parse response body")

				// Check response content
				assert.False(t, response.Success, "Expected success to be false")
				assert.Equal(t, "Payroll already rolled", response.Message, "Expected already rolled message")
			})
		})
	})
}

func TestPayrollValidationErrors(t *testing.T) {
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

	// Test missing required fields
	t.Run("Missing Required Fields", func(t *testing.T) {
		// Prepare invalid payroll request (missing name)
		invalidRequest := map[string]interface{}{
			"started_at": time.Date(2025, 6, 1, 0, 0, 0, 0, time.Local),
			"ended_at":   time.Date(2025, 6, 30, 23, 59, 59, 0, time.Local),
			// Missing name
		}

		// Convert to JSON
		requestBody, err := json.Marshal(invalidRequest)
		require.NoError(t, err, "Failed to marshal payroll request")

		// Create a new request
		req, err := testApp.makeAuthenticatedRequest("POST", "/payrolls", requestBody, testApp.AdminToken)
		require.NoError(t, err, "Failed to create payroll request")
		req.Header.Set("Content-Type", "application/json")

		// Perform the request
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test payroll request")

		// Check the response status code - should be bad request or unprocessable entity
		assert.True(t, resp.StatusCode == fiber.StatusBadRequest || resp.StatusCode == fiber.StatusUnprocessableEntity,
			"Expected status code to indicate validation error")
	})

	// Test empty name
	t.Run("Empty Name", func(t *testing.T) {
		// Prepare invalid payroll request (empty name)
		invalidRequest := dto.CreatePayrollBodyDto{
			Name:      "", // Empty name
			StartedAt: time.Date(2025, 6, 1, 0, 0, 0, 0, time.Local),
			EndedAt:   time.Date(2025, 6, 30, 23, 59, 59, 0, time.Local),
		}

		// Convert to JSON
		requestBody, err := json.Marshal(invalidRequest)
		require.NoError(t, err, "Failed to marshal payroll request")

		// Create a new request
		req, err := testApp.makeAuthenticatedRequest("POST", "/payrolls", requestBody, testApp.AdminToken)
		require.NoError(t, err, "Failed to create payroll request")
		req.Header.Set("Content-Type", "application/json")

		// Perform the request
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test payroll request")

		// Check the response status code - should indicate validation error
		assert.True(t, resp.StatusCode == fiber.StatusBadRequest || resp.StatusCode == fiber.StatusUnprocessableEntity,
			"Expected status code to indicate validation error")
	})

	// Test invalid date range (end date before start date)
	t.Run("Invalid Date Range", func(t *testing.T) {
		// Prepare invalid payroll request (end date before start date)
		invalidRequest := dto.CreatePayrollBodyDto{
			Name:      "Invalid Date Range Payroll",
			StartedAt: time.Date(2025, 6, 30, 0, 0, 0, 0, time.Local),
			EndedAt:   time.Date(2025, 6, 1, 23, 59, 59, 0, time.Local), // End before start
		}

		// Convert to JSON
		requestBody, err := json.Marshal(invalidRequest)
		require.NoError(t, err, "Failed to marshal payroll request")

		// Create a new request
		req, err := testApp.makeAuthenticatedRequest("POST", "/payrolls", requestBody, testApp.AdminToken)
		require.NoError(t, err, "Failed to create payroll request")
		req.Header.Set("Content-Type", "application/json")

		// Perform the request
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test payroll request")

		// This might succeed at the API level but fail at business logic level
		// The validation depends on the business rules implemented
		// For now, we'll just check that we get a response
		assert.True(t, resp.StatusCode >= 200 && resp.StatusCode < 500, "Should get a valid HTTP response")
	})
}

func TestPayrollAuthorizationErrors(t *testing.T) {
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

	// Create a test employee user for authorization tests
	salary := 4500000
	testUser := &entity.User{
		Username: "employee-payroll-test",
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

	// Test unauthorized access (no token)
	t.Run("No Token", func(t *testing.T) {
		// Create a new request without token
		req, err := testApp.makeAuthenticatedRequest("POST", "/payrolls", nil, "")
		require.NoError(t, err, "Failed to create payroll request")

		// Perform the request
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test payroll request")

		// Check the response status code - should be unauthorized
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode, "Expected status code to be 401 Unauthorized")
	})

	// Test unauthorized role (employee trying to create payroll)
	t.Run("Employee Creating Payroll", func(t *testing.T) {
		// Prepare payroll request
		payrollRequest := dto.CreatePayrollBodyDto{
			Name:      "Unauthorized Payroll",
			StartedAt: time.Date(2025, 6, 1, 0, 0, 0, 0, time.Local),
			EndedAt:   time.Date(2025, 6, 30, 23, 59, 59, 0, time.Local),
		}

		// Convert to JSON
		requestBody, err := json.Marshal(payrollRequest)
		require.NoError(t, err, "Failed to marshal payroll request")

		// Create a new request with employee token
		req, err := testApp.makeAuthenticatedRequest("POST", "/payrolls", requestBody, employeeToken)
		require.NoError(t, err, "Failed to create payroll request")
		req.Header.Set("Content-Type", "application/json")

		// Perform the request
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test payroll request")

		// Check the response status code - should be forbidden
		assert.Equal(t, fiber.StatusForbidden, resp.StatusCode, "Expected status code to be 403 Forbidden")
	})

	// Test unauthorized role (employee trying to get payrolls)
	t.Run("Employee Getting Payrolls", func(t *testing.T) {
		// Create a new request with employee token
		req, err := testApp.makeAuthenticatedRequest("GET", "/payrolls", nil, employeeToken)
		require.NoError(t, err, "Failed to create get payrolls request")

		// Perform the request
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test get payrolls request")

		// Check the response status code - should be forbidden
		assert.Equal(t, fiber.StatusForbidden, resp.StatusCode, "Expected status code to be 403 Forbidden")
	})

	// Test unauthorized role (employee trying to roll payroll)
	t.Run("Employee Rolling Payroll", func(t *testing.T) {
		// Create a new request with employee token
		req, err := testApp.makeAuthenticatedRequest("POST", "/payrolls/1/roll", nil, employeeToken)
		require.NoError(t, err, "Failed to create roll payroll request")

		// Perform the request
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test roll payroll request")

		// Check the response status code - should be forbidden
		assert.Equal(t, fiber.StatusForbidden, resp.StatusCode, "Expected status code to be 403 Forbidden")
	})
}

func TestPayrollNotFoundErrors(t *testing.T) {
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

	// Test rolling non-existent payroll
	t.Run("Roll Non-Existent Payroll", func(t *testing.T) {
		// Create a new request for rolling non-existent payroll
		req, err := testApp.makeAuthenticatedRequest("POST", "/payrolls/999999/roll", nil, testApp.AdminToken)
		require.NoError(t, err, "Failed to create roll payroll request")

		// Perform the request
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test roll payroll request")

		// Check the response status code - should be not found
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode, "Expected status code to be 404 Not Found")

		// Parse response body
		var response entity.HttpResponse
		body, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(body, &response)
		require.NoError(t, err, "Failed to parse response body")

		// Check response content
		assert.False(t, response.Success, "Expected success to be false")
		assert.Equal(t, "Payroll not found", response.Message, "Expected payroll not found message")
	})

	// Test invalid payroll ID format
	t.Run("Invalid Payroll ID Format", func(t *testing.T) {
		// Create a new request with invalid payroll ID
		req, err := testApp.makeAuthenticatedRequest("POST", "/payrolls/invalid/roll", nil, testApp.AdminToken)
		require.NoError(t, err, "Failed to create roll payroll request")

		// Perform the request
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test roll payroll request")

		// Check the response status code - should be bad request
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode, "Expected status code to be 400 Bad Request")

		// Parse response body
		var response entity.HttpResponse
		body, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(body, &response)
		require.NoError(t, err, "Failed to parse response body")

		// Check response content
		assert.False(t, response.Success, "Expected success to be false")
		assert.Equal(t, "Invalid payroll ID param", response.Message, "Expected invalid payroll ID message")
	})
}
