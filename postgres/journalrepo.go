package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"

	"bitbucket.org/sketchground/journal/journal"
)

type journalRepo struct {
	db *sqlx.DB
}

// NewJournalRepo returns a new implementation of the Journal postgres Repository interface
func NewJournalRepo(db *sqlx.DB) journal.Repository {
	return &journalRepo{db: db}
}

func (jr *journalRepo) Create(ctx context.Context, journal *journal.Journal) (*journal.Journal, error) {
	var id int64
	err := jr.db.Get(&id, "INSERT INTO journal(UserID, Public, Title, Description, Created) VALUES($1, $2, $3, $4, $5)", journal.UserID, journal.Public, journal.Title, journal.Description, journal.Created)
	if err != nil {
		return nil, err
	}
	journal.ID = id
	return journal, nil
}

func (jr *journalRepo) FindAll(ctx context.Context) ([]*journal.Journal, error) {
	journals := []*journal.Journal{}
	err := jr.db.Select(journals, "SELECT * FROM journal")
	if err != nil {
		return nil, err
	}
	return journals, nil
}

func (jr *journalRepo) AddEntry(ctx context.Context, entry *journal.Entry, journalID int64) (*journal.Entry, error) {
	panic("Not implemented")
}

func (jr *journalRepo) UpdateEntry(ctx context.Context, entry *journal.Entry) error {
	panic("Not implemented")
}

func (jr *journalRepo) Entries(ctx context.Context, journalID int64) ([]*journal.Entry, error) {
	panic("Not implemented")
}
