package journal_test

import (
	"context"
	"testing"

	"bitbucket.org/sketchground/journal/journal"
	"bitbucket.org/sketchground/journal/user"
	"github.com/gorilla/mux"
)

func TestTransport(t *testing.T) {
	m := mux.NewRouter()
	jr := NewInmemRepo()
	js := journal.NewService(jr)
	journal.SetupHandler(m, js)
	// TODO: Test actual calls to routes...
}

func TestIntegration(t *testing.T) {
	// TODO: Implement integration test with real database...
}

func TestService(t *testing.T) {
	jr := NewInmemRepo()
	js := journal.NewService(jr)

	u := &user.User{
		ID:       201,
		Username: "jzs",
	}
	ctx := user.TestContextWithUser(u)
	journals, err := js.MyJournals(ctx)
	if err != nil {
		t.Fatalf("Expected to fetch journals, got: %v", err.Error())
	}
	if len(journals) > 0 {
		t.Fatalf("Expected empty journals, got: %v", len(journals))
	}

	title := "My first journal"
	jrnl, err := js.Create(ctx, &journal.Journal{Title: title})
	if err != nil {
		t.Fatalf("Expected insert journal, got: %v", err.Error())
	}
	if jrnl.Title != title {
		t.Fatalf("Expected Title %v , got: %v", title, jrnl.Title)
	}
	if jrnl.UserID != u.ID {
		t.Fatalf("Expected UserID %v , got: %v", u.ID, jrnl.UserID)
	}

	if jrnl.ID != 1 {
		t.Fatalf("Expected that journal got id 1, got: %v", err.Error())
	}
	_, err = js.Create(ctx, jrnl)
	if err == nil {
		t.Fatalf("Expected that journal create fail but got a newly created journal. ID must not be set")
	}
}

// In memory repository of journal

type journalRepo struct {
	journals []*journal.Journal
	entries  []*journal.Entry
	id       int64
	eid      int64
}

func NewInmemRepo() journal.Repository {
	repo := &journalRepo{
		journals: []*journal.Journal{},
		entries:  []*journal.Entry{},
		id:       1,
		eid:      1,
	}
	return repo
}

func (jr *journalRepo) Create(c context.Context, journal *journal.Journal) (*journal.Journal, error) {
	journal.ID = jr.id
	jr.journals = append(jr.journals, journal)
	jr.id = jr.id + 1
	return journal, nil
}

func (jr *journalRepo) FindAll(ctx context.Context, userid int64) ([]*journal.Journal, error) {
	js := []*journal.Journal{}
	for _, j := range jr.journals {
		if j.UserID == userid {
			js = append(js, j)
		}
	}
	return js, nil
}

func (jr *journalRepo) AddEntry(c context.Context, entry *journal.Entry, journalID int64) (*journal.Entry, error) {
	entry.ID = jr.eid
	jr.entries = append(jr.entries, entry)
	jr.eid = jr.eid + 1
	return entry, nil
}

func (jr *journalRepo) UpdateEntry(c context.Context, entry *journal.Entry) error {
	for i, e := range jr.entries {
		if e.ID == entry.ID {
			jr.entries[i] = entry
			break
		}
	}
	return nil
}

func (jr *journalRepo) Entries(ctx context.Context, journalID int64) ([]*journal.Entry, error) {
	entries := []*journal.Entry{}
	for _, e := range jr.entries {
		if e.JournalID == journalID {
			entries = append(entries, e)
		}
	}
	return entries, nil
}
