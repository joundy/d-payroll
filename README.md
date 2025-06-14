# D Payroll System

## Overview

The Dealls Payroll System is a backend application designed to manage employee payroll, including attendance, overtime, reimbursements, and salary calculations. It provides a RESTful API for various operations.

## Features

*   User Management (Admin and Employee roles)
*   Authentication (JWT-based)
*   Attendance Tracking (Check-in/Check-out)
*   Overtime Request and Approval
*   Reimbursement Request and Approval
*   Automated Payroll Processing
*   Payslip Generation

## Tech Stack

*   **Language:** Go
*   **Framework:** Fiber (for HTTP server)
*   **Database:** PostgreSQL
*   **ORM:** GORM
*   **Authentication:** JWT

## Project Structure

The project follows a standard Go project layout, emphasizing a clean separation of concerns:

```text
/
├── cmd/                     # Main applications (entry points)
│   └── app/                 # Main application server (main.go)
├── config/                  # Configuration loading (e.g., from .env)
├── controller/              # HTTP request handlers and input/output structuring
│   ├── http/                # HTTP specific controllers for API routes
│   │   └── dto/             # Data Transfer Objects for request/response bodies
├── db/                      # Database related files
│   └── migrations/          # Database migration scripts (using go-migrate)
├── entity/                  # Core domain models/structs (e.g., User, Payroll)
├── internal-error/          # Custom error types for the application
├── repository/              # Data access layer (interacts with the database via GORM)
├── service/                 # Business logic layer, orchestrates operations
├── tests/                   # Test files
│   └── integration/         # Integration tests (using testcontainers)
├── utils/                   # Utility/helper functions
├── .env                     # Environment variables (local, gitignored)
├── .env.example             # Example environment variables template
├── go.mod                   # Go module definitions and dependencies
├── go.sum                   # Go module checksums
├── Makefile                 # Make commands for common tasks (build, run, test)
├── docker-compose.yaml      # Docker Compose for development services (e.g., DB)
└── README.md                # This file
```

## Prerequisites

*   Go (version 1.24.4 or higher)
*   PostgreSQL (version 12 or higher recommended, please use your installed version)
*   Docker (optional, for running PostgreSQL)
*   Git

## Installation

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/joundy/d-payroll.git
    cd dealls-payroll
    ```

2.  **Set up environment variables:**
    Create a `.env` file in the root directory by copying `.env.example` (if it exists) and fill in the necessary configuration (e.g., database credentials, JWT secret).
    ```bash
    cp .env.example .env 
    # Edit .env with your configuration
    ```
    *Note: If `.env.example` is not provided, you'll need to create `.env` manually with required variables like `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `JWT_SECRET`.*

3.  **Install dependencies:**
    ```bash
    go mod tidy
    ```

