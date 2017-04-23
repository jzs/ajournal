package journal

import (
	"context"
	"net/http"
	"time"

	"bitbucket.org/sketchground/ajournal/user"
	"bitbucket.org/sketchground/ajournal/utils"
	"github.com/pkg/errors"
	"github.com/russross/blackfriday"
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
		return nil, utils.NewAPIError(nil, http.StatusForbidden, "Cannot create a journal without a user context")
	}

	if journal.ID != 0 {
		return nil, utils.NewAPIError(nil, http.StatusBadRequest, "Cannot create a journal with an already existing id set")
	}
	journal.UserID = usr.ID
	jrnl, err := s.repo.Create(ctx, journal)
	if err != nil {
		return nil, errors.Wrap(err, "JournalService:Create")
	}
	jrnl.Entries = []*Entry{}
	return jrnl, nil
}

func (s *service) MyJournals(ctx context.Context) ([]*Journal, error) {
	usr := user.FromContext(ctx)
	if usr == nil {
		return nil, utils.NewAPIError(nil, http.StatusForbidden, "Cannot fetch journals without a user context")
	}

	journals, err := s.repo.FindAll(ctx, usr.ID)
	if err != nil {
		return nil, errors.Wrap(err, "MyJournals")
	}
	return journals, nil
}

func (s *service) Journal(ctx context.Context, id int64) (*Journal, error) {
	usr := user.FromContext(ctx)
	if usr == nil {
		return nil, utils.NewAPIError(nil, http.StatusForbidden, "Cannot fetch journals without a user context")
	}

	journal, err := s.repo.FindByID(ctx, id)
	switch {
	case err == ErrJournalNotExist:
		return nil, utils.NewAPIError(err, http.StatusNotFound, "Journal does not exist")
	case err != nil:
		return nil, err
	}

	// If it's another users journal and it is not public, then do not return it.
	if journal.UserID != usr.ID && journal.Public == false {
		return nil, utils.NewAPIError(errors.New("User trying to access another users private journal"), http.StatusNotFound, ErrJournalNotExist.Error())
	}

	entries, err := s.repo.FindAllEntries(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "Journal:FindAllEntries")
	}
	journal.Entries = entries

	return journal, nil
}

func (s *service) Journals(ctx context.Context, userid int64) ([]*Journal, error) {
	journals, err := s.repo.FindAll(ctx, userid)
	if err != nil {
		return nil, errors.Wrap(err, "Journals")
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
		return nil, utils.NewAPIError(nil, http.StatusBadRequest, "ID must not be set when creating a new entry")
	}
	if entry.JournalID == 0 {
		return nil, utils.NewAPIError(nil, http.StatusBadRequest, "JournalID must be set when creating a new entry")
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
		return nil, utils.NewAPIError(errors.New("Trying to create enty on another users journal"), http.StatusNotFound, ErrJournalNotExist.Error())
	}

	entry.Created = time.Now()

	return s.repo.AddEntry(ctx, entry)
}

func (s *service) UpdateEntry(ctx context.Context, entry *Entry) (*Entry, error) {
	if entry.ID == 0 {
		return nil, utils.NewAPIError(nil, http.StatusBadRequest, ErrEntryNotExist.Error())
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

	return entry, nil
}

func (s *service) Entry(ctx context.Context, id int64) (*Entry, error) {
	usr := user.FromContext(ctx)

	entry, err := s.repo.FindEntryByID(ctx, id)
	if err != nil {
		return nil, utils.NewAPIError(err, http.StatusNotFound, ErrEntryNotExist.Error())
	}

	j, err := s.repo.FindByID(ctx, entry.JournalID)
	if err != nil {
		return nil, utils.NewAPIError(nil, http.StatusNotFound, ErrEntryNotExist.Error())
	}
	// If it is private and you are not logged in
	if !j.Public && usr == nil {
		return nil, utils.NewAPIError(errors.New("Journal belongs to another user"), http.StatusNotFound, ErrEntryNotExist.Error())
	}
	// If it is private and you are logged in as different user
	if !j.Public && (usr != nil && usr.ID != j.UserID) {
		return nil, utils.NewAPIError(errors.New("Journal belongs to another user"), http.StatusNotFound, ErrEntryNotExist.Error())
	}

	// Render content
	entry.HTMLContent = string(blackfriday.MarkdownCommon([]byte(entry.Content)))

	return entry, nil
}
