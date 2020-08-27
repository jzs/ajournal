package blob

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jzs/ajournal/utils"
	"github.com/jzs/ajournal/utils/logger"
)

// SetupHandler sets up the blob endpoint
func SetupHandler(router *mux.Router, bs Service, l logger.Logger) {

	router.Path("/blobs/{key:.*}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]

		details, err := bs.Details(key)
		if err != nil {
			utils.JSONResp(r.Context(), l, r, w, nil, err)
			return
		}
		w.Header().Add("Content-Type", details.MIMEType)
		w.WriteHeader(http.StatusOK)

		// Fetch image and stream it...
		reader, err := bs.Value(key)
		if err != nil {
			utils.JSONResp(r.Context(), l, r, w, nil, err)
			return
		}

		// make a buffer to keep chunks that are read
		buf := make([]byte, 1024)
		for {
			// read a chunk
			n, err := reader.Read(buf)
			if err != nil && err != io.EOF {
				panic(err)
			}
			if n == 0 {
				break
			}

			// write a chunk
			if _, err := w.Write(buf[:n]); err != nil {
				panic(err)
			}
		}
	})

	router.Path("/blobs/details/{key}").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]
		bd, err := bs.Details(key)
		if err != nil {
			utils.JSONResp(r.Context(), l, r, w, nil, err)
			return
		}
		utils.JSONResp(r.Context(), l, r, w, bd, nil)
	})

	router.Path("/blobs/{key}").Methods("PUT").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		key := vars["key"]

		file, header, err := r.FormFile("blob")
		if err != nil {
			utils.JSONResp(r.Context(), l, r, w, nil, err)
		}
		lr := io.LimitReader(file, 10000000)
		bd, err := bs.Create(key, header.Header[""][0], lr)
		if err != nil {
			utils.JSONResp(r.Context(), l, r, w, nil, err)
		}
		utils.JSONResp(r.Context(), l, r, w, bd, nil)
	})
}
