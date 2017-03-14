package profile

import (
	"net/http"

	"bitbucket.org/sketchground/journal/utils"

	"github.com/gorilla/mux"
)

// SetupHandler sets up the handler routes for the user service
func SetupHandler(r *mux.Router, ps Service) {
	// Create user
	r.Path("/profile").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		profile, err := ps.Profile(r.Context())
		utils.JSONResp(w, profile, err)
	})
}
