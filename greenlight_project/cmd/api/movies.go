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
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrResponse(w, r, err)
		}
		return
	}

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

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	if err = app.models.Movies.Delete(id); err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrResponse(w, r, err)
		}
		return
	}

	if err = app.writeJson(w, *r, http.StatusOK, envelope{"message": "Delete Movie Succesfully"}, nil); err != nil {
		app.serverErrResponse(w, r, err)
	}
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
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrResponse(w, r, err)
		}
		return
	}
	var userInput struct {
		Title   *string       `json:"title"`
		Year    *int32        `json:"year"`
		Runtime *data.Runtime `json:"runtime"`
		Genres  []string      `json:"genres"`
	}
	err = app.readJSON(w, r, &userInput)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if userInput.Title != nil {
		movie.Title = *userInput.Title
	}
	if userInput.Year != nil {
		movie.Year = *userInput.Year
	}
	if userInput.Runtime != nil {
		movie.Runtime = *userInput.Runtime
	}
	if userInput.Genres != nil {
		movie.Genres = userInput.Genres
	}

	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err = app.models.Movies.Update(movie); err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflitResponse(w, r)
		default:
			app.serverErrResponse(w, r, err)
		}
		return
	}

	if err = app.writeJson(w, *r, http.StatusOK, envelope{"movie": movie}, nil); err != nil {
		app.serverErrResponse(w, r, err)
	}
}

func (app *application) listMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Genres []string
		data.Filters
	}
	v := validator.New()
	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Genres = app.readCSV(qs, "genres", []string{})
	input.Page = app.readInt(qs, "page", 1, v)
	input.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Sort = app.readString(qs, "sort", "id")
	input.SortSafeList = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	movies, metadata, err := app.models.Movies.GetAll(input.Title, input.Genres, input.Filters)
	if err != nil {
		app.serverErrResponse(w, r, err)
		return
	}

	err = app.writeJson(w, *r, http.StatusOK, envelope{"metadata": metadata, "movies": movies}, nil)
	if err != nil {
		app.serverErrResponse(w, r, err)
	}
}
