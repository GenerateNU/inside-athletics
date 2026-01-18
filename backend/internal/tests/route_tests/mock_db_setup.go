package routeTests

import (
	"context"
	"encoding/json"
	"fmt"
	"inside-athletics/internal/server"
	"log"
	"net/http/httptest"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestDatabase holds the test database connection and container
type TestDatabase struct {
	DB        *gorm.DB
	Container *postgres.PostgresContainer
	API       humatest.TestAPI
}

// SetupTestDB creates a new PostgreSQL container and returns a connection
func SetupTestDB(t *testing.T) *TestDatabase {
	ctx := context.Background()

	// Create PostgreSQL container
	postgresContainer, err := postgres.Run(ctx,
		"postgres:17.4",
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("test_user"),
		postgres.WithPassword("test_password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second)),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %s", err)
	}

	// Get connection string
	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}

	// Connect to database and setup API
	api, db := SetupTestAPI(t, connStr)

	// Verify connection
	sqlDb, err := db.DB()

	if err != nil {
		t.Fatal(err)
	}

	if err := sqlDb.Ping(); err != nil {
		t.Fatalf("failed to ping database: %s", err)
	}

	testDB := &TestDatabase{
		DB:        db,
		Container: postgresContainer,
		API:       api,
	}

	// Run migrations to sync schemas with temporary DB
	testDB.RunMigrations(t)

	return testDB
}

// Teardown cleans up the test database
func (td *TestDatabase) Teardown(t *testing.T) {
	ctx := context.Background()

	if td.DB != nil {
		sqlDb, err := td.DB.DB()

		if err != nil {
			t.Fatal(err)
		}

		sqlDb.Close()
	}

	if td.Container != nil {
		if err := td.Container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}
}

func (td *TestDatabase) RunMigrations(t *testing.T) {
	ctx := context.Background()

	// Get the connection string for Atlas
	connStr, err := td.Container.ConnectionString(ctx, "sslmode=disable", "search_path=public")
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}

	_, filename, _, _ := runtime.Caller(0)
	// Go up from current file to project root
	backendDir := filepath.Join(filepath.Dir(filename), "..", "..")
	migrationDir := filepath.Join(backendDir, "migrations")

	// Run Atlas migrations using exec
	cmd := exec.Command("atlas", "migrate", "apply",
		"--dir", fmt.Sprintf("file://%s", migrationDir),
		"--url", connStr,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to run atlas migrations: %s\nOutput: %s", err, output)
	}
}

// GENERIC HELPER FUNCS

/*
Decode the given response JSON into the given struct entity. Reads the value into the struct
*/
func DecodeTo[T any](entity *T, resp *httptest.ResponseRecorder) {
	dec := json.NewDecoder(resp.Body)

	err := dec.Decode(entity)
	if err != nil {
		log.Fatalf("decode error: %v", err)
	}
}

/*
Create API routing with test DB connection based on given dbUrl
*/
func SetupTestAPI(t *testing.T, dbUrl string) (humatest.TestAPI, *gorm.DB) {
	_, api := humatest.New(t) // setup test API

	db, err := gorm.Open(gormPostgres.Open(dbUrl), &gorm.Config{})

	if err != nil {
		t.Errorf("Unable to connect to DB: %v", err)
	}

	server.CreateRoutes(db, api)

	return api, db
}
