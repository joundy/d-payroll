package integration

import (
	"d-payroll/controller/http/dto"
	"d-payroll/entity"
	"d-payroll/utils"
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAndApproveReimbursement(t *testing.T) {
	// Set up the test app
	testApp, err := SetupTestApp(t)
	if err != nil {
		t.Fatalf("Failed to set up test app: %v", err)
	}
	defer testApp.TeardownTestApp()

	// Create a test employee user
	salary := 4500000
	testUser := &entity.User{
		Username: "employee-reimbursement",
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

	// Test successful reimbursement creation
	t.Run("Successful Reimbursement Creation", func(t *testing.T) {
		// Prepare reimbursement request
		reimbursementRequest := dto.CreateReimbursementBodyDto{
			Description: "Office supplies purchase",
			Amount:      150000, // 150,000 (in the smallest currency unit, e.g., cents or sen)
		}

		// Convert to JSON
		requestBody, err := json.Marshal(reimbursementRequest)
		require.NoError(t, err, "Failed to marshal reimbursement request")

		// Create a new request for reimbursement creation
		req, err := testApp.makeAuthenticatedRequest("POST", "/reimbursements", requestBody, employeeToken)
		require.NoError(t, err, "Failed to create reimbursement request")
		req.Header.Set("Content-Type", "application/json")

		// Perform the request
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test reimbursement request")

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

		// Verify the reimbursement data
		var reimbursementData map[string]interface{}
		reimbursementDataJson, err := json.Marshal(response.Data)
		require.NoError(t, err, "Failed to marshal reimbursement data")
		err = json.Unmarshal(reimbursementDataJson, &reimbursementData)
		require.NoError(t, err, "Failed to unmarshal reimbursement data")

		// Check reimbursement properties
		assert.Equal(t, "Office supplies purchase", reimbursementData["description"], "Description should match")
		assert.Equal(t, float64(150000), reimbursementData["amount"], "Amount should match")
		assert.Equal(t, false, reimbursementData["is_approved"], "Reimbursement should not be approved initially")
		assert.NotNil(t, reimbursementData["id"], "Reimbursement ID should not be nil")

		// Store reimbursement ID for approval test
		reimbursementId := reimbursementData["id"].(float64)

		// Test reimbursement approval by admin
		t.Run("Successful Reimbursement Approval", func(t *testing.T) {
			// Create a new request for reimbursement approval
			approvalPath := fmt.Sprintf("/reimbursements/%d/approve", int(reimbursementId))
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
				// Create a new request for reimbursement approval again
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
				assert.Equal(t, "Reimbursement already approved", response.Message, "Expected already approved message")
			})
		})
	})
}

func TestReimbursementValidationErrors(t *testing.T) {
	// Set up the test app
	testApp, err := SetupTestApp(t)
	if err != nil {
		t.Fatalf("Failed to set up test app: %v", err)
	}
	defer testApp.TeardownTestApp()

	// Create a test employee user
	salary := 4500000
	testUser := &entity.User{
		Username: "employee-reimbursement-validation",
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
		// Prepare invalid reimbursement request (missing description)
		invalidRequest := map[string]interface{}{
			"amount": 150000,
			// Missing description
		}

		// Convert to JSON
		requestBody, err := json.Marshal(invalidRequest)
		require.NoError(t, err, "Failed to marshal reimbursement request")

		// Create a new request
		req, err := testApp.makeAuthenticatedRequest("POST", "/reimbursements", requestBody, employeeToken)
		require.NoError(t, err, "Failed to create reimbursement request")
		req.Header.Set("Content-Type", "application/json")

		// Perform the request
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test reimbursement request")

		// Check the response status code - should be bad request or unprocessable entity
		assert.True(t, resp.StatusCode == fiber.StatusBadRequest || resp.StatusCode == fiber.StatusUnprocessableEntity,
			"Expected status code to indicate validation error")
	})

	// Test invalid amount
	t.Run("Invalid Amount", func(t *testing.T) {
		// Prepare invalid reimbursement request (amount <= 0)
		invalidRequest := dto.CreateReimbursementBodyDto{
			Description: "Office supplies",
			Amount:      0, // Invalid amount
		}

		// Convert to JSON
		requestBody, err := json.Marshal(invalidRequest)
		require.NoError(t, err, "Failed to marshal reimbursement request")

		// Create a new request
		req, err := testApp.makeAuthenticatedRequest("POST", "/reimbursements", requestBody, employeeToken)
		require.NoError(t, err, "Failed to create reimbursement request")
		req.Header.Set("Content-Type", "application/json")

		// Perform the request
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test reimbursement request")

		// Check the response status code - should indicate validation error
		assert.True(t, resp.StatusCode == fiber.StatusBadRequest || resp.StatusCode == fiber.StatusUnprocessableEntity,
			"Expected status code to indicate validation error")
	})
}

