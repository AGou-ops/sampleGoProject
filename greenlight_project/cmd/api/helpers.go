package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 0, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

func (app *application) formatJson(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}
	json = append(json, '\n')

	for k,v := range headers {
		w.Header()[k] = v
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	w.Write(json)

	return nil
}
