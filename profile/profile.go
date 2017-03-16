package profile

import (
	"context"
	"errors"
)

// ErrProfileNotExist custom error
var ErrProfileNotExist error

func init() {
	ErrProfileNotExist = errors.New("Profile does not exist")
}

type Repository interface {
	Create(ctx context.Context, p *Profile) (*Profile, error)
	FindByID(ctx context.Context, id int64) (*Profile, error)
}

type Profile struct {
	ID    int64
	Name  string
	Email string
}
