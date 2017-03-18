package profile

import (
	"context"
	"errors"

	"bitbucket.org/sketchground/ajournal/user"
)

type Service interface {
	Profile(ctx context.Context) (*Profile, error)
	UpdateProfile(ctx context.Context, p *Profile) (*Profile, error)
	Subscribe(ctx context.Context, sub *Subscription) error
}

func NewService(pr Repository) Service {
	return &service{pr: pr}
}

type service struct {
	pr Repository
}

func (s *service) Profile(ctx context.Context) (*Profile, error) {
	usr := user.FromContext(ctx)
	if usr == nil {
		return nil, errors.New("Cannot create a journal without a user context")
	}

	pro, err := s.pr.FindByID(ctx, usr.ID)
	if err == ErrProfileNotExist {
		// Create profile and return that.
		pro, err = s.pr.Create(ctx, &Profile{ID: usr.ID, Plan: PlanFree})
		if err != nil {
			return nil, errors.New("Could not create profile for user")
		}
		return pro, nil
	}
	return pro, nil
}

func (s *service) UpdateProfile(ctx context.Context, p *Profile) (*Profile, error) {
	usr := user.FromContext(ctx)
	if usr == nil {
		return nil, errors.New("Cannot create a journal without a user context")
	}
	if usr.ID != p.ID {
		return nil, errors.New("Cannot update another users profile")
	}

	return s.pr.Update(ctx, p)
}

func (s *service) Subscribe(ctx context.Context, sub *Subscription) error {
	return errors.New("Not implemented")
}
