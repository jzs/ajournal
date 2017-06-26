package profile_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sketchground/ajournal/profile"
	"github.com/sketchground/ajournal/utils/logger"
)

func TestTransport(t *testing.T) {
	m := mux.NewRouter()
	pr := NewInmemRepo()
	sr := NewInmemSubRepo()
	ps := profile.NewService(pr, sr)
	profile.SetupHandler(m, ps, logger.NewTestLogger())

	posts := []struct {
		URL      string
		Code     int
		Type     string
		PostBody string
	}{
		{
			URL:      "/profile",
			Code:     http.StatusForbidden,
			Type:     "PUT",
			PostBody: "{}",
		},
		{
			URL:  "/profile",
			Code: http.StatusForbidden,
			Type: "GET",
		},
		{
			URL:      "/profile",
			Code:     http.StatusForbidden,
			Type:     "POST",
			PostBody: "{}",
		},
	}

	for _, p := range posts {
		var req *http.Request
		switch p.Type {
		case "GET":
			req, _ = http.NewRequest(p.Type, p.URL, nil)
		case "POST":
			req, _ = http.NewRequest(p.Type, p.URL, strings.NewReader(p.PostBody))
		case "PUT":
			req, _ = http.NewRequest(p.Type, p.URL, strings.NewReader(p.PostBody))
		default:
			req, _ = http.NewRequest(p.Type, p.URL, nil)
		}

		rw := httptest.NewRecorder()
		m.ServeHTTP(rw, req)
		if rw.Code != p.Code {
			t.Errorf("Expected %v on url %v, got %v", p.Code, p.URL, rw.Code)
		}
	}
}
