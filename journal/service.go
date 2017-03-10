package journal

import (
	"context"
	"errors"
	"time"

	"bitbucket.org/sketchground/journal/user"
)

// Service interface for journals
type Service interface {
	/* Create creates new journal with no entries in.*/
	Create(ctx context.Context, journal *Journal) (*Journal, error)
	MyJournals(ctx context.Context) ([]*Journal, error)
	Journal(ctx context.Context, id int64) (*Journal, error)
	// Interfaces for entry creation
	CreateEntry(ctx context.Context, entry *Entry) (*Entry, error)
	Entry(ctx context.Context, id int64) (*Entry, error)
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
	jrnl.Entries = []*Entry{}
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

func (s *service) Journal(ctx context.Context, id int64) (*Journal, error) {
	usr := user.FromContext(ctx)
	if usr == nil {
		return nil, errors.New("Cannot fetch journals without a user context")
	}

	journal, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	entries, err := s.repo.FindAllEntries(ctx, id)
	if err != nil {
		return nil, err
	}
	journal.Entries = entries

	// If it's another users journal and it is not public, then do not return it.
	if journal.UserID != usr.ID && journal.Public == false {
		return nil, ErrJournalNotExist
	}

	return journal, nil
}

func (s *service) CreateEntry(ctx context.Context, entry *Entry) (*Entry, error) {
	if entry.ID != 0 {
		return nil, errors.New("ID must not be set when creating a new entry")
	}
	if entry.JournalID == 0 {
		return nil, errors.New("JournalID must be set when creating a new entry")
	}

	usr := user.FromContext(ctx)
	if usr == nil {
		return nil, errors.New("Cannot fetch journals without a user context")
	}

	journal, err := s.repo.FindByID(ctx, entry.JournalID)
	if err != nil {
		return nil, err
	}
	if journal.UserID != usr.ID {
		return nil, ErrJournalNotExist
	}

	entry.Created = time.Now()

	return s.repo.AddEntry(ctx, entry)
}

func (s *service) Entry(ctx context.Context, id int64) (*Entry, error) {
	panic("Not implemented")
}
