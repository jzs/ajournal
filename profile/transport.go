package profile

import (
	"encoding/json"
	"net/http"

	"bitbucket.org/sketchground/ajournal/utils"

	"github.com/gorilla/mux"
)

// SetupHandler sets up the handler routes for the user service
func SetupHandler(r *mux.Router, ps Service) {
	// Handler for presenting a users profile
	r.Path("/profile").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		profile, err := ps.Profile(r.Context())
		utils.JSONResp(w, profile, err)
	})

	r.Path("/profile").Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prof := &Profile{}
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&prof)
		if err != nil {
			utils.JSONResp(w, nil, err)
			return
		}

		profile, err := ps.UpdateProfile(r.Context(), prof)
		utils.JSONResp(w, profile, err)
	})

	// Handler for subscribing a plan
	r.Path("/profile/signup").Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sub := &Subscription{}
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&sub)
		if err != nil {
			utils.JSONResp(w, nil, err)
			return
		}

		// Handle signup form
		err = ps.Subscribe(r.Context(), sub)
		utils.JSONResp(w, nil, err)
	})
}
