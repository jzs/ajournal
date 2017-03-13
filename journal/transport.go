package journal

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// SetupHandler sets up routes for the journal service
func SetupHandler(router *mux.Router, js Service) {
	router.Path("/users/{id}/journals").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idstr := vars["id"]
		id, err := strconv.ParseInt(idstr, 10, 64)
		if err != nil {
			JSONResp(w, nil, err)
			return
		}
		journals, err := js.Journals(r.Context(), id)
		err = JSONResp(w, journals, err)
		if err != nil {
			// TODO Handle error...
		}
	})

	router.Path("/journals").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		journals, err := js.MyJournals(r.Context())
		err = JSONResp(w, journals, err)
		if err != nil {
			// TODO Log this error or panic
		}
	})

	router.Path("/journals/{id}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idstr := vars["id"]
		id, err := strconv.ParseInt(idstr, 10, 64)
		if err != nil {
			JSONResp(w, nil, err)
			return
		}
		journal, err := js.Journal(r.Context(), id)
		JSONResp(w, journal, err)
	})

	router.Path("/journals/{id}/entries").Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idstr := vars["id"]
		id, err := strconv.ParseInt(idstr, 10, 64)
		if err != nil {
			JSONResp(w, nil, err)
			return
		}
		ntry := &Entry{}
		dec := json.NewDecoder(r.Body)
		err = dec.Decode(&ntry)
		if err != nil {
			JSONResp(w, nil, err)
			return
		}
		if ntry.JournalID != id {
			JSONResp(w, nil, errors.New("Mismatch between journal id's"))
			return
		}

		ntry, err = js.CreateEntry(r.Context(), ntry)
		JSONResp(w, ntry, err)
	})
	router.Path("/journals/{jid}/entries/{id}").Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idstr := vars["id"]
		id, err := strconv.ParseInt(idstr, 10, 64)
		if err != nil {
			JSONResp(w, nil, err)
			return
		}
		ntry := &Entry{}
		dec := json.NewDecoder(r.Body)
		err = dec.Decode(&ntry)
		if err != nil {
			JSONResp(w, nil, err)
			return
		}
		if ntry.ID != id {
			JSONResp(w, nil, errors.New("Mismatch between id's"))
			return
		}

		ntry, err = js.UpdateEntry(r.Context(), ntry)
		JSONResp(w, ntry, err)
	})

	router.Path("/journals/{jid}/entries/{id}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idstr := vars["id"]
		id, err := strconv.ParseInt(idstr, 10, 64)
		if err != nil {
			JSONResp(w, nil, err)
			return
		}
		ntry, err := js.Entry(r.Context(), id)
		JSONResp(w, ntry, err)
	})

	router.Path("/journals").Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse request if needed.
		journal := &Journal{}
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&journal)
		if err != nil {
			JSONResp(w, nil, err)
			return
		}

		// Call service method for data.
		journal, err = js.Create(r.Context(), journal)
		// Output response in proper format.
		err = JSONResp(w, journal, err)
		if err != nil {
			// TODO Log this error or panic
		}
	})
}

// Json response handling and error handling!

type jsonresp struct {
	Data   interface{}
	Status int64
	Error  string
}

// JSONResp formats responses in json
func JSONResp(w http.ResponseWriter, data interface{}, err error) error {
	w.Header().Set("content-type", "application/json")
	enc := json.NewEncoder(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := jsonresp{
			Status: http.StatusInternalServerError,
			Error:  err.Error(),
		}
		err = enc.Encode(resp)
		if err != nil {
			// Log this error or panic!
			return err
		}
		return nil
	}

	resp := jsonresp{
		Data:   data,
		Status: http.StatusOK,
	}
	err = enc.Encode(resp)
	if err != nil {
		return err
	}
	return nil
}
