package journal

import (
	"context"
	"errors"
	"log"
	"time"

	"bitbucket.org/sketchground/journal/user"
)

// Service interface for journals
type Service interface {
	/* Create creates new journal with no entries in.*/
	Create(ctx context.Context, journal *Journal) (*Journal, error)
	MyJournals(ctx context.Context) ([]*Journal, error)
	Journal(ctx context.Context, id int64) (*Journal, error)
	Journals(ctx context.Context, userid int64) ([]*Journal, error)
	// Interfaces for entry creation
	CreateEntry(ctx context.Context, entry *Entry) (*Entry, error)
	UpdateEntry(ctx context.Context, entry *Entry) (*Entry, error)
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

func (s *service) Journals(ctx context.Context, userid int64) ([]*Journal, error) {
	journals, err := s.repo.FindAll(ctx, userid)
	if err != nil {
		return nil, err
	}

	result := []*Journal{}
	for _, j := range journals {
		if j.Public {
			result = append(result, j)
		}
	}
	return result, nil
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

func (s *service) UpdateEntry(ctx context.Context, entry *Entry) (*Entry, error) {
	if entry.ID == 0 {
		return nil, ErrEntryNotExist
	}

	usr := user.FromContext(ctx)
	if usr == nil {
		return nil, errors.New("Cannot update entry without a user context")
	}

	ntry, err := s.repo.FindEntryByID(ctx, entry.ID)
	if err != nil {
		return nil, ErrEntryNotExist
	}

	journal, err := s.repo.FindByID(ctx, ntry.JournalID)
	if err != nil {
		return nil, err
	}
	if journal.UserID != usr.ID {
		return nil, ErrEntryNotExist
	}

	err = s.repo.UpdateEntry(ctx, entry)
	if err != nil {
		return nil, ErrEntryNotExist
	}

	return ntry, nil

}

func (s *service) Entry(ctx context.Context, id int64) (*Entry, error) {
	usr := user.FromContext(ctx)

	entry, err := s.repo.FindEntryByID(ctx, id)
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("Cannot fetch entry")
	}

	j, err := s.repo.FindByID(ctx, entry.JournalID)
	if err != nil {
		log.Println(err.Error())
		return nil, ErrEntryNotExist
	}
	// If it is private and you are not logged in
	if !j.Public && usr == nil {
		return nil, ErrEntryNotExist
	}
	// If it is private and you are logged in as different user
	if !j.Public && (usr != nil && usr.ID != j.UserID) {
		return nil, ErrEntryNotExist
	}

	return entry, nil
}
