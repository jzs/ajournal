package user

import (
	"context"
	"errors"
	"time"
)

// ErrUserExist custom error
var ErrUserExist error

// ErrUserNotExist custom error
var ErrUserNotExist error

// ErrTokenNotExist custom error
var ErrTokenNotExist error

func init() {
	ErrUserExist = errors.New("User exists")
	ErrUserNotExist = errors.New("User does not exist")
	ErrTokenNotExist = errors.New("Token does not exist")
}

// Repository interface for User repository methods
type Repository interface {
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByToken(ctx context.Context, token string) (*User, error)
	Create(ctx context.Context, u *User) (*User, error)
	CreateToken(ctx context.Context, t *Token) error
	DeleteToken(ctx context.Context, token string) error
}

// User domain model
type User struct {
	ID       int64
	Username string
	Password string
	Active   bool
	Created  time.Time
}

// Token domain model
type Token struct {
	Token   string
	UserID  int64
	Expires time.Time
}
