package postgres_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/jzs/ajournal/postgres"
	"github.com/jzs/ajournal/user"
	"github.com/jzs/ajournal/utils/logger"
	"github.com/jzs/ajournal/utils/testhelpers"
)

func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	l := logger.NewTestLogger()
	td := testhelpers.SetupTestDB(ctx, l)
	defer td()

	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%v dbname=%v sslmode=disable", "jzs", "ajournal_test"))
	if err != nil {
		t.Fatalf("Could not connect to database! %v", err)
	}
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ur := postgres.NewUserRepo(db)

	now := time.Now()

	u, err := ur.Create(ctx, &user.User{Username: "bobo", Password: "bobo", Active: true, Created: now})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if u.Username != "bobo" {
		t.Errorf("Expected username bobo, got %v", u.Username)
	}

	if u.Password != "bobo" {
		t.Errorf("Expected password bobo, got %v", u.Username)
	}

	if !u.Created.Equal(now) {
		t.Errorf("Expected created %v, got %v", now, u.Created)
	}
}
