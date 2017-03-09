package journal

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// SetupHandler sets up routes for the journal service
func SetupHandler(router *mux.Router, js Service) {

	router.Path("/journals").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		journals, err := js.Journals(r.Context())
		err = JSONResp(w, journals, err)
		if err != nil {
			// TODO Log this error or panic
		}
	})

	router.Path("/journals").Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Post on journals")
		// Parse request if needed.
		// Call service method for data.
		// Output response in proper format.
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