4.  **Database Setup:**
    Ensure your PostgreSQL server is running and accessible. Create the database specified in your `.env` file.
    Run database migrations:
    The project uses `go-migrate` for managing database schema changes. Ensure you have `go-migrate` installed (see [go-migrate installation](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#installation)).
    The migration files are located in the `db/migrations` directory.

    To apply all pending migrations (UP):
    ```bash
    migrate -source file://db/migrations -database "postgresql://YOUR_USER:YOUR_PASSWORD@YOUR_HOST:YOUR_PORT/YOUR_DB_NAME?sslmode=disable" up
    ```
    *Replace `YOUR_USER`, `YOUR_PASSWORD`, etc., with your actual database connection details from the `.env` file.*

    To roll back the last applied migration (DOWN):
    ```bash
    migrate -source file://db/migrations -database "postgresql://YOUR_USER:YOUR_PASSWORD@YOUR_HOST:YOUR_PORT/YOUR_DB_NAME?sslmode=disable" down 1
    ```

## Running the Application

1.  **Start the server:**
    ```bash
    go run cmd/app/main.go
    ```
    Alternatively, this project includes a `Makefile` which may provide convenient commands for building and running the application (e.g., `make run` or `make build`). Please refer to the `Makefile` for available targets.

    The application should now be running on the configured port (e.g., `http://localhost:8080`).

## Testing

The project includes integration tests located in the `test/integration` directory. These tests can be used to verify the functionality of various API endpoints and demonstrate happy path flows.

**Important:** The integration tests utilize test containers (e.g., via a library like `testcontainers-go`) to spin up a dedicated test database instance. Therefore, **Docker must be installed and running** on your system to execute these tests successfully.

To run the tests, you might use standard `go test` commands targeting this directory, or check the `Makefile` for specific test execution targets (e.g., `make test` or `make test-integration`).

Example (navigate to the project root):
```bash
# (Ensure Docker is running)
# go test ./tests/integration
```

## API Documentation

This section details the available API endpoints, request formats, and response examples.

### Authentication

#### Login

*   **Endpoint:** `POST /login`
*   **Description:** Authenticates a user and returns a JWT token.
*   **Request Body:** `application/json`
    ```json
    {
        "username": "your_username",
        "password": "your_password"
    }
    ```
*   **Response (Success 200 OK):** `application/json`
    ```json
    {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
    ```
*   **Responses (Error):**
    *   `400 Bad Request`: Invalid request body.
        ```json
        {
            "message": "Invalid request body"
        }
        ```
    *   `404 Not Found`: User not found.
        ```json
        {
            "message": "User not found"
        }
        ```
    *   `401 Unauthorized`: Invalid credentials.
        ```json
        {
            "message": "Invalid credentials"
        }
        ```

### User Management

All User Management endpoints require Admin privileges.

#### Create User

*   **Endpoint:** `POST /users`
*   **Description:** Creates a new user (Admin or Employee).
*   **Authentication:** Required (Admin role).
*   **Request Body:** `application/json`
    ```json
    {
        "username": "newuser",
        "password": "securepassword123",
        "role": "EMPLOYEE", // or "ADMIN"
        "user_info": {
            "monthly_salary": 5000000 
        }
    }
    ```
    *Note: `user_info` and `monthly_salary` are optional for an Admin user but generally required for an Employee if salary is managed.*
*   **Response (Success 200 OK):** `application/json`
    ```json
    {
        "id": 123,
        "username": "newuser",
        "role": "EMPLOYEE",
        "user_info": {
            "monthly_salary": 5000000
        },
        "created_at": "2023-10-27T10:00:00Z",
        "updated_at": "2023-10-27T10:00:00Z"
    }
    ```
*   **Responses (Error):**
    *   `400 Bad Request`: Invalid request body or validation error.
    *   `401 Unauthorized`: Missing or invalid token.
    *   `403 Forbidden`: User does not have Admin privileges.

#### Get User by ID

*   **Endpoint:** `GET /users/:id`
*   **Description:** Retrieves a specific user by their ID.
*   **Authentication:** Required (Admin role).
*   **Path Parameters:**
    *   `id` (integer, required): The ID of the user to retrieve.
*   **Response (Success 200 OK):** `application/json`
    ```json
    {
        "id": 123,
        "username": "existinguser",
        "role": "EMPLOYEE",
        "user_info": {
            "monthly_salary": 6000000
        },
        "created_at": "2023-01-15T09:30:00Z",
        "updated_at": "2023-05-20T14:45:00Z"
    }
    ```
*   **Responses (Error):**
    *   `400 Bad Request`: Invalid ID parameter.
    *   `401 Unauthorized`: Missing or invalid token.
    *   `403 Forbidden`: User does not have Admin privileges.
    *   `404 Not Found`: User with the specified ID not found.

### Attendance Management

#### Check-in

*   **Endpoint:** `POST /attendances/checkin`
*   **Description:** Allows an authenticated employee to record their check-in time.
*   **Authentication:** Required (Employee role).
*   **Request Body:** None.
*   **Response (Success 200 OK):** `application/json`
    ```json
    {
        "id": 1,
        "type": "CHECK_IN",
        "created_at": "2023-10-27T09:00:00Z",
        "updated_at": "2023-10-27T09:00:00Z"
    }
    ```
*   **Responses (Error):**
    *   `401 Unauthorized`: Missing or invalid token.
    *   `403 Forbidden`: User does not have Employee privileges.
    *   `409 Conflict`: "User already checked in".
    *   `422 Unprocessable Entity`: "User cannot checked in on weekend".

#### Check-out

*   **Endpoint:** `POST /attendances/checkout`
*   **Description:** Allows an authenticated employee to record their check-out time.
*   **Authentication:** Required (Employee role).
*   **Request Body:** None.
*   **Response (Success 200 OK):** `application/json`
    ```json
    {
        "id": 2,
        "type": "CHECK_OUT",
        "created_at": "2023-10-27T17:30:00Z",
        "updated_at": "2023-10-27T17:30:00Z"
    }
    ```
*   **Responses (Error):**
    *   `401 Unauthorized`: Missing or invalid token.
    *   `403 Forbidden`: User does not have Employee privileges.
    *   `409 Conflict`: "User already checked out".
    *   `422 Unprocessable Entity`: "User cannot checked out because it is not checked in".

#### Get Attendances by User ID

*   **Endpoint:** `GET /attendances`
*   **Description:** Retrieves a list of attendance records for a specified user. Employees can only fetch their own records. Admins can fetch records for any user.
*   **Authentication:** Required (Employee or Admin role).
*   **Query Parameters:**
    *   `user_id` (integer, required): The ID of the user whose attendances are to be fetched.
*   **Response (Success 200 OK):** `application/json`
    ```json
    [
        {
            "id": 1,
            "type": "CHECK_IN",
            "created_at": "2023-10-27T09:00:00Z",
            "updated_at": "2023-10-27T09:00:00Z"
        },
        {
            "id": 2,
            "type": "CHECK_OUT",
            "created_at": "2023-10-27T17:30:00Z",
            "updated_at": "2023-10-27T17:30:00Z"
        }
        // ... more attendance records
    ]
    ```
*   **Responses (Error):**
    *   `400 Bad Request`: "Invalid user ID query".
    *   `401 Unauthorized`: Missing or invalid token, or Employee attempting to access another user's data.
    *   `403 Forbidden`: User does not have sufficient privileges.

### Overtime Management

#### Submit Overtime Request

*   **Endpoint:** `POST /overtimes`
*   **Description:** Allows an authenticated employee to submit an overtime request.
*   **Authentication:** Required (Employee role).
*   **Request Body:** `application/json`
    ```json
    {
        "description": "Urgent bug fix for production issue",
        "overtime_at": "2023-10-27T18:00:00Z", // Date and start time of overtime
        "duration_milis": 7200000 // Duration in milliseconds (e.g., 2 hours)
    }
    ```
*   **Response (Success 200 OK):** `application/json`
    ```json
    {
        "id": 101,
        "user_id": 45,
        "description": "Urgent bug fix for production issue",
        "overtime_at": "2023-10-27T18:00:00Z",
        "duration_milis": 7200000,
        "is_approved": false,
        "updated_by_user_id": null,
        "created_at": "2023-10-27T17:35:00Z",
        "updated_at": "2023-10-27T17:35:00Z"
    }
    ```
*   **Responses (Error):**
    *   `400 Bad Request`: Invalid request body.
    *   `401 Unauthorized`: Missing or invalid token.
    *   `403 Forbidden`: User does not have Employee privileges.
    *   `422 Unprocessable Entity`: "Overtime exceeds limit" or "Overtime submit before checkout".

#### Approve Overtime Request

*   **Endpoint:** `POST /overtimes/:overtimeId/approve`
*   **Description:** Allows an authenticated admin to approve a pending overtime request.
*   **Authentication:** Required (Admin role).
*   **Path Parameters:**
    *   `overtimeId` (integer, required): The ID of the overtime request to approve.
*   **Request Body:** None.
*   **Response (Success 200 OK):** No content.
*   **Responses (Error):**
    *   `400 Bad Request`: "Invalid overtime ID param".
    *   `401 Unauthorized`: Missing or invalid token.
    *   `403 Forbidden`: User does not have Admin privileges.
    *   `404 Not Found`: "Overtime not found".
    *   `409 Conflict`: "Overtime already approved".

#### Get User Overtime Requests

*   **Endpoint:** `GET /overtimes`
*   **Description:** Retrieves a list of overtime requests for a specified user. Employees can only fetch their own records. Admins can fetch records for any user.
*   **Authentication:** Required (Employee or Admin role).
*   **Query Parameters:**
    *   `user_id` (integer, required): The ID of the user whose overtime requests are to be fetched.
*   **Response (Success 200 OK):** `application/json`
    ```json
    [
        {
            "id": 101,
            "user_id": 45,
            "description": "Urgent bug fix for production issue",
            "overtime_at": "2023-10-27T18:00:00Z",
            "duration_milis": 7200000,
            "is_approved": true,
            "updated_by_user_id": 10, // Admin user ID who approved
            "created_at": "2023-10-27T17:35:00Z",
            "updated_at": "2023-10-27T19:05:00Z"
        },
        {
            "id": 102,
            "user_id": 45,
            "description": "Completing quarterly report",
            "overtime_at": "2023-10-28T19:00:00Z",
            "duration_milis": 3600000,
            "is_approved": false,
            "updated_by_user_id": null,
            "created_at": "2023-10-28T10:00:00Z",
            "updated_at": "2023-10-28T10:00:00Z"
        }
        // ... more overtime records
    ]
    ```
*   **Responses (Error):**
    *   `400 Bad Request`: "Invalid user ID query".
    *   `401 Unauthorized`: Missing or invalid token, or Employee attempting to access another user's data.
    *   `403 Forbidden`: User does not have sufficient privileges.

### Reimbursement Management

#### Submit Reimbursement Request

*   **Endpoint:** `POST /reimbursements`
*   **Description:** Allows an authenticated employee to submit a reimbursement request.
*   **Authentication:** Required (Employee role).
*   **Request Body:** `application/json`
    ```json
    {
        "description": "Client meeting transportation costs",
        "amount": 150000 // Amount in the smallest currency unit (e.g., cents, or full units if not using decimals)
    }
    ```
*   **Response (Success 200 OK):** `application/json`
    ```json
    {
        "id": 201,
        "user_id": 45,
        "description": "Client meeting transportation costs",
        "amount": 150000,
        "is_approved": false,
        "updated_by_user_id": null,
        "created_at": "2023-10-29T11:00:00Z",
        "updated_at": "2023-10-29T11:00:00Z"
    }
    ```
*   **Responses (Error):**
    *   `400 Bad Request`: Invalid request body.
    *   `401 Unauthorized`: Missing or invalid token.
    *   `403 Forbidden`: User does not have Employee privileges.

#### Approve Reimbursement Request

*   **Endpoint:** `POST /reimbursements/:reimbursementId/approve`
*   **Description:** Allows an authenticated admin to approve a pending reimbursement request.
*   **Authentication:** Required (Admin role).
*   **Path Parameters:**
    *   `reimbursementId` (integer, required): The ID of the reimbursement request to approve.
*   **Request Body:** None.
*   **Response (Success 200 OK):** No content.
*   **Responses (Error):**
    *   `400 Bad Request`: "Invalid reimbursement ID param".
    *   `401 Unauthorized`: Missing or invalid token.
    *   `403 Forbidden`: User does not have Admin privileges.
    *   `404 Not Found`: "Reimbursement not found".
    *   `409 Conflict`: "Reimbursement already approved".

#### Get User Reimbursement Requests

*   **Endpoint:** `GET /reimbursements`
*   **Description:** Retrieves a list of reimbursement requests for a specified user. Employees can only fetch their own records. Admins can fetch records for any user.
*   **Authentication:** Required (Employee or Admin role).
*   **Query Parameters:**
    *   `user_id` (integer, required): The ID of the user whose reimbursement requests are to be fetched.
*   **Response (Success 200 OK):** `application/json`
    ```json
    [
        {
            "id": 201,
            "user_id": 45,
            "description": "Client meeting transportation costs",
            "amount": 150000,
            "is_approved": true,
            "updated_by_user_id": 10, // Admin user ID who approved
            "created_at": "2023-10-29T11:00:00Z",
            "updated_at": "2023-10-29T14:30:00Z"
        },
        {
            "id": 202,
            "user_id": 45,
            "description": "Software license purchase",
            "amount": 500000,
            "is_approved": false,
            "updated_by_user_id": null,
            "created_at": "2023-10-30T09:15:00Z",
            "updated_at": "2023-10-30T09:15:00Z"
        }
        // ... more reimbursement records
    ]
    ```
*   **Responses (Error):**
    *   `400 Bad Request`: "Invalid user ID query".
    *   `401 Unauthorized`: Missing or invalid token, or Employee attempting to access another user's data.
    *   `403 Forbidden`: User does not have sufficient privileges.

### Payroll Management

All Payroll Management endpoints require Admin privileges, except for fetching one's own payslip.

#### Create Payroll Period

*   **Endpoint:** `POST /payrolls`
*   **Description:** Creates a new payroll period for processing.
*   **Authentication:** Required (Admin role).
*   **Request Body:** `application/json`
    ```json
    {
        "name": "November 2023 Payroll",
        "started_at": "2023-11-01T00:00:00Z",
        "ended_at": "2023-11-30T23:59:59Z"
    }
    ```
*   **Response (Success 200 OK):** `application/json`
    ```json
    {
        "id": 301,
        "name": "November 2023 Payroll",
        "started_at": "2023-11-01T00:00:00Z",
        "ended_at": "2023-11-30T23:59:59Z",
        "is_rolled": false,
        "updated_by_user_id": null,
        "created_by_user_id": 1, // Admin user ID who created
        "created_at": "2023-10-27T10:00:00Z",
        "updated_at": "2023-10-27T10:00:00Z"
    }
    ```
*   **Responses (Error):**
    *   `400 Bad Request`: Invalid request body.
    *   `401 Unauthorized`: Missing or invalid token.
    *   `403 Forbidden`: User does not have Admin privileges.

#### Get All Payroll Periods

*   **Endpoint:** `GET /payrolls`
*   **Description:** Retrieves a list of all payroll periods.
*   **Authentication:** Required (Admin role).
*   **Response (Success 200 OK):** `application/json`
    ```json
    [
        {
            "id": 301,
            "name": "November 2023 Payroll",
            "started_at": "2023-11-01T00:00:00Z",
            "ended_at": "2023-11-30T23:59:59Z",
            "is_rolled": false,
            "updated_by_user_id": null,
            "created_by_user_id": 1,
            "created_at": "2023-10-27T10:00:00Z",
            "updated_at": "2023-10-27T10:00:00Z"
        },
        {
            "id": 300,
            "name": "October 2023 Payroll",
            "started_at": "2023-10-01T00:00:00Z",
            "ended_at": "2023-10-31T23:59:59Z",
            "is_rolled": true,
            "updated_by_user_id": 2, // Admin user ID who rolled
            "created_by_user_id": 1,
            "created_at": "2023-09-27T10:00:00Z",
            "updated_at": "2023-11-05T11:00:00Z"
        }
        // ... more payroll periods
    ]
    ```
*   **Responses (Error):**
    *   `401 Unauthorized`: Missing or invalid token.
    *   `403 Forbidden`: User does not have Admin privileges.

#### Roll Payroll Period

*   **Endpoint:** `POST /payrolls/:payrollId/roll`
*   **Description:** Finalizes a payroll period, calculating all payslips. This action is irreversible for the given payroll period.
*   **Authentication:** Required (Admin role).
*   **Path Parameters:**
    *   `payrollId` (integer, required): The ID of the payroll period to roll.
*   **Request Body:** None.
*   **Response (Success 200 OK):** No content.
*   **Responses (Error):**
    *   `400 Bad Request`: "Invalid payroll ID param".
    *   `401 Unauthorized`: Missing or invalid token.
    *   `403 Forbidden`: User does not have Admin privileges.
    *   `404 Not Found`: "Payroll not found".
    *   `409 Conflict`: "Payroll already rolled".

#### Get User Payslip

*   **Endpoint:** `POST /payrolls/:payrollId/payslips`
*   **Description:** Retrieves the payslip for a specific user within a rolled payroll period. Employees can only fetch their own payslips. Admins can fetch for any user.
*   **Authentication:** Required (Employee or Admin role).
*   **Path Parameters:**
    *   `payrollId` (integer, required): The ID of the rolled payroll period.
*   **Query Parameters:**
    *   `user_id` (integer, required): The ID of the user whose payslip is to be fetched.
*   **Response (Success 200 OK):** `application/json`
    ```json
    {
        "payroll_id": 300,
        "user_id": 45,
        "salary": 5000000,
        "pro_rate": 1.0,
        "attendance": {
            "details": [
                {
                    "checkin_at": "2023-10-02T09:00:00Z",
                    "checkout_at": "2023-10-02T17:30:00Z",
                    "duration_milis": 30600000
                }
                // ... more attendance details
            ],
            "total_duration_milis": 612000000, // Example total for the period
            "total_amount": 5000000 
        },
        "overtime": {
            "details": [
                {
                    "overtime_at": "2023-10-05T18:00:00Z",
                    "description": "Urgent fix",
                    "duration_milis": 7200000,
                    "created_at": "2023-10-05T17:00:00Z"
                }
                // ... more overtime details
            ],
            "total_duration_milis": 14400000,
            "total_amount": 250000 
        },
        "reimburse": {
            "details": [
                {
                    "description": "Transport for client meeting",
                    "amount": 50000,
                    "created_at": "2023-10-10T10:00:00Z"
                }
                // ... more reimbursement details
            ],
            "total_amount": 50000
        },
        "take_home_pay": 5300000 
    }
    ```
*   **Responses (Error):**
    *   `400 Bad Request`: "Invalid payroll ID param" or "Invalid user ID query".
    *   `401 Unauthorized`: Missing or invalid token, or Employee attempting to access another user's payslip.
    *   `403 Forbidden`: User does not have sufficient privileges.
    *   `422 Unprocessable Entity`: "Payroll is not rolled yet".

#### Get Payslip Summaries for Payroll Period

*   **Endpoint:** `POST /payrolls/:payrollId/payslip-summaries`
*   **Description:** Retrieves a summary of payslips (user ID and total take-home pay) for all users in a rolled payroll period.
*   **Authentication:** Required (Admin role).
*   **Path Parameters:**
    *   `payrollId` (integer, required): The ID of the rolled payroll period.
*   **Response (Success 200 OK):** `application/json`
    ```json
    [
        {
            "payroll_id": 300,
            "user_id": 45,
            "total_take_home_pay": 5300000,
            "created_at": "2023-11-05T11:00:00Z",
            "updated_at": "2023-11-05T11:00:00Z"
        },
        {
            "payroll_id": 300,
            "user_id": 46,
            "total_take_home_pay": 6250000,
            "created_at": "2023-11-05T11:00:00Z",
            "updated_at": "2023-11-05T11:00:00Z"
        }
        // ... more user payslip summaries
    ]
    ```
*   **Responses (Error):**
    *   `400 Bad Request`: "Invalid payroll ID param".
    *   `401 Unauthorized`: Missing or invalid token.
    *   `403 Forbidden`: User does not have Admin privileges.
    *   `422 Unprocessable Entity`: "Payroll is not rolled yet".

#### Get Total Take-Home Pay for Payroll Period

*   **Endpoint:** `POST /payrolls/:payrollId/total-take-home-pay`
*   **Description:** Calculates and retrieves the sum of all take-home pay for a rolled payroll period.
*   **Authentication:** Required (Admin role).
*   **Path Parameters:**
    *   `payrollId` (integer, required): The ID of the rolled payroll period.
*   **Response (Success 200 OK):** `application/json`
    ```json
    {
        "data": 11550000 // Example sum for all users
    }
    ```
    *(Note: The exact structure of this response might vary; the example above assumes a simple object. The actual service returns an integer directly, which Fiber wraps in JSON)*
*   **Responses (Error):**
    *   `400 Bad Request`: "Invalid payroll ID param".
    *   `401 Unauthorized`: Missing or invalid token.
    *   `403 Forbidden`: User does not have Admin privileges.
    *   `422 Unprocessable Entity`: "Payroll is not rolled yet".

---

## Important Notes & Future Improvements

This project provides a foundational payroll system. While it covers core functionalities, there are several areas for improvement and consideration for future development:

1.  **Concurrency and Locking:**
    *   To optimize performance and prevent race conditions, especially during payroll processing or when multiple admins might perform overlapping actions, implementing robust locking mechanisms is crucial. This could involve using mutexes at the application level or leveraging PostgreSQL's transaction locks (e.g., `FOR UPDATE`, `SHARE` locks).

2.  **Scalability of Payslip Generation/Summary:**
    *   The current payslip summary generation might become a bottleneck with a large number of users. For better scalability, consider:
        *   Implementing a worker pool to process payslips in parallel.
        *   Using a message queue (e.g., RabbitMQ, Kafka) to offload payslip generation and summary calculations to background processes. This would make the API response faster and the system more resilient.

3.  **Caching Strategies:**
    *   Certain data, like generated user payslips or frequently accessed user information, could be cached (e.g., using Redis or an in-memory cache) to improve read performance and reduce database load.

4.  **Edge Case Handling:**
    *   There are numerous edge cases to consider for a production-grade payroll system:
        *   **Attendance:** How to handle scenarios where an employee checks in but forgets to check out?
        *   **Payroll Period Overlaps:** What if a payroll period starts or ends in the middle of an employee's active session or attendance record?
        *   **Prorated Salaries:** The basis for proration needs clear definition (e.g., based on a fixed number of workdays like 20 or 22 per month, or actual calendar days).
        *   Employee onboarding/offboarding mid-period.

5.  **Timezone Handling:**
    *   Currently, the system likely relies on the server's local timezone. For consistency and to support distributed teams, all timestamps should ideally be stored and processed in UTC. Unix timestamps can also be an alternative. Display layers can then convert UTC to the user's local timezone.

6.  **General Refinements:**
    *   This initial version focuses on core functionality. Further refactoring, more comprehensive error handling, and additional validation are areas for ongoing improvement.

These notes serve as a reminder for future iterations and to highlight areas where the system can be made more robust, scalable, and feature-rich.

This concludes the API documentation.