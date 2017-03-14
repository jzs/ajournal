package profile

import (
	"context"
	"errors"
)

type Service interface {
	Profile(ctx context.Context) (*Profile, error)
}

func NewService(pr Repository) Service {
	return &service{}
}

type service struct {
}

func (s *service) Profile(ctx context.Context) (*Profile, error) {
	return nil, errors.New("Not implemented")
}
