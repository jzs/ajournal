package testhelpers

import (
	"context"
	"fmt"
	"net/http/httptest"
	"os"
	"os/exec"

	"github.com/mattes/migrate"
	ajournal "github.com/sketchground/ajournal/app"
	"github.com/sketchground/ajournal/utils/logger"

	_ "github.com/mattes/migrate/database/postgres" // support postgres migrate
	_ "github.com/mattes/migrate/source/file"
)

// SetupTestDB Creates a new empty test database, migrates it tear down function
func SetupTestDB(ctx context.Context, log logger.Logger) (teardown func() error) {
	// Recreate test db.
	if err := exec.Command("dropdb", "ajournal_test").Run(); err != nil {
		log.Fatalf(ctx, "Failed dropping test database: %v", err)
	}
	if err := exec.Command("createdb", "ajournal_test").Run(); err != nil {
		log.Fatalf(ctx, "Failed creating test database: %v", err)
	}

	ppath := fmt.Sprintf("%v/src/github.com/sketchground/ajournal", os.Getenv("GOPATH"))
	mig, err := migrate.New(fmt.Sprintf("file://%v/%v", ppath, "db/migrations"), "postgres://jzs@localhost:5432/ajournal_test?sslmode=disable")
	if err != nil {
		log.Fatalf(ctx, "Failed migrating database reason: %v", err)
	}

	if err := mig.Up(); err != nil {
		log.Fatalf(ctx, "Failed migrating database reason: %v", err)
	}

	return func() error {
		// Drop/clean database.
		err := mig.Down()
		if err != nil {
			return err
		}
		return nil
	}
}

// InitTestServer initiates a test server with a test transport
func InitTestServer() (*httptest.Server, logger.Logger, func() error) {
	ctx := context.Background()

	// TODO: Replace with logger that can be dumped to stdout only when a test fails.
	log := logger.NewTestLogger()

	teardown := SetupTestDB(ctx, log)

	ppath := fmt.Sprintf("%v/src/github.com/sketchground/ajournal", os.Getenv("GOPATH"))
	// Setenv for all needed configuration...
	cfg := ajournal.Configuration{
		TranslateFolder: fmt.Sprintf("%v/%v", ppath, "i18n"),
		DBName:          "ajournal_test",
		DBUser:          "jzs",
		S3Mock:          true,
	}

	h := ajournal.Setup(ctx, cfg, log)
	s := httptest.NewServer(h)

	return s, log, teardown
}
