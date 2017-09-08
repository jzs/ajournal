package profile

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sketchground/ajournal/utils"
	"github.com/sketchground/ajournal/utils/logger"

	"github.com/gorilla/mux"
)

// SetupHandler sets up the handler routes for the user service
func SetupHandler(r *mux.Router, ps Service, l logger.Logger) {
	// Handler for presenting a users profile
	r.Path("/users/{userid}/profile").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idstr := vars["userid"]
		id, err := strconv.ParseInt(idstr, 10, 64)
		if err != nil {
			utils.JSONResp(r.Context(), l, r, w, nil, err)
			return
		}
		profile, err := ps.UserProfile(r.Context(), id)
		utils.JSONResp(r.Context(), l, r, w, profile, err)
	})

	r.Path("/profile").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		profile, err := ps.Profile(r.Context())
		utils.JSONResp(r.Context(), l, r, w, profile, err)
	})

	r.Path("/profile").Methods("PUT").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prof := &Profile{}
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&prof)
		if err != nil {
			utils.JSONResp(r.Context(), l, r, w, nil, err)
			return
		}

		profile, err := ps.Create(r.Context(), prof)
		utils.JSONResp(r.Context(), l, r, w, profile, err)
	})

	r.Path("/profile").Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prof := &Profile{}
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&prof)
		if err != nil {
			utils.JSONResp(r.Context(), l, r, w, nil, err)
			return
		}

		profile, err := ps.UpdateProfile(r.Context(), prof)
		utils.JSONResp(r.Context(), l, r, w, profile, err)
	})

	r.Path("/profile/{userid}/shortname/{name}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate short name. Make sure it doesn't exist. Return possible alternatives.
		vars := mux.Vars(r)
		idstr := vars["userid"]
		name := vars["name"]
		id, err := strconv.ParseInt(idstr, 10, 64)
		if err != nil {
			utils.JSONResp(r.Context(), l, r, w, nil, err)
			return
		}

		valid := ps.ValidateShortName(r.Context(), id, name)
		utils.JSONResp(r.Context(), l, r, w, valid, nil)
	})

	// Handler for subscribing a plan
	r.Path("/profile/signup").Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sub := &Subscription{}
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&sub)
		if err != nil {
			utils.JSONResp(r.Context(), l, r, w, nil, utils.NewAPIError(nil, http.StatusBadRequest, "Invalid json"))
			return
		}

		// Handle signup form
		err = ps.Subscribe(r.Context(), sub)
		utils.JSONResp(r.Context(), l, r, w, nil, err)
	})
}
