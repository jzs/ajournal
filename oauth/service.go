package oauth

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/sketchground/ajournal/profile"
	"github.com/sketchground/ajournal/user"
)

// Service interface for oauth support
type Service interface {
	Login(ctx context.Context, username string, provider Provider) (*user.Token, error)
	Register(ctx context.Context, u *user.User, p *profile.Profile) error
}

// NewService returns a service implementation of the oauth service
func NewService(or Repository, ur user.Repository, pr profile.Repository) Service {
	return &service{
		repo: or,
		ur:   ur,
		pr:   pr,
	}
}

type service struct {
	repo Repository
	ur   user.Repository
	pr   profile.Repository
}

func (s *service) Login(ctx context.Context, username string, provider Provider) (*user.Token, error) {
	u, err := s.ur.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	providers, err := s.repo.Find(ctx, u.ID)
	if err != nil {
		return nil, err
	}

	for _, p := range providers {
		if p.Provider == provider {
			token := user.GenerateToken(u.ID)
			err = s.ur.CreateToken(ctx, token)
			if err != nil {
				return nil, errors.Wrap(err, "Oauth Login")
			}
			return token, nil
		}
	}

	return nil, errors.New("Could not find a valid provider")
}

func (s *service) Register(ctx context.Context, u *user.User, p *profile.Profile) error {
	u.Created = time.Now()
	cu, err := s.ur.Create(ctx, u)
	if err != nil {
		return errors.Wrap(err, "Oauth.Register")
	}

	p.ID = cu.ID
	_, err = s.pr.Create(ctx, p)
	if err != nil {
		return errors.Wrap(err, "Oauth.Register")
	}

	err = s.repo.Create(ctx, UserProvider{
		UserID:           cu.ID,
		Provider:         ProviderGoogle,
		ProviderUsername: u.Username,
	})
	if err != nil {
		return errors.Wrap(err, "Oauth.Register")
	}
	return nil
}
