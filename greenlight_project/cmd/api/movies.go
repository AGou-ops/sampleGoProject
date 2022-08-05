package main

import (
	"fmt"
	"net/http"
	"time"

	"greenlight.agou-ops.cn/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var userInput struct {
		Title   string       `json:"title"`
		Year    int          `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}
	err := app.readJSON(w, r, &userInput)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	fmt.Fprintf(w, "%+v\n", userInput)
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		w.Write([]byte(err.Error()))
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Let's Go Further!",
		Year:      2020,
		Runtime:   99,
		Type:      []string{"Go", "GOlang", "web"},
		Version:   0,
	}
	if err := app.writeJson(w, *r, http.StatusOK, envelope{"movie": movie}, nil); err != nil {
		app.serverErrResponse(w, r, err)
	}
	w.Header().Set("content-type", "application/json")
}
