package journal_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"bitbucket.org/sketchground/ajournal/journal"
	"bitbucket.org/sketchground/ajournal/user"

	"github.com/gorilla/mux"
)

type logger struct{}

func (l *logger) Error(ctx context.Context, err error)                                    {}
func (l *logger) Errorf(ctx context.Context, format string, args ...interface{})          {}
func (l *logger) Print(ctx context.Context, err error)                                    {}
func (l *logger) Printf(ctx context.Context, format string, args ...interface{})          {}
func (l *logger) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) { next(w, r) }

func TestTransport(t *testing.T) {
	m := mux.NewRouter()
	jr := NewInmemRepo()
	js := journal.NewService(jr)
	journal.SetupHandler(m, js, &logger{})

	testTransport(t, m, "/users/1/journals", http.StatusOK)
	testTransport(t, m, "/journals", http.StatusForbidden)
	testTransport(t, m, "/journals/1", http.StatusForbidden)
	testTransport(t, m, "/journals/1/entries/1", http.StatusNotFound)
}

func testTransport(t *testing.T, m *mux.Router, url string, response int) {
	req, _ := http.NewRequest("GET", url, nil)
	rw := httptest.NewRecorder()
	m.ServeHTTP(rw, req)
	if rw.Code != response {
		t.Errorf("Expected %v on url %v, got %v", response, url, rw.Code)
	}
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
	if jrnl.Entries == nil {
		t.Fatalf("Expected empty Entries list. got nil object")
	}

	if jrnl.ID != 1 {
		t.Fatalf("Expected that journal got id 1, got: %v", err.Error())
	}
	_, err = js.Create(ctx, jrnl)
	if err == nil {
		t.Fatalf("Expected that journal create fail but got a newly created journal. ID must not be set")
	}

	// Test create entry in journal...
	entry := &journal.Entry{}
	ntry, err := js.CreateEntry(ctx, entry)
	if err == nil {
		t.Fatalf("Expected error creating entry, got: %v", ntry)
	}
	entry.ID = 1
	ntry, err = js.CreateEntry(ctx, entry)
	if err == nil {
		t.Fatalf("Expected error creating entry bad arg, got: %v", ntry)
	}
	entry.ID = 0
	entry.JournalID = 1
	entry.Content = "#Hello World"
	ntry, err = js.CreateEntry(ctx, entry)
	if err != nil {
		t.Fatalf("Expected sucessful creation of entry , got: %v", err.Error())
	}
	if ntry.Title != entry.Title {
		t.Fatalf("Expected entry title %v, got: %v", entry.Title, ntry.Title)
	}

	ntry, err = js.Entry(ctx, ntry.ID)
	if err != nil {
		t.Fatalf("Expected to fetch entry just created, got: %v", err.Error())
	}
	if ntry.Title != entry.Title {
		t.Fatalf("Expected title %v, got: %v", entry.Title, ntry.Title)
	}
	if ntry.Content == "" {
		t.Fatalf("Expected content: %v, got: %v", entry.Content, ntry.Content)
	}
	if entry.HTMLContent == "" {
		t.Fatalf("Expected Html rendered content, got empty string")
	}

	// Test if we can update the entry
	ntry, err = js.UpdateEntry(ctx, ntry)
	if err != nil {
		t.Fatalf("Could not update the entry, got: %v", err.Error())
	}
	if ntry.Title != entry.Title {
		t.Fatalf("Expected title: %v, got %v", entry.Title, ntry.Title)
	}

	// Test if we now return 1 journal
	journals, err = js.MyJournals(ctx)
	if err != nil {
		t.Fatalf("Expected to fetch journals, got: %v", err.Error())
	}
	if len(journals) != 1 {
		t.Fatalf("Expected 1 journal, got: %v", len(journals))
	}

	// Test if we can fetch that one journal
	j, err := js.Journal(ctx, jrnl.ID)
	if err != nil {
		t.Fatalf("Expected to fetch the journal, got: %v", err.Error())
	}
	if j.Title != title {
		t.Fatalf("Expected Title %v , got: %v", title, jrnl.Title)
	}
	if j.UserID != u.ID {
		t.Fatalf("Expected UserID %v , got: %v", u.ID, jrnl.UserID)
	}
	if j.ID != 1 {
		t.Fatalf("Expected that journal got id 1, got: %v", err.Error())
	}
	if len(j.Entries) != 1 {
		t.Fatalf("Expected one entry in journal, got: %v", len(j.Entries))
	}

	// Test if we don't return other users journals
	u = &user.User{
		ID:       500,
		Username: "Bob",
	}
	ctx = user.TestContextWithUser(u)
	journals, err = js.MyJournals(ctx)
	if err != nil {
		t.Fatalf("Expected to fetch 0 journals, got: %v", err.Error())
	}
	if len(journals) > 0 {
		t.Fatalf("Expected to fetch 0 journals, got: %v", len(journals))
	}
	j, err = js.Journal(ctx, jrnl.ID)
	if err == nil && j.Public == false {
		t.Fatalf("Expected to not find journal, but got it anyways")
	}

	// Test entry creation on other persons journal
	entry = &journal.Entry{
		Title:     "hello",
		JournalID: 1,
	}
	ntry, err = js.CreateEntry(ctx, entry)
	if err == nil {
		t.Fatalf("Expected error creating entry, got: %v", ntry)
	}

	// Checking access fetching entry in other users journal
	ntry, err = js.Entry(ctx, 1) // 1 is other users journal
	if err == nil {
		t.Fatalf("Expected error denied fetching entry, got: %v", ntry)
	}
	// Check if we can update an entry that is not ours
	ntry, err = js.UpdateEntry(ctx, &journal.Entry{ID: 1, Title: "world"})
	if err == nil {
		t.Fatalf("Expected access denied, got: %v", ntry)
	}

	//TODO Test Journals for user...
	journals, err = js.Journals(ctx, u.ID)
	if err != nil {
		t.Fatalf("Expected to get result, got %v", err.Error())
	}

	// Test create with no context
	_, err = js.Create(context.Background(), &journal.Journal{Title: title})
	if err == nil {
		t.Fatalf("Expected no err, got %v", err.Error())
	}
	_, err = js.MyJournals(context.Background())
	if err == nil {
		t.Fatalf("Expected no err, got %v", err.Error())
	}
	_, err = js.Journal(context.Background(), 1)
	if err == nil {
		t.Fatalf("Expected no err, got %v", err.Error())
	}
	_, err = js.CreateEntry(context.Background(), &journal.Entry{ID: 1, JournalID: 1})
	if err == nil {
		t.Fatalf("Expected no err, got %v", err.Error())
	}
	_, err = js.UpdateEntry(context.Background(), &journal.Entry{ID: 1})
	if err == nil {
		t.Fatalf("Expected no err, got %v", err.Error())
	}
	_, err = js.Entry(context.Background(), 1)
	if err == nil {
		t.Fatalf("Expected no err, got %v", err.Error())
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

func (jr *journalRepo) FindByID(ctx context.Context, id int64) (*journal.Journal, error) {
	for _, j := range jr.journals {
		if j.ID == id {
			return j, nil
		}
	}
	return nil, journal.ErrJournalNotExist
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

func (jr *journalRepo) AddEntry(c context.Context, entry *journal.Entry) (*journal.Entry, error) {
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

func (jr *journalRepo) FindEntryByID(ctx context.Context, id int64) (*journal.Entry, error) {
	for _, e := range jr.entries {
		if e.ID == id {
			return e, nil
		}
	}
	return nil, journal.ErrEntryNotExist
}

func (jr *journalRepo) FindAllEntries(ctx context.Context, journalID int64) ([]*journal.Entry, error) {
	entries := []*journal.Entry{}
	for _, e := range jr.entries {
		if e.JournalID == journalID {
			entries = append(entries, e)
		}
	}
	return entries, nil
}
