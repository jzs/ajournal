package testhelpers

import (
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/mattes/migrate"
	pg "github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file" // for db migrations
	"github.com/sketchground/ajournal/app"
	"github.com/sketchground/ajournal/utils"
	"github.com/sketchground/ajournal/utils/logger"
)

// SetupTestingServer sets up a server for testing purposes
func SetupTestingServer(t testing.TB) (*httptest.Server, logger.Logger) {
	l := logger.NewTestLogger()
	// TODO: Consider flushing database instead. since dev db scripts suck...
	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%v dbname=%v sslmode=disable", "jzs", "aj_test"))
	if err != nil {
		t.Fatalf("Could not connect to database! %v", err)
	}

	driver, err := pg.WithInstance(db.DB, &pg.Config{})
	mig, err := migrate.NewWithDatabaseInstance(
		os.Getenv("BS_MIGRATIONDIR"),
		"postgres", driver)
	if err != nil {
		t.Fatal(err)
	}
	err = mig.Down()
	if err != nil && err != migrate.ErrNoChange {
		t.Fatal(err)
	}
	err = mig.Up()
	if err != nil && err != migrate.ErrNoChange {
		t.Fatal(err)
	}

	translator, err := utils.NewTranslator("../i18n", l)

	params := app.Params{
		WWWDir:    "",
		StripeKey: "",
	}

	m := app.SetupRouter(db, l, translator, params)
	s := httptest.NewServer(m)
	return s, l
}
