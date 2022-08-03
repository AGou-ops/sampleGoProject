package main

import (
	"fmt"
	"net/http"
	"time"

	"greenlight.agou-ops.cn/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "created a movie")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w,r)
		w.Write([]byte(err.Error()))
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Let's Go Further!",
		Year:      2020,
		Runtime:   0,
		Type:      []string{"Go", "GOlang", "web"},
		Version:   0,
	}
	if err := app.formatJson(w, http.StatusOK, movie, nil); err != nil {
		app.logger.Println(err.Error())
		http.Error(w, "format json error", http.StatusInternalServerError)
	}
	w.Header().Set("content-type", "application/json")
}
