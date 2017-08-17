package journal

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/russross/blackfriday"
	"github.com/sketchground/ajournal/common"
	"github.com/sketchground/ajournal/user"
	"github.com/sketchground/ajournal/utils"
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
	Entries(ctx context.Context, id int64, args common.PaginationArgs) (*Entries, error)
	LatestJournals(ctx context.Context, count uint64) (*LatestJournals, error)
}

// Entries represent a paginated version of entries
type Entries struct {
	common.Pagination
	Entries []*Entry
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
	jrnl.Entries = 0
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
	journal, err := s.repo.FindByID(ctx, id)
	switch {
	case err == ErrJournalNotExist:
		return nil, utils.NewAPIError(err, http.StatusNotFound, "Journal does not exist")
	case err != nil:
		return nil, err
	}

	usr := user.FromContext(ctx)
	// If it's another users journal and it is not public, then do not return it.
	if !journal.Public && usr != nil && journal.UserID != usr.ID {
		return nil, utils.NewAPIError(errors.New("User trying to access another users private journal"), http.StatusNotFound, ErrJournalNotExist.Error())
	}

	entries, err := s.repo.FindAllEntries(ctx, id, common.PaginationArgs{Limit: 10000})
	if err != nil {
		return nil, errors.Wrap(err, "Journal:FindAllEntries")
	}
	journal.Entries = uint64(len(entries))

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

func (s *service) Entries(ctx context.Context, journalID int64, args common.PaginationArgs) (*Entries, error) {
	usr := user.FromContext(ctx)
	j, err := s.repo.FindByID(ctx, journalID)
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

	// Return all entries based on
	args.Limit++ // Fetch one more than we need, to support has next...
	entries, err := s.repo.FindAllEntries(ctx, journalID, args)
	if err != nil {
		return nil, errors.Wrap(err, "Journal:FindAllEntries")
	}

	res := &Entries{Entries: entries}
	if uint64(len(entries)) > args.Limit-1 {
		res.HasNext = true
		res.Entries = res.Entries[:args.Limit-1]
		from, err := strconv.Atoi(args.From)
		if err != nil {
			from = 0
		}
		res.Next = fmt.Sprint(uint64(from) + args.Limit - 1)
	}
	return res, nil
}

// LatestJournals is a view model for returning the x latest journals
type LatestJournals struct {
	common.Pagination
	Journals []*LatestJournal
}

// LatestJournal is a view model for returning the x latest journals
type LatestJournal struct {
	ID    int64
	Title string
	Entry *Entry
}

// LatestJournals returns a view of the latest x journals with the newest entry on.
func (s *service) LatestJournals(ctx context.Context, count uint64) (*LatestJournals, error) {
	pag := common.DefaultPagination()
	pag.Limit = count
	journals, err := s.repo.FindNewest(ctx, pag)
	if err != nil {
		return nil, err
	}
	result := []*LatestJournal{}
	for _, j := range journals {
		if !j.Public {
			return nil, errors.New("Expected no private journals from repository")
		}

		entries, err := s.repo.FindAllEntries(ctx, j.ID, common.PaginationArgs{Limit: 1})
		if err != nil {
			return nil, err
		}
		if len(entries) != 1 {
			return nil, errors.New("Expected one entry in journal, found none")
		}

		result = append(result, &LatestJournal{
			ID:    j.ID,
			Title: j.Title,
			Entry: entries[0],
		})
	}

	return &LatestJournals{Journals: result}, nil
}
