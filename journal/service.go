package journal

import (
	"context"
	"errors"

	"bitbucket.org/sketchground/journal/user"
)

// Service interface for journals
type Service interface {
	/* Create creates new journal with no entries in.*/
	Create(c context.Context, journal *Journal) (*Journal, error)
	MyJournals(c context.Context) ([]*Journal, error)
}

type service struct {
	repo Repository
}

// NewService returns a new implementation of the Service interface
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, journal *Journal) (*Journal, error) {
	usr := user.FromContext(ctx)
	if usr == nil {
		return nil, errors.New("Cannot create a journal without a user context")
	}

	if journal.ID != 0 {
		return nil, errors.New("Cannot create a journal with an already existing id set")
	}
	journal.UserID = usr.ID
	jrnl, err := s.repo.Create(ctx, journal)
	return jrnl, err
}

func (s *service) MyJournals(ctx context.Context) ([]*Journal, error) {
	usr := user.FromContext(ctx)
	if usr == nil {
		return nil, errors.New("Cannot fetch journals without a user context")
	}

	// TODO: Use userid when looking up journals!
	journals, err := s.repo.FindAll(ctx, usr.ID)
	if err != nil {
		return nil, err
	}
	return journals, nil
}
