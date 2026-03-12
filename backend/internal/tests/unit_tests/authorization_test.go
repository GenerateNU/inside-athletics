package unitTests

import (
	"context"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"inside-athletics/internal/models"
	"inside-athletics/internal/server"

	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestIsOwnerOfPostOrComment(t *testing.T) {
	testDB := setupAuthTestDB(t)
	defer testDB.Teardown(t)

	userID := uuid.New()
	user := models.User{
		ID:                      userID,
		FirstName:               "Owner",
		LastName:                "User",
		Email:                   "owner@example.com",
		Username:                "owner",
		Account_Type:            false,
		Verified_Athlete_Status: models.VerifiedAthleteStatusPending,
	}
	if err := testDB.DB.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	otherUserID := uuid.New()
	otherUser := models.User{
		ID:                      otherUserID,
		FirstName:               "Other",
		LastName:                "User",
		Email:                   "other@example.com",
		Username:                "other",
		Account_Type:            false,
		Verified_Athlete_Status: models.VerifiedAthleteStatusPending,
	}
	if err := testDB.DB.Create(&otherUser).Error; err != nil {
		t.Fatalf("failed to create other user: %v", err)
	}

	post := models.Post{
		AuthorID: userID,
		Title:    "Test Post",
		Content:  "Test Content",
	}
	if err := testDB.DB.Create(&post).Error; err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	comment := models.Comment{
		UserID:      userID,
		PostID:      post.ID,
		Description: "Test Comment",
	}
	if err := testDB.DB.Create(&comment).Error; err != nil {
		t.Fatalf("failed to create comment: %v", err)
	}

	owned, err := server.IsOwnerOfPostOrComment(testDB.DB, userID, post.ID.String(), "post")
	if err != nil || !owned {
		t.Fatalf("expected owner for post, got owned=%v err=%v", owned, err)
	}

	owned, err = server.IsOwnerOfPostOrComment(testDB.DB, otherUserID, post.ID.String(), "post")
	if err != nil || owned {
		t.Fatalf("expected non-owner for post, got owned=%v err=%v", owned, err)
	}

	owned, err = server.IsOwnerOfPostOrComment(testDB.DB, userID, comment.ID.String(), "comment")
	if err != nil || !owned {
		t.Fatalf("expected owner for comment, got owned=%v err=%v", owned, err)
	}

	owned, err = server.IsOwnerOfPostOrComment(testDB.DB, userID, post.ID.String(), "sport")
	if err == nil || owned {
		t.Fatalf("expected error for unsupported resource, got owned=%v err=%v", owned, err)
	}

	owned, err = server.IsOwnerOfPostOrComment(testDB.DB, userID, "not-a-uuid", "post")
	if err == nil || owned {
		t.Fatalf("expected error for invalid resource ID, got owned=%v err=%v", owned, err)
	}
}

type authTestDB struct {
	DB        *gorm.DB
	Container *postgres.PostgresContainer
}

func setupAuthTestDB(t *testing.T) *authTestDB {
	t.Helper()

	ctx := context.Background()
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

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}

	db, err := gorm.Open(gormPostgres.Open(connStr), &gorm.Config{})
	if err != nil {
		t.Fatalf("unable to connect to DB: %v", err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	if err := sqlDb.Ping(); err != nil {
		t.Fatalf("failed to ping database: %s", err)
	}

	td := &authTestDB{
		DB:        db,
		Container: postgresContainer,
	}

	td.runMigrations(t)
	return td
}

func (td *authTestDB) Teardown(t *testing.T) {
	t.Helper()
	ctx := context.Background()

	if td.DB != nil {
		sqlDb, err := td.DB.DB()
		if err != nil {
			t.Fatal(err)
		}
		if err := sqlDb.Close(); err != nil {
			t.Fatalf("unable to close DB connection %s", err.Error())
		}
	}

	if td.Container != nil {
		if err := td.Container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}
}

func (td *authTestDB) runMigrations(t *testing.T) {
	t.Helper()
	ctx := context.Background()

	connStr, err := td.Container.ConnectionString(ctx, "sslmode=disable", "search_path=public")
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}

	_, filename, _, _ := runtime.Caller(0)
	backendDir := filepath.Join(filepath.Dir(filename), "..", "..")
	migrationDir := filepath.Join("internal", "migrations")

	cmd := exec.Command("atlas", "migrate", "apply",
		"--dir", "file://"+filepath.ToSlash(migrationDir),
		"--url", connStr,
	)
	cmd.Dir = backendDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to run atlas migrations: %s\nOutput: %s", err, output)
	}
}
