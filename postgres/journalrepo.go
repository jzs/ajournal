package postgres

import (
	"context"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"bitbucket.org/sketchground/journal/journal"
)

type DBEntry struct {
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
	db *sqlx.DB
}

// NewJournalRepo returns a new implementation of the Journal postgres Repository interface
func NewJournalRepo(db *sqlx.DB) journal.Repository {
	return &journalRepo{db: db}
}

func (jr *journalRepo) Create(ctx context.Context, journal *journal.Journal) (*journal.Journal, error) {
	var id int64
	err := jr.db.Get(&id, "INSERT INTO journal(UserID, Public, Title, Description, Created) VALUES($1, $2, $3, $4, $5) RETURNING id", journal.UserID, journal.Public, journal.Title, journal.Description, journal.Created)
	if err != nil {
		return nil, err
	}
	journal.ID = id
	return journal, nil
}

func (jr *journalRepo) FindByID(ctx context.Context, id int64) (*journal.Journal, error) {
	j := &journal.Journal{}
	err := jr.db.Get(j, "SELECT * FROM Journal WHERE id=$1", id)
	if err != nil {
		log.Println(err.Error())
		return nil, journal.ErrJournalNotExist
	}
	return j, nil
}

func (jr *journalRepo) FindAll(ctx context.Context, userid int64) ([]*journal.Journal, error) {
	journals := []*journal.Journal{}
	err := jr.db.Select(&journals, "SELECT * FROM journal WHERE userid=$1", userid)
	if err != nil {
		return nil, err
	}
	return journals, nil
}

func (jr *journalRepo) AddEntry(ctx context.Context, entry *journal.Entry) (*journal.Entry, error) {
	var id int64
	err := jr.db.Get(&id, "INSERT INTO Entry(JournalID, Date, Title, Content, Created, IsPublished) VALUES($1, $2, $3, $4, $5, $6) RETURNING id", entry.JournalID, entry.Date, entry.Title, entry.Content, entry.Created, entry.IsPublished)
	if err != nil {
		log.Println(err)
		return nil, err
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
		return err
	}
	return nil
}

func (jr *journalRepo) FindEntryByID(ctx context.Context, id int64) (*journal.Entry, error) {
	e := &DBEntry{}
	err := jr.db.Get(e, "SELECT * FROM Entry WHERE ID=$1", id)
	if err != nil {
		log.Println(err.Error())
		return nil, journal.ErrEntryNotExist
	}

	result := mapToEntry(e)

	return result, nil
}

func (jr *journalRepo) FindAllEntries(ctx context.Context, journalID int64) ([]*journal.Entry, error) {
	entries := []*DBEntry{}
	err := jr.db.Select(&entries, "SELECT * FROM Entry WHERE journalid=$1 ORDER BY Created DESC", journalID)
	if err != nil {
		return nil, err
	}
	result := []*journal.Entry{}
	for _, e := range entries {
		result = append(result, mapToEntry(e))
	}
	return result, nil
}

func mapToEntry(e *DBEntry) *journal.Entry {
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
