package postgres

import (
	"context"
	"database/sql"
	"time"

	"bitbucket.org/sketchground/ajournal/journal"
	"bitbucket.org/sketchground/ajournal/utils/logger"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type dbEntry struct {
	ID          int64
	JournalID   int64
	Date        time.Time
	Title       string      // The title of the entry in the journal
	Content     string      // Content in markdown format
	Tags        []string    // Tags of entry
	Created     time.Time   // The creation time of the entry
	Published   pq.NullTime // The publish time of the entry
	IsPublished bool        // Marks whether the entry is published or not
}

type journalRepo struct {
	db     *sqlx.DB
	logger logger.Logger
}

// NewJournalRepo returns a new implementation of the Journal postgres Repository interface
func NewJournalRepo(db *sqlx.DB, logger logger.Logger) journal.Repository {
	return &journalRepo{db: db, logger: logger}
}

func (jr *journalRepo) Create(ctx context.Context, journal *journal.Journal) (*journal.Journal, error) {
	var id int64
	err := jr.db.Get(&id, "INSERT INTO journal(UserID, Public, Title, Description, Created) VALUES($1, $2, $3, $4, $5) RETURNING id", journal.UserID, journal.Public, journal.Title, journal.Description, journal.Created)
	if err != nil {
		return nil, errors.Wrap(err, "JournalRepo failed Create")
	}
	journal.ID = id
	return journal, nil
}

func (jr *journalRepo) FindByID(ctx context.Context, id int64) (*journal.Journal, error) {
	j := &journal.Journal{}
	err := jr.db.Get(j, "SELECT * FROM Journal WHERE id=$1", id)
	switch {
	case err == sql.ErrNoRows:
		return nil, journal.ErrJournalNotExist
	case err != nil:
		return nil, errors.Wrap(err, "JournalRepo:FindByID")
	}

	return j, nil
}

func (jr *journalRepo) FindAll(ctx context.Context, userid int64) ([]*journal.Journal, error) {
	journals := []*journal.Journal{}
	err := jr.db.Select(&journals, "SELECT * FROM journal WHERE userid=$1", userid)
	if err != nil {
		return nil, errors.Wrap(err, "JournalRepo:FindAll failed")
	}
	return journals, nil
}

func (jr *journalRepo) AddEntry(ctx context.Context, entry *journal.Entry) (*journal.Entry, error) {
	var id int64
	err := jr.db.Get(&id, "INSERT INTO Entry(JournalID, Date, Title, Content, Created, IsPublished) VALUES($1, $2, $3, $4, $5, $6) RETURNING id", entry.JournalID, entry.Date, entry.Title, entry.Content, entry.Created, entry.IsPublished)
	if err != nil {
		return nil, errors.Wrap(err, "JournalRepo:AddEntry failed")
	}
	entry.ID = id
	return entry, nil
}

func (jr *journalRepo) UpdateEntry(ctx context.Context, entry *journal.Entry) error {

	var published *time.Time
	if !entry.Published.IsZero() {
		published = &entry.Published
	}
	_, err := jr.db.Exec("UPDATE Entry SET Date=$1, Title=$2, Content=$3, Published=$4, IsPublished=$5 WHERE id=$6", entry.Date, entry.Title, entry.Content, published, entry.IsPublished, entry.ID)
	if err != nil {
		return errors.Wrap(err, "JournalRepo:UpdateEntry failed")
	}
	return nil
}

func (jr *journalRepo) FindEntryByID(ctx context.Context, id int64) (*journal.Entry, error) {
	e := &dbEntry{}
	err := jr.db.Get(e, "SELECT * FROM Entry WHERE ID=$1", id)
	switch {
	case err == sql.ErrNoRows:
		return nil, journal.ErrEntryNotExist
	case err != nil:
		return nil, errors.Wrap(err, "JournalRepo:FindEntryByID")
	}

	result := mapToEntry(e)

	return result, nil
}

func (jr *journalRepo) FindAllEntries(ctx context.Context, journalID int64) ([]*journal.Entry, error) {
	entries := []*dbEntry{}
	err := jr.db.Select(&entries, "SELECT * FROM Entry WHERE journalid=$1 ORDER BY Created DESC", journalID)
	if err != nil {
		return nil, errors.Wrap(err, "JournalRepo:FindAllEntries failed")
	}
	result := []*journal.Entry{}
	for _, e := range entries {
		result = append(result, mapToEntry(e))
	}
	return result, nil
}

func mapToEntry(e *dbEntry) *journal.Entry {
	var date time.Time
	if e.Published.Valid {
		date = e.Published.Time
	}
	result := &journal.Entry{
		ID:          e.ID,
		JournalID:   e.JournalID,
		Date:        e.Date,
		Title:       e.Title,   // The title of the entry in the journal
		Content:     e.Content, // Content in markdown format
		Created:     e.Created, // The creation time of the entry
		Published:   date,
		IsPublished: e.IsPublished, // Marks whether the entry is published or not
	}
	return result
}
