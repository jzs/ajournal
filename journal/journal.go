package journal

import (
	"context"
	"errors"
	"time"

	"github.com/sketchground/ajournal/common"
)

// ErrJournalNotExist error
var ErrJournalNotExist error

// ErrEntryNotExist error
var ErrEntryNotExist error

func init() {
	ErrJournalNotExist = errors.New("Journal does not exist")
	ErrEntryNotExist = errors.New("Entry does not exist")
}

// Journal describes a travel journal describing ones travel adventure.
type Journal struct {
	ID          int64
	UserID      int64      // The user that the journal belongs to
	Public      bool       // Describes whether the journal is public or private
	Title       string     // The title of the journal.
	Description string     // A description of the journal. It could be the initiary of the trip or the goal
	Tags        []string   // Tags of journal.
	Entries     uint64     // a count of how many entries in the journal
	From        *time.Time // Starting time of the journal. Especially interesting for a travel journal
	To          *time.Time // The ending time of the journal.
	Locations   []string   // Locations that was visited during the journal.
	Created     time.Time  // The creation time of the journal
}

// Entry domain model
type Entry struct {
	ID          int64
	JournalID   int64
	Date        time.Time
	Title       string    // The title of the entry in the journal
	Content     string    // Content in markdown format
	HTMLContent string    `json:"HtmlContent"` // A rendered html version of the markdown content
	Tags        []string  // Tags of entry
	Created     time.Time // The creation time of the entry
	Published   time.Time // The publish time of the entry
	IsPublished bool      // Marks whether the entry is published or not
}

// Repository interface for journal repositories
type Repository interface {
	Create(ctx context.Context, journal *Journal) (*Journal, error)
	FindByID(ctx context.Context, id int64) (*Journal, error)
	FindAll(ctx context.Context, userid int64) ([]*Journal, error)
	AddEntry(ctx context.Context, entry *Entry) (*Entry, error)
	UpdateEntry(ctx context.Context, entry *Entry) error
	FindEntryByID(ctx context.Context, id int64) (*Entry, error)
	FindAllEntries(ctx context.Context, journalID int64, args common.PaginationArgs) ([]*Entry, error)
}
