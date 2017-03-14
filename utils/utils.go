package utils

import (
	"encoding/json"
	"net/http"
)

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

type jsonresp struct {
	Data   interface{}
	Status int64
	Error  string
}
