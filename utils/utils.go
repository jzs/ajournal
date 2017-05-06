package utils

import (
	"context"
	"encoding/json"
	"net/http"

	"bitbucket.org/sketchground/ajournal/utils/logger"
)

// JSONResp formats responses in json
func JSONResp(ctx context.Context, l logger.Logger, w http.ResponseWriter, data interface{}, err error) {
	w.Header().Set("content-type", "application/json")
	enc := json.NewEncoder(w)
	if err != nil {
		var resp jsonresp

		switch err.(type) {
		case APIError:
			apierr := err.(APIError)
			w.WriteHeader(apierr.Status)
			resp = jsonresp{
				Status: apierr.Status,
				Error:  apierr.Desc,
			}
			l.Printf(ctx, err.Error())
			l.Print(ctx, apierr.Cause())
		default:
			w.WriteHeader(http.StatusInternalServerError)
			resp = jsonresp{
				Status: http.StatusInternalServerError,
				Error:  "Internal server error",
			}
			// Log the error...
			l.Error(ctx, err)
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
	Status int
	Error  string
}
