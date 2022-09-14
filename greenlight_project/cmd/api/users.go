package main

import (
	"errors"
	"net/http"

	"greenlight.agou-ops.cn/internal/data"
	"greenlight.agou-ops.cn/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err := user.Password.Set(input.Password)
	if err != nil {
		app.serverErrResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrResponse(w, r, err)
		}
		return
	}

	app.background(func() {
		err = app.dialer.Send(user.Email, "../../internal/mailer/templates/user_welcome.tmpl", user)
		if err != nil {
			app.serverErrResponse(w, r, err)
		}
	})

	if err = app.writeJson(w, *r, http.StatusCreated, envelope{"User": user}, nil); err != nil {
		app.serverErrResponse(w, r, err)
	}
}
