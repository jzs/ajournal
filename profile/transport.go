package profile

import (
	"net/http"

	"bitbucket.org/sketchground/ajournal/utils"

	"github.com/gorilla/mux"
)

// SetupHandler sets up the handler routes for the user service
func SetupHandler(r *mux.Router, ps Service) {
	r.Path("/profile").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		profile, err := ps.Profile(r.Context())
		utils.JSONResp(w, profile, err)
	})
}
