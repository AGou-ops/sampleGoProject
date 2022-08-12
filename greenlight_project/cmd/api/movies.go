package main

import (
	"errors"
	"fmt"
	"net/http"

	"greenlight.agou-ops.cn/internal/data"
	"greenlight.agou-ops.cn/internal/validator"
)

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		w.Write([]byte(err.Error()))
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrorRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrResponse(w, r, err)
		}
		return
	}

	// movie := data.Movie{
	// 	ID:        id,
	// 	CreatedAt: time.Now(),
	// 	Title:     "Let's Go Further!",
	// 	Year:      2020,
	// 	Runtime:   99,
	// 	Genres:    []string{"Go", "GOlang", "web"},
	// 	Version:   0,
	// }
	if err := app.writeJson(w, *r, http.StatusOK, envelope{"movie": movie}, nil); err != nil {
		app.serverErrResponse(w, r, err)
	}
	w.Header().Set("content-type", "application/json")
}

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var userInput struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}
	err := app.readJSON(w, r, &userInput)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie := &data.Movie{
		Title:   userInput.Title,
		Year:    userInput.Year,
		Runtime: userInput.Runtime,
		Genres:  userInput.Genres,
	}

	v := validator.New()

	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
	}

	if err = app.models.Movies.Insert(movie); err != nil {
		app.serverErrResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	if err = app.writeJson(w, *r, http.StatusCreated, envelope{"movie": movie}, headers); err != nil {
		app.serverErrResponse(w, r, err)
	}

	fmt.Fprintf(w, "%+v\n", userInput)
}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrorRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrResponse(w, r, err)
		}
		return
	}
	var userInput struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}
	err = app.readJSON(w, r, &userInput)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie.Title = userInput.Title
	movie.Year = userInput.Year
	movie.Runtime = userInput.Runtime
	movie.Genres = userInput.Genres

	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.serverErrResponse(w, r, err)
	}

	if err = app.models.Movies.Update(movie); err != nil {
		app.serverErrResponse(w, r, err)
		return
	}

	if err = app.writeJson(w, *r, http.StatusOK, envelope{"movie": movie}, nil); err != nil {
		app.serverErrResponse(w, r, err)
	}
}
