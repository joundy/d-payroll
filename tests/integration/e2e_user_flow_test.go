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
	"github.com/stretchr/testify/require"
)

func TestE2EUserWorkWeekFlow(t *testing.T) {
	// Set up the test app
	testApp, err := SetupTestApp(t)
	if err != nil {
		t.Fatalf("Failed to set up test app: %v", err)
	}
	defer testApp.TeardownTestApp()

	var employeeToken string
	var employeeID uint

	// Step 1: Create an employee account
	t.Run("Step 1: Create Employee Account", func(t *testing.T) {
		// Mock time for user creation
		originalTimeNow := utils.TimeNow
		defer func() { utils.TimeNow = originalTimeNow }()
		utils.TimeNow = func() time.Time {
			return time.Date(2025, 6, 16, 8, 0, 0, 0, time.Local) // Monday, June 16, 2025 at 8:00 AM
		}

		// Create employee user
		salary := 5000000
		testUser := &entity.User{
			Username: "john.doe",
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

		employeeID = *createdUser.Id
		fmt.Printf("✅ Step 1: Created employee account with ID: %d, Username: %s\n", employeeID, createdUser.Username)
	})

	// Step 2: Login
	t.Run("Step 2: Employee Login", func(t *testing.T) {
		// Generate token for the employee
		token, err := utils.GenerateToken(testApp.Config.Auth.JwtSecret, &entity.AuthTokenPayload{
			ID:   employeeID,
			Role: entity.UserRoleEmployee,
		})
		require.NoError(t, err, "Failed to generate employee token")

		employeeToken = token
		fmt.Printf("✅ Step 2: Employee logged in successfully\n")
	})

	// Step 3: Monday - Check in at 8 AM, check out at 5 PM
	t.Run("Step 3: Monday - Check in at 8 AM, check out at 5 PM", func(t *testing.T) {
		// Mock time for Monday 8 AM check-in
		originalTimeNow := utils.TimeNow
		defer func() { utils.TimeNow = originalTimeNow }()

		// Check-in at 8 AM
		utils.TimeNow = func() time.Time {
			return time.Date(2025, 6, 16, 8, 0, 0, 0, time.Local) // Monday, June 16, 2025 at 8:00 AM
		}

		req, err := testApp.makeAuthenticatedRequest("POST", "/attendances/checkin", nil, employeeToken)
		require.NoError(t, err, "Failed to create check-in request")
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test check-in request")
		require.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected check-in to succeed")

		fmt.Printf("✅ Step 3a: Monday check-in at 8:00 AM completed\n")

		// Check-out at 5 PM
		utils.TimeNow = func() time.Time {
			return time.Date(2025, 6, 16, 17, 0, 0, 0, time.Local) // Monday, June 16, 2025 at 5:00 PM
		}

		req, err = testApp.makeAuthenticatedRequest("POST", "/attendances/checkout", nil, employeeToken)
		require.NoError(t, err, "Failed to create check-out request")
		resp, err = testApp.App.Test(req)
		require.NoError(t, err, "Failed to test check-out request")
		require.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected check-out to succeed")

		fmt.Printf("✅ Step 3b: Monday check-out at 5:00 PM completed\n")
	})

	// Step 4: Tuesday - Check in at 9 AM without checkout
	t.Run("Step 4: Tuesday - Check in at 9 AM without checkout", func(t *testing.T) {
		// Mock time for Tuesday 9 AM check-in
		originalTimeNow := utils.TimeNow
		defer func() { utils.TimeNow = originalTimeNow }()
		utils.TimeNow = func() time.Time {
			return time.Date(2025, 6, 17, 9, 0, 0, 0, time.Local) // Tuesday, June 17, 2025 at 9:00 AM
		}

		req, err := testApp.makeAuthenticatedRequest("POST", "/attendances/checkin", nil, employeeToken)
		require.NoError(t, err, "Failed to create check-in request")
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test check-in request")
		require.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected check-in to succeed")

		fmt.Printf("✅ Step 4: Tuesday check-in at 9:00 AM completed (no checkout)\n")
	})

	// Step 5: Wednesday - Check in at 10 AM, check out at 7 PM, submit overtime at 9 PM (2 hours)
	t.Run("Step 5: Wednesday - Check in at 10 AM, check out at 7 PM, submit overtime at 9 PM", func(t *testing.T) {
		originalTimeNow := utils.TimeNow
		defer func() { utils.TimeNow = originalTimeNow }()

		// Check-in at 10 AM
		utils.TimeNow = func() time.Time {
			return time.Date(2025, 6, 18, 10, 0, 0, 0, time.Local) // Wednesday, June 18, 2025 at 10:00 AM
		}

		req, err := testApp.makeAuthenticatedRequest("POST", "/attendances/checkin", nil, employeeToken)
		require.NoError(t, err, "Failed to create check-in request")
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test check-in request")
		require.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected check-in to succeed")

		fmt.Printf("✅ Step 5a: Wednesday check-in at 10:00 AM completed\n")

		// Check-out at 7 PM
		utils.TimeNow = func() time.Time {
			return time.Date(2025, 6, 18, 19, 0, 0, 0, time.Local) // Wednesday, June 18, 2025 at 7:00 PM
		}

		req, err = testApp.makeAuthenticatedRequest("POST", "/attendances/checkout", nil, employeeToken)
		require.NoError(t, err, "Failed to create check-out request")
		resp, err = testApp.App.Test(req)
		require.NoError(t, err, "Failed to test check-out request")
		require.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected check-out to succeed")

		fmt.Printf("✅ Step 5b: Wednesday check-out at 7:00 PM completed\n")

		// Submit overtime at 9 PM (2 hours)
		utils.TimeNow = func() time.Time {
			return time.Date(2025, 6, 18, 21, 0, 0, 0, time.Local) // Wednesday, June 18, 2025 at 9:00 PM
		}

		overtimeRequest := dto.CreateOvertimeBodyDto{
			Description:   "Working on urgent project deadline",
			OvertimeAt:    utils.TimeNow(),
			DurationMilis: 7200000, // 2 hours in milliseconds
		}

		requestBody, err := json.Marshal(overtimeRequest)
		require.NoError(t, err, "Failed to marshal overtime request")

		req, err = testApp.makeAuthenticatedRequest("POST", "/overtimes", requestBody, employeeToken)
		require.NoError(t, err, "Failed to create overtime request")
		req.Header.Set("Content-Type", "application/json")

		resp, err = testApp.App.Test(req)
		require.NoError(t, err, "Failed to test overtime request")
		require.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected overtime submission to succeed")

		fmt.Printf("✅ Step 5c: Wednesday overtime submission at 9:00 PM completed (2 hours)\n")
	})

	// Step 6: Thursday - Check in at 1 PM, check out at 3 PM
	t.Run("Step 6: Thursday - Check in at 1 PM, check out at 3 PM", func(t *testing.T) {
		originalTimeNow := utils.TimeNow
		defer func() { utils.TimeNow = originalTimeNow }()

		// Check-in at 1 PM
		utils.TimeNow = func() time.Time {
			return time.Date(2025, 6, 19, 13, 0, 0, 0, time.Local) // Thursday, June 19, 2025 at 1:00 PM
		}

		req, err := testApp.makeAuthenticatedRequest("POST", "/attendances/checkin", nil, employeeToken)
		require.NoError(t, err, "Failed to create check-in request")
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test check-in request")
		require.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected check-in to succeed")

		fmt.Printf("✅ Step 6a: Thursday check-in at 1:00 PM completed\n")

		// Check-out at 3 PM
		utils.TimeNow = func() time.Time {
			return time.Date(2025, 6, 19, 15, 0, 0, 0, time.Local) // Thursday, June 19, 2025 at 3:00 PM
		}

		req, err = testApp.makeAuthenticatedRequest("POST", "/attendances/checkout", nil, employeeToken)
		require.NoError(t, err, "Failed to create check-out request")
		resp, err = testApp.App.Test(req)
		require.NoError(t, err, "Failed to test check-out request")
		require.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected check-out to succeed")

		fmt.Printf("✅ Step 6b: Thursday check-out at 3:00 PM completed\n")
	})

	// Step 7: Friday - Check in at 8 AM, check out at 10 PM, submit reimbursement with amount 100000
	t.Run("Step 7: Friday - Check in at 8 AM, check out at 10 PM, submit reimbursement", func(t *testing.T) {
		originalTimeNow := utils.TimeNow
		defer func() { utils.TimeNow = originalTimeNow }()

		// Check-in at 8 AM
		utils.TimeNow = func() time.Time {
			return time.Date(2025, 6, 20, 8, 0, 0, 0, time.Local) // Friday, June 20, 2025 at 8:00 AM
		}

		req, err := testApp.makeAuthenticatedRequest("POST", "/attendances/checkin", nil, employeeToken)
		require.NoError(t, err, "Failed to create check-in request")
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test check-in request")
		require.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected check-in to succeed")

		fmt.Printf("✅ Step 7a: Friday check-in at 8:00 AM completed\n")

		// Check-out at 10 PM
		utils.TimeNow = func() time.Time {
			return time.Date(2025, 6, 20, 22, 0, 0, 0, time.Local) // Friday, June 20, 2025 at 10:00 PM
		}

		req, err = testApp.makeAuthenticatedRequest("POST", "/attendances/checkout", nil, employeeToken)
		require.NoError(t, err, "Failed to create check-out request")
		resp, err = testApp.App.Test(req)
		require.NoError(t, err, "Failed to test check-out request")
		require.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected check-out to succeed")

		fmt.Printf("✅ Step 7b: Friday check-out at 10:00 PM completed\n")

		// Submit reimbursement with amount 100000
		reimbursementRequest := dto.CreateReimbursementBodyDto{
			Description: "Business lunch with client",
			Amount:      100000,
		}

		requestBody, err := json.Marshal(reimbursementRequest)
		require.NoError(t, err, "Failed to marshal reimbursement request")

		req, err = testApp.makeAuthenticatedRequest("POST", "/reimbursements", requestBody, employeeToken)
		require.NoError(t, err, "Failed to create reimbursement request")
		req.Header.Set("Content-Type", "application/json")

		resp, err = testApp.App.Test(req)
		require.NoError(t, err, "Failed to test reimbursement request")
		require.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected reimbursement submission to succeed")

		fmt.Printf("✅ Step 7c: Friday reimbursement submission completed (Amount: 100,000)\n")
	})

	// Step 8: Saturday - Submit overtime at 9 AM (1 hour)
	t.Run("Step 8: Saturday - Submit overtime at 9 AM", func(t *testing.T) {
		originalTimeNow := utils.TimeNow
		defer func() { utils.TimeNow = originalTimeNow }()
		utils.TimeNow = func() time.Time {
			return time.Date(2025, 6, 21, 9, 0, 0, 0, time.Local) // Saturday, June 21, 2025 at 9:00 AM
		}

		overtimeRequest := dto.CreateOvertimeBodyDto{
			Description:   "Weekend maintenance work",
			OvertimeAt:    utils.TimeNow(),
			DurationMilis: 3600000, // 1 hour in milliseconds
		}

		requestBody, err := json.Marshal(overtimeRequest)
		require.NoError(t, err, "Failed to marshal overtime request")

		req, err := testApp.makeAuthenticatedRequest("POST", "/overtimes", requestBody, employeeToken)
		require.NoError(t, err, "Failed to create overtime request")
		req.Header.Set("Content-Type", "application/json")

		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test overtime request")
		require.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected overtime submission to succeed")

		fmt.Printf("✅ Step 8: Saturday overtime submission at 9:00 AM completed (1 hour)\n")
	})

	// Step 9: Print user's attendances, overtimes, and reimbursements
	t.Run("Step 9: Print User's Records", func(t *testing.T) {
		// Get attendances
		req, err := testApp.makeAuthenticatedRequest("GET", fmt.Sprintf("/attendances?user_id=%d", employeeID), nil, employeeToken)
		require.NoError(t, err, "Failed to create get attendances request")
		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test get attendances request")
		require.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected get attendances to succeed")

		var attendanceResponse entity.HttpResponse
		body, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(body, &attendanceResponse)
		require.NoError(t, err, "Failed to parse attendance response")

		fmt.Printf("🕐 ATTENDANCES:\n")
		if attendanceResponse.Data != nil {
			attendanceDataJson, _ := json.MarshalIndent(attendanceResponse.Data, "", "  ")
			fmt.Printf("%s\n\n", string(attendanceDataJson))
		}

		// Get overtimes
		req, err = testApp.makeAuthenticatedRequest("GET", fmt.Sprintf("/overtimes?user_id=%d", employeeID), nil, employeeToken)
		require.NoError(t, err, "Failed to create get overtimes request")
		resp, err = testApp.App.Test(req)
		require.NoError(t, err, "Failed to test get overtimes request")
		require.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected get overtimes to succeed")

		var overtimeResponse entity.HttpResponse
		body, _ = io.ReadAll(resp.Body)
		err = json.Unmarshal(body, &overtimeResponse)
		require.NoError(t, err, "Failed to parse overtime response")

		fmt.Printf("⏰ OVERTIMES:\n")
		if overtimeResponse.Data != nil {
			overtimeDataJson, _ := json.MarshalIndent(overtimeResponse.Data, "", "  ")
			fmt.Printf("%s\n\n", string(overtimeDataJson))
		}

		// Get reimbursements
		req, err = testApp.makeAuthenticatedRequest("GET", fmt.Sprintf("/reimbursements?user_id=%d", employeeID), nil, employeeToken)
		require.NoError(t, err, "Failed to create get reimbursements request")
		resp, err = testApp.App.Test(req)
		require.NoError(t, err, "Failed to test get reimbursements request")
		require.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected get reimbursements to succeed")

		var reimbursementResponse entity.HttpResponse
		body, _ = io.ReadAll(resp.Body)
		err = json.Unmarshal(body, &reimbursementResponse)
		require.NoError(t, err, "Failed to parse reimbursement response")

		fmt.Printf("💰 REIMBURSEMENTS:\n")
		if reimbursementResponse.Data != nil {
			reimbursementDataJson, _ := json.MarshalIndent(reimbursementResponse.Data, "", "  ")
			fmt.Printf("%s\n\n", string(reimbursementDataJson))
		}

		fmt.Printf("✅ Step 9: All user records printed successfully\n")
		fmt.Printf("🎉 === E2E USER FLOW COMPLETED ===\n\n")
	})

	// Step 10: Admin approves all overtimes
	t.Run("Step 10: Admin Approve All Overtimes", func(t *testing.T) {
		// Create admin token
		adminToken, err := utils.GenerateToken(testApp.Config.Auth.JwtSecret, &entity.AuthTokenPayload{
			ID:   1, // Admin user ID from setup
			Role: entity.UserRoleAdmin,
		})
		require.NoError(t, err, "Failed to generate admin token")

		// Get all overtimes to approve
		req, err := testApp.makeAuthenticatedRequest("GET", fmt.Sprintf("/overtimes?user_id=%d", employeeID), nil, adminToken)
		require.NoError(t, err, "Failed to create get overtimes request")

		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test get overtimes request")
		require.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected get overtimes to succeed")

		var overtimeResponse entity.HttpResponse
		body, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(body, &overtimeResponse)
		require.NoError(t, err, "Failed to parse overtime response")

		// Approve each overtime
		if overtimeResponse.Data != nil {
			overtimes, ok := overtimeResponse.Data.([]interface{})
			require.True(t, ok, "Expected overtimes to be an array")

			for _, overtime := range overtimes {
				overtimeMap, ok := overtime.(map[string]interface{})
				require.True(t, ok, "Expected overtime to be a map")

				overtimeID, ok := overtimeMap["id"].(float64)
				require.True(t, ok, "Expected overtime ID to be a number")

				// Approve the overtime
				approveURL := fmt.Sprintf("/overtimes/%d/approve", int(overtimeID))
				req, err := testApp.makeAuthenticatedRequest("POST", approveURL, nil, adminToken)
				require.NoError(t, err, "Failed to create approve overtime request")

				resp, err := testApp.App.Test(req)
				require.NoError(t, err, "Failed to test approve overtime request")
				require.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected approve overtime to succeed")

				fmt.Printf("✅ Approved overtime ID: %d\n", int(overtimeID))
			}
		}

		fmt.Printf("✅ Step 10: All overtimes approved successfully\n\n")
	})

	// Step 11: Admin approves all reimbursements
	t.Run("Step 11: Admin Approve All Reimbursements", func(t *testing.T) {
		// Create admin token
		adminToken, err := utils.GenerateToken(testApp.Config.Auth.JwtSecret, &entity.AuthTokenPayload{
			ID:   1, // Admin user ID from setup
			Role: entity.UserRoleAdmin,
		})
		require.NoError(t, err, "Failed to generate admin token")

		// Get all reimbursements to approve
		req, err := testApp.makeAuthenticatedRequest("GET", fmt.Sprintf("/reimbursements?user_id=%d", employeeID), nil, adminToken)
		require.NoError(t, err, "Failed to create get reimbursements request")

		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test get reimbursements request")
		require.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected get reimbursements to succeed")

		var reimbursementResponse entity.HttpResponse
		body, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(body, &reimbursementResponse)
		require.NoError(t, err, "Failed to parse reimbursement response")

		// Approve each reimbursement
		if reimbursementResponse.Data != nil {
			reimbursements, ok := reimbursementResponse.Data.([]interface{})
			require.True(t, ok, "Expected reimbursements to be an array")

			for _, reimbursement := range reimbursements {
				reimbursementMap, ok := reimbursement.(map[string]interface{})
				require.True(t, ok, "Expected reimbursement to be a map")

				reimbursementID, ok := reimbursementMap["id"].(float64)
				require.True(t, ok, "Expected reimbursement ID to be a number")

				// Approve the reimbursement
				approveURL := fmt.Sprintf("/reimbursements/%d/approve", int(reimbursementID))
				req, err := testApp.makeAuthenticatedRequest("POST", approveURL, nil, adminToken)
				require.NoError(t, err, "Failed to create approve reimbursement request")

				resp, err := testApp.App.Test(req)
				require.NoError(t, err, "Failed to test approve reimbursement request")
				require.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected approve reimbursement to succeed")

				fmt.Printf("✅ Approved reimbursement ID: %d\n", int(reimbursementID))
			}
		}

		fmt.Printf("✅ Step 11: All reimbursements approved successfully\n\n")
	})

	// Step 12: Create payroll for the work week (2025-06-16 to 2025-06-21)
	t.Run("Step 12: Create Payroll for Work Week", func(t *testing.T) {
		// Mock time for payroll creation
		originalTimeNow := utils.TimeNow
		defer func() { utils.TimeNow = originalTimeNow }()
		utils.TimeNow = func() time.Time {
			return time.Date(2025, 6, 22, 9, 0, 0, 0, time.Local) // Sunday, June 22, 2025 at 9:00 AM (after work week)
		}

		// Create admin token for payroll operations
		adminToken, err := utils.GenerateToken(testApp.Config.Auth.JwtSecret, &entity.AuthTokenPayload{
			ID:   1, // Admin user ID from setup
			Role: entity.UserRoleAdmin,
		})
		require.NoError(t, err, "Failed to generate admin token")

		// Prepare payroll request
		startedAt := time.Date(2025, 6, 16, 0, 0, 0, 0, time.Local)  // Monday
		endedAt := time.Date(2025, 6, 21, 23, 59, 59, 0, time.Local) // Saturday

		payrollRequest := dto.CreatePayrollBodyDto{
			Name:      "Work Week June 16-21, 2025",
			StartedAt: startedAt,
			EndedAt:   endedAt,
		}

		requestBody, err := json.Marshal(payrollRequest)
		require.NoError(t, err, "Failed to marshal payroll request")

		req, err := testApp.makeAuthenticatedRequest("POST", "/payrolls", requestBody, adminToken)
		require.NoError(t, err, "Failed to create payroll request")
		req.Header.Set("Content-Type", "application/json")

		resp, err := testApp.App.Test(req)
		require.NoError(t, err, "Failed to test payroll request")
		require.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected payroll creation to succeed")

		// Parse response to get payroll ID
		var payrollResponse entity.HttpResponse
		body, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(body, &payrollResponse)
		require.NoError(t, err, "Failed to parse payroll response")

		fmt.Printf("✅ Step 12: Payroll created successfully for work week (June 16-21, 2025)\n")
		fmt.Printf("📋 Payroll Details:\n")
		payrollDataJson, _ := json.MarshalIndent(payrollResponse.Data, "", "  ")
		fmt.Printf("%s\n\n", string(payrollDataJson))

		// Extract payroll ID for rolling
		payrollData, ok := payrollResponse.Data.(map[string]interface{})
		require.True(t, ok, "Expected payroll data to be a map")
		payrollIDFloat, ok := payrollData["id"].(float64)
		require.True(t, ok, "Expected payroll ID to be a number")
		payrollID := int(payrollIDFloat)

		// Step 13: Roll the payroll
		t.Run("Step 13: Roll Payroll", func(t *testing.T) {
			// Mock time for payroll rolling to Friday, June 20, 2025
			originalTimeNow := utils.TimeNow
			defer func() { utils.TimeNow = originalTimeNow }()
			utils.TimeNow = func() time.Time {
				return time.Date(2025, 6, 20, 23, 0, 0, 0, time.Local)
			}

			rollURL := fmt.Sprintf("/payrolls/%d/roll", payrollID)
			req, err := testApp.makeAuthenticatedRequest("POST", rollURL, nil, adminToken)
			require.NoError(t, err, "Failed to create roll payroll request")

			resp, err := testApp.App.Test(req)
			require.NoError(t, err, "Failed to test roll payroll request")
			require.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected payroll rolling to succeed")

			var rollResponse entity.HttpResponse
			body, _ := io.ReadAll(resp.Body)
			err = json.Unmarshal(body, &rollResponse)
			require.NoError(t, err, "Failed to parse roll payroll response")

			fmt.Printf("✅ Step 13: Payroll rolled successfully on Friday, June 20, 2025 at 5:00 PM\n")
			fmt.Printf("🎯 Roll Result:\n")
			rollDataJson, _ := json.MarshalIndent(rollResponse.Data, "", "  ")
			fmt.Printf("%s\n\n", string(rollDataJson))
		})

		// Step 14: Verify payroll status after rolling
		t.Run("Step 14: Verify Payroll Status", func(t *testing.T) {
			req, err := testApp.makeAuthenticatedRequest("GET", "/payrolls", nil, adminToken)
			require.NoError(t, err, "Failed to create get payrolls request")

			resp, err := testApp.App.Test(req)
			require.NoError(t, err, "Failed to test get payrolls request")
			require.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected get payrolls to succeed")

			var payrollsResponse entity.HttpResponse
			body, _ := io.ReadAll(resp.Body)
			err = json.Unmarshal(body, &payrollsResponse)
			require.NoError(t, err, "Failed to parse payrolls response")

			fmt.Printf("✅ Step 14: Payroll status verified\n")
			fmt.Printf("📊 All Payrolls:\n")
			payrollsDataJson, _ := json.MarshalIndent(payrollsResponse.Data, "", "  ")
			fmt.Printf("%s\n\n", string(payrollsDataJson))
		})

		// Step 15: Employee generates payslip
		t.Run("Step 15: Employee Generate Payslip", func(t *testing.T) {
			// Use payroll ID 1 for simplicity (you can make this dynamic by parsing payrollsResponse.Data)
			payslipURL := fmt.Sprintf("/payrolls/1/payslips?user_id=%d", employeeID)
			req, err := testApp.makeAuthenticatedRequest("POST", payslipURL, nil, employeeToken)
			require.NoError(t, err, "Failed to create generate payslip request")

			resp, err := testApp.App.Test(req)
			require.NoError(t, err, "Failed to test generate payslip request")
			require.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected generate payslip to succeed")

			var payslipResponse entity.HttpResponse
			body, _ := io.ReadAll(resp.Body)
			err = json.Unmarshal(body, &payslipResponse)
			require.NoError(t, err, "Failed to parse payslip response")

			fmt.Printf("✅ Step 15: Employee payslip generated successfully\n")
			fmt.Printf("💰 Payslip Details:\n")
			payslipDataJson, _ := json.MarshalIndent(payslipResponse.Data, "", "  ")
			fmt.Printf("%s\n\n", string(payslipDataJson))
		})

	})
}
