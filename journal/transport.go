package journal

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sketchground/ajournal/utils"
	"github.com/sketchground/ajournal/utils/logger"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// SetupHandler sets up routes for the journal service
func SetupHandler(router *mux.Router, js Service, l logger.Logger) {
	router.Path("/users/{id}/journals").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idstr := vars["id"]
		id, err := strconv.ParseInt(idstr, 10, 64)
		if err != nil {
			utils.JSONResp(r.Context(), l, r, w, nil, err)
			return
		}
		journals, err := js.Journals(r.Context(), id)
		utils.JSONResp(r.Context(), l, r, w, journals, err)
	})

	router.Path("/journals").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		journals, err := js.MyJournals(r.Context())
		utils.JSONResp(r.Context(), l, r, w, journals, err)
	})

	router.Path("/journals/{id}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idstr := vars["id"]
		id, err := strconv.ParseInt(idstr, 10, 64)
		if err != nil {
			utils.JSONResp(r.Context(), l, r, w, nil, err)
			return
		}
		journal, err := js.Journal(r.Context(), id)
		utils.JSONResp(r.Context(), l, r, w, journal, err)
	})

	router.Path("/journals/{id}/entries").Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idstr := vars["id"]
		id, err := strconv.ParseInt(idstr, 10, 64)
		if err != nil {
			utils.JSONResp(r.Context(), l, r, w, nil, errors.Wrap(err, "Router"))
			return
		}
		ntry := &Entry{}
		dec := json.NewDecoder(r.Body)
		err = dec.Decode(&ntry)
		if err != nil {
			utils.JSONResp(r.Context(), l, r, w, nil, errors.Wrap(err, "Router: err decoding json"))
			return
		}
		if ntry.JournalID != id {
			utils.JSONResp(r.Context(), l, r, w, nil, utils.NewAPIError(nil, http.StatusBadRequest, "Mismatch between journal id's"))
			return
		}

		ntry, err = js.CreateEntry(r.Context(), ntry)
		utils.JSONResp(r.Context(), l, r, w, ntry, err)
	})
	router.Path("/journals/{jid}/entries/{id}").Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idstr := vars["id"]
		id, err := strconv.ParseInt(idstr, 10, 64)
		if err != nil {
			utils.JSONResp(r.Context(), l, r, w, nil, err)
			return
		}
		ntry := &Entry{}
		dec := json.NewDecoder(r.Body)
		err = dec.Decode(&ntry)
		if err != nil {
			utils.JSONResp(r.Context(), l, r, w, nil, err)
			return
		}
		if ntry.ID != id {
			utils.JSONResp(r.Context(), l, r, w, nil, utils.NewAPIError(nil, http.StatusBadRequest, "Mismatch between id's"))
			return
		}

		ntry, err = js.UpdateEntry(r.Context(), ntry)
		utils.JSONResp(r.Context(), l, r, w, ntry, err)
	})

	router.Path("/journals/{jid}/entries/{id}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idstr := vars["id"]
		id, err := strconv.ParseInt(idstr, 10, 64)
		if err != nil {
			utils.JSONResp(r.Context(), l, r, w, nil, err)
			return
		}
		ntry, err := js.Entry(r.Context(), id)
		utils.JSONResp(r.Context(), l, r, w, ntry, err)
	})

	router.Path("/journals").Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse request if needed.
		journal := &Journal{}
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&journal)
		if err != nil {
			utils.JSONResp(r.Context(), l, r, w, nil, err)
			return
		}

		// Call service method for data.
		journal, err = js.Create(r.Context(), journal)
		// Output response in proper format.
		utils.JSONResp(r.Context(), l, r, w, journal, err)
	})
}
