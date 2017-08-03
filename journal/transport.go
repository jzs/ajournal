package journal

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/sketchground/ajournal/blob"
	"github.com/sketchground/ajournal/common"
	"github.com/sketchground/ajournal/user"
	"github.com/sketchground/ajournal/utils"
	"github.com/sketchground/ajournal/utils/logger"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// SetupHandler sets up routes for the journal service
func SetupHandler(router *mux.Router, js Service, bs blob.Service, l logger.Logger) {
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

	router.Path("/journals/{id}/entries").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idstr := vars["id"]
		id, err := strconv.ParseInt(idstr, 10, 64)
		if err != nil {
			utils.JSONResp(r.Context(), l, r, w, nil, errors.Wrap(err, "Router"))
			return
		}
		args := common.ParsePagination(r)
		entries, err := js.Entries(r.Context(), id, args)
		utils.JSONResp(r.Context(), l, r, w, entries, err)

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

	router.Path("/journals/{id}/blobs").Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idstr := vars["id"]
		id, err := strconv.ParseInt(idstr, 10, 64)
		if err != nil {
			utils.JSONResp(r.Context(), l, r, w, nil, err)
			return
		}

		usr := user.FromContext(r.Context())
		if usr == nil {
			utils.JSONResp(r.Context(), l, r, w, nil, utils.NewAPIError(nil, http.StatusForbidden, "Unauthorized"))
			return
		}

		//TODO: Create/update a new blob on a given hash
		// Check if we have access to the journal.
		j, err := js.Journal(r.Context(), id)
		if err != nil {
			utils.JSONResp(r.Context(), l, r, w, nil, err)
			return
		}
		if j.UserID != usr.ID {
			utils.JSONResp(r.Context(), l, r, w, nil, utils.NewAPIError(nil, http.StatusForbidden, "Unauthorized"))
			return
		}

		// Create blob
		file, header, err := r.FormFile("blobs")
		if err != nil {
			utils.JSONResp(r.Context(), l, r, w, nil, utils.NewAPIError(err, http.StatusBadRequest, "File in multipart form does not exist"))
			return
		}
		hash := md5.New()
		if _, err := io.Copy(hash, file); err != nil {
			log.Fatal(err)
		}
		f, err := bs.Create(fmt.Sprintf("journals/%v/blobs/%x", id, hash.Sum(nil)), header.Header["Content-Type"][0], file)
		utils.JSONResp(r.Context(), l, r, w, f, err)

	})

	router.Path("/journals/{id}/blobs").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idstr := vars["id"]
		id, err := strconv.ParseInt(idstr, 10, 64)
		if err != nil {
			utils.JSONResp(r.Context(), l, r, w, nil, err)
			return
		}

		// Check if we have access to the journal.
		if _, err := js.Journal(r.Context(), id); err != nil {
			utils.JSONResp(r.Context(), l, r, w, nil, err)
			return
		}
		files, err := bs.List(fmt.Sprintf("journals/%v/blobs/", id))
		utils.JSONResp(r.Context(), l, r, w, files, err)
	})

	router.Path("/journals/{id}/blobs/{hash}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idstr := vars["id"]
		hash := vars["hash"]
		id, err := strconv.ParseInt(idstr, 10, 64)
		if err != nil {
			utils.JSONResp(r.Context(), l, r, w, nil, err)
			return
		}
		details, err := bs.Details(fmt.Sprintf("journals/%v/blobs/%v", id, hash))
		if err != nil {
			utils.JSONResp(r.Context(), l, r, w, nil, err)
			return
		}
		w.Header().Set("Content-Type", details.MIMEType)
		_, err = io.Copy(w, details.Reader)
		if err != nil {
			l.Printf(r.Context(), err.Error())
			return
		}
	})
}
