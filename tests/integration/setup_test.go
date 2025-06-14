package integration

import (
	"bytes"
	"context"
	"d-payroll/config"
	"d-payroll/controller/http"
	"d-payroll/entity"
	repository "d-payroll/repository/db"
	attendanceservice "d-payroll/service/attendance"
	authservice "d-payroll/service/auth"
	overtimeservice "d-payroll/service/overtime"
	payrollservice "d-payroll/service/payroll"
	reimbursementservice "d-payroll/service/reimbursement"
	userservice "d-payroll/service/user"
	"d-payroll/utils"
	"fmt"
	nethttp "net/http"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	testcontainers "github.com/testcontainers/testcontainers-go"
	postgrescontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
)

// TestApp holds the test application instance and configuration
type TestApp struct {
	App                  *fiber.App
	Config               *config.Config
	DB                   *repository.Db
	PgContainer          testcontainers.Container
	UserService          userservice.UserService
	AuthService          authservice.AuthService
	AttendanceService    attendanceservice.AttendanceService
	OvertimeService      overtimeservice.OvertimeService
	PayrollService       payrollservice.PayrollService
	ReimbursementService reimbursementservice.ReimbursementService
	AdminToken           string
	ctx                  context.Context
}

// SetupTestApp creates a new test application instance with PostgreSQL container
func SetupTestApp(t *testing.T) (*TestApp, error) {
	ctx := context.Background()

	// Start PostgreSQL container
	postgresContainer, pgHost, pgPort, err := startPostgresContainer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	// Use test-specific configuration with container details
	cfg := &config.Config{
		Postgres: &config.PostgresConfig{
			Host:     pgHost,
			Port:     int32(pgPort),
			User:     "postgres",
			Password: "postgres",
			Db:       "postgres",
		},
		Http: &config.HttpConfig{
			Host: "localhost",
			Port: 3000,
		},
		Auth: &config.AuthConfig{
			JwtSecret: "test-secret",
		},
		AdminUser: &config.AdminUserConfig{
			Username: "admin-test",
			Password: "test-password",
		},
		Overtime: &config.OvertimeConfig{
			MaxDurationPerDayMilis: 1000 * 60 * 60 * 3,
		},
	}

	// Connect to the database
	db, err := repository.NewDBHelper(*cfg)
	if err != nil {
		// Clean up container if DB connection fails
		postgresContainer.Terminate(ctx)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Apply database migrations
	err = applyMigrations(db.DB)
	if err != nil {
		// Clean up if migrations fail
		db.Close()
		postgresContainer.Terminate(ctx)
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	// Initialize repositories
	userDB := repository.NewUserDB(db.DB)
	attendanceDB := repository.NewAttendanceDB(db.DB)
	reimbursementDB := repository.NewReimbursementDB(db.DB)
	overtimeDB := repository.NewOvertimeDB(db.DB)
	payrollDB := repository.NewPayrollDB(db.DB)

	// Initialize services
	userSvc := userservice.NewUserService(userDB)
	authSvc := authservice.NewAuthService(cfg, userSvc)
	attendanceSvc := attendanceservice.NewAttendanceService(attendanceDB)
	reimbursementSvc := reimbursementservice.NewReimbursementService(reimbursementDB)
	overtimeSvc := overtimeservice.NewOvertimeService(cfg, overtimeDB, attendanceSvc)
	payrollSvc := payrollservice.NewPayrollService(cfg, payrollDB)

	// Initialize HTTP app
	httpApp := http.NewHttpApp(cfg)

	http.NewUserHttp(httpApp, userSvc)
	http.NewAuthHttp(httpApp, authSvc)
	http.NewAttendanceHttp(httpApp, attendanceSvc)
	http.NewReimbursementHttp(httpApp, reimbursementSvc)
	http.NewOvertimeHttp(httpApp, overtimeSvc)
	http.NewPayrollHttp(httpApp, payrollSvc)

	// Create test app
	testApp := &TestApp{
		App:                  httpApp.App,
		Config:               cfg,
		DB:                   db,
		PgContainer:          postgresContainer,
		UserService:          userSvc,
		AuthService:          authSvc,
		AttendanceService:    attendanceSvc,
		OvertimeService:      overtimeSvc,
		PayrollService:       payrollSvc,
		ReimbursementService: reimbursementSvc,
		ctx:                  ctx,
	}

	// Create admin user and get token
	adminToken, err := testApp.createAdminUser()
	if err != nil {
		testApp.TeardownTestApp()
		return nil, fmt.Errorf("failed to create admin user: %w", err)
	}
	testApp.AdminToken = adminToken

	return testApp, nil
}

// TeardownTestApp cleans up resources after tests
func (app *TestApp) TeardownTestApp() {
	// Close database connection
	if app.DB != nil {
		app.DB.Close()
	}

	// Stop PostgreSQL container
	if app.PgContainer != nil {
		app.PgContainer.Terminate(app.ctx)
	}
}

// startPostgresContainer starts a PostgreSQL container for testing
func startPostgresContainer(ctx context.Context) (testcontainers.Container, string, int, error) {
	postgresContainer, err := postgrescontainer.RunContainer(ctx,
		testcontainers.WithImage("postgres:14-alpine"),
		postgrescontainer.WithDatabase("postgres"),
		postgrescontainer.WithUsername("postgres"),
		postgrescontainer.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		return nil, "", 0, err
	}

	// Get host and port
	host, err := postgresContainer.Host(ctx)
	if err != nil {
		postgresContainer.Terminate(ctx)
		return nil, "", 0, err
	}

	mappedPort, err := postgresContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
		postgresContainer.Terminate(ctx)
		return nil, "", 0, err
	}

	return postgresContainer, host, mappedPort.Int(), nil
}

// createAdminUser creates an admin user for testing and returns a JWT token
func (app *TestApp) createAdminUser() (string, error) {
	// Create admin user
	salary := 5000000
	adminUser := &entity.User{
		Username: "admin",
		Password: "admin123",
		Role:     entity.UserRoleAdmin,
		UserInfo: &entity.UserInfo{
			MonthlySalary: &salary,
		},
	}

	// Create the user
	createdUser, err := app.UserService.CreateUser(app.ctx, adminUser)
	if err != nil {
		return "", err
	}

	// Generate token
	token, err := utils.GenerateToken(app.Config.Auth.JwtSecret, &entity.AuthTokenPayload{
		ID:   *createdUser.Id,
		Role: createdUser.Role,
	})
	if err != nil {
		return "", err
	}

	return token, nil
}

// Helper method to make authenticated requests
func (app *TestApp) makeAuthenticatedRequest(method, path string, body []byte, token string) (*nethttp.Request, error) {
	var req *nethttp.Request
	var err error

	if body != nil {
		req, err = nethttp.NewRequest(method, path, bytes.NewReader(body))
	} else {
		req, err = nethttp.NewRequest(method, path, nil)
	}

	if err != nil {
		return nil, err
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	return req, nil
}

func TestMain(m *testing.M) {
	// This is where we would do global setup if needed
	code := m.Run()
	os.Exit(code)
}

// Helper function to require a test container is running properly
// func requireContainerRunning(t *testing.T, container testcontainers.Container) {
// 	require.NotNil(t, container, "Container should not be nil")
// }

// applyMigrations runs the database migrations from the db/migrations folder
func applyMigrations(db *gorm.DB) error {
	// Get the migration files
	upSqlFile := "../../db/migrations/20250611154246_init_database.up.sql"

	// Read the migration file
	sqlBytes, err := os.ReadFile(upSqlFile)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Execute the SQL migration
	return db.Exec(string(sqlBytes)).Error
}
