package user_test

import (
	"context"
	"testing"

	"github.com/sketchground/ajournal/user"
)

func TestGetUserUnauthorized(t *testing.T) {
	ur := NewInmemRepo()
	us := user.NewService(ur)
	ctx := context.Background()
	u := &user.User{Username: "bob@cat.de", Password: "ewifj"}
	err := us.Register(ctx, u)
	if err != nil {
		t.Fatalf("Expected create user success, got: %v", err)
	}
	guser, err := us.User(ctx, u.Username)
	if err != nil {
		t.Fatalf("Expected user, found none")
	}
	if guser.Password != "" {
		t.Fatalf("Expected password not to be set, got: %v", guser.Password)
	}
}

func TestService(t *testing.T) {
	ur := NewInmemRepo()
	us := user.NewService(ur)

	ctx := context.Background()

	// create user
	err := us.Register(ctx, &user.User{})
	if err == nil {
		t.Fatalf("Expected error no username, got no error")
	}
	err = us.Register(ctx, &user.User{Username: "bob@cat.de"})
	if err == nil {
		t.Fatalf("Expected error no password, got no error")
	}

	username := "bob@cat.de"
	pass := "12345"

	err = us.Register(ctx, &user.User{Username: username, Password: pass})
	if err != nil {
		t.Fatalf("Expected create user, got %v", err.Error())
	}
	// TODO: Check if properties are set correctly on newly created user...

	err = us.Register(ctx, &user.User{Username: username, Password: pass})
	if err == nil {
		t.Fatalf("Expected error user already exist, got no error")
	}

	u, err := us.User(ctx, username)
	if err != nil {
		t.Fatalf("Expected user, found none")
	}
	if u.Password != "" {
		t.Fatalf("Expected no password, got %v", u.Password)
	}
	if !u.Active {
		t.Fatalf("Expected newly created user to be active. It Isn't")
	}
	if u.Created.IsZero() {
		t.Fatalf("Expected created time set. It looks like it isn't")
	}

	// Perform login with previously created user.
	token, err := us.Login(ctx, username, pass)
	if err != nil {
		t.Fatalf("Expected successful login, got %v", err.Error())
	}
	if token.Token == "" {
		t.Fatalf("Expected token, got empty token")
	}
	if token.UserID != u.ID {
		t.Fatalf("Expected token for same user, got other users token")
	}

	err = us.Logout(ctx, token.Token)
	if err != nil {
		t.Fatalf("Expected successful logout, got: %v", err.Error())
	}
	user, err := us.UserWithToken(ctx, token.Token)
	if err == nil {
		t.Fatalf("Expected user to be logged out. But found user %v", user)
	}
}

type userRepo struct {
	users  []*user.User
	tokens []*user.Token
	id     int64
}

func NewInmemRepo() user.Repository {
	repo := &userRepo{
		users:  []*user.User{},
		tokens: []*user.Token{},
		id:     1,
	}
	return repo
}

func (ur *userRepo) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	for _, u := range ur.users {
		if u.Username == username {
			u2 := *u        // Dereference to create a copy.
			return &u2, nil // Return pointer to new copy
		}
	}
	return nil, user.ErrUserNotExist
}

func (ur *userRepo) FindByToken(ctx context.Context, token string) (*user.User, error) {
	for _, t := range ur.tokens {
		if t.Token == token {
			// We found the right token!
			for _, u := range ur.users {
				if u.ID == t.UserID {
					return u, nil
				}
			}
		}
	}
	return nil, user.ErrTokenNotExist
}

func (ur *userRepo) Create(ctx context.Context, u *user.User) (*user.User, error) {
	u.ID = ur.id
	ur.users = append(ur.users, u)
	ur.id = ur.id + 1
	return u, nil
}

func (ur *userRepo) CreateToken(ctx context.Context, t *user.Token) error {
	ur.tokens = append(ur.tokens, t)
	return nil
}

func (ur *userRepo) DeleteToken(ctx context.Context, token string) error {
	for i, t := range ur.tokens {
		if t.Token == token {
			// Delete it
			ur.tokens = append(ur.tokens[:i], ur.tokens[i+1:]...)
			return nil
		}
	}
	return user.ErrTokenNotExist
}