func TestReimbursementAuthorizationErrors(t *testing.T) {
	// Set up the test app
	testApp, err := SetupTestApp(t)
	if err != nil {
		t.Fatalf("Failed to set up test app: %v", err)
	}
	defer testApp.TeardownTestApp()

	// Test unauthorized access (no token)
	t.Run("No Token", func(t *testing.T) {
		// Create a new request without token
		req, err := testApp.makeAuthenticatedRequest("POST", "/reimbursements", nil, "")
		require.NoError(t, err, "Failed to create reimbursement request")

		// Perform the request
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test reimbursement request")

		// Check the response status code - should be unauthorized
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode, "Expected status code to be 401 Unauthorized")
	})

	// Test unauthorized role (admin trying to create reimbursement)
	t.Run("Admin Creating Reimbursement", func(t *testing.T) {
		// Prepare reimbursement request
		reimbursementRequest := dto.CreateReimbursementBodyDto{
			Description: "Office supplies purchase",
			Amount:      150000,
		}

		// Convert to JSON
		requestBody, err := json.Marshal(reimbursementRequest)
		require.NoError(t, err, "Failed to marshal reimbursement request")

		// Create a new request with admin token
		req, err := testApp.makeAuthenticatedRequest("POST", "/reimbursements", requestBody, testApp.AdminToken)
		require.NoError(t, err, "Failed to create reimbursement request")
		req.Header.Set("Content-Type", "application/json")

		// Perform the request
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test reimbursement request")

		// Check the response status code - should be forbidden
		assert.Equal(t, fiber.StatusForbidden, resp.StatusCode, "Expected status code to be 403 Forbidden")
	})

	// Test unauthorized role (employee trying to approve reimbursement)
	t.Run("Employee Approving Reimbursement", func(t *testing.T) {
		// Create a test employee user
		salary := 4500000
		testUser := &entity.User{
			Username: "employee-approve-attempt-reimbursement",
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

		// Create a new request for reimbursement approval with employee token
		req, err := testApp.makeAuthenticatedRequest("POST", "/reimbursements/1/approve", nil, employeeToken)
		require.NoError(t, err, "Failed to create approval request")

		// Perform the request
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test approval request")

		// Check the response status code - should be forbidden
		assert.Equal(t, fiber.StatusForbidden, resp.StatusCode, "Expected status code to be 403 Forbidden")
	})
}

func TestNonExistentReimbursement(t *testing.T) {
	// Set up the test app
	testApp, err := SetupTestApp(t)
	if err != nil {
		t.Fatalf("Failed to set up test app: %v", err)
	}
	defer testApp.TeardownTestApp()

	// Test approving a non-existent reimbursement
	t.Run("Approve Non-Existent Reimbursement", func(t *testing.T) {
		// Create a new request for reimbursement approval with admin token
		// Use a very large ID that's unlikely to exist
		req, err := testApp.makeAuthenticatedRequest("POST", "/reimbursements/999999/approve", nil, testApp.AdminToken)
		require.NoError(t, err, "Failed to create approval request")

		// Perform the request
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test approval request")

		// Check the response status code - should be not found
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode, "Expected status code to be 404 Not Found")

		// Parse response body
		var response entity.HttpResponse
		body, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(body, &response)
		require.NoError(t, err, "Failed to parse response body")

		// Check response content
		assert.False(t, response.Success, "Expected success to be false")
		assert.Equal(t, "Reimbursement not found", response.Message, "Expected not found message")
	})
}
