package profile_test

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sketchground/ajournal/profile"
	"github.com/sketchground/ajournal/utils/logger"
	"github.com/sketchground/ajournal/utils/testhelpers"
)

func TestTransport(t *testing.T) {
	m := mux.NewRouter()
	pr := NewInmemRepo()
	sr := NewInmemSubRepo()
	ps := profile.NewService(pr, sr)
	profile.SetupHandler(m, ps, logger.NewTestLogger())

	posts := []testhelpers.Request{
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
	testhelpers.PerformRequests(t, m, posts)
}
