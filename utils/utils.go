package utils

import (
	"encoding/json"
	"net/http"
)

type ErrBadArgs struct {
	error
}

func NewErrBadArgs() error {
	return &ErrBadArgs{}
}

// JSONResp formats responses in json
func JSONResp(w http.ResponseWriter, data interface{}, err error) {
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
			panic(err)
		}
		return
	}

	resp := jsonresp{
		Data:   data,
		Status: http.StatusOK,
	}
	err = enc.Encode(resp)
	if err != nil {
		panic(err)
	}
	return
}

type jsonresp struct {
	Data   interface{}
	Status int64
	Error  string
}
