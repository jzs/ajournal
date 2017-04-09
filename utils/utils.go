package utils

import (
	"context"
	"encoding/json"
	"net/http"

	"bitbucket.org/sketchground/ajournal/utils/logger"
)

type ErrBadArgs struct {
	error
}

func NewErrBadArgs() error {
	return &ErrBadArgs{}
}

// JSONResp formats responses in json
func JSONResp(ctx context.Context, l logger.Logger, w http.ResponseWriter, data interface{}, err error) {
	w.Header().Set("content-type", "application/json")
	enc := json.NewEncoder(w)
	if err != nil {
		// TODO: Build out this to check for error kinds.
		l.Error(ctx, err)
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
