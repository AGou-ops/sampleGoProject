package main

import (
	"errors"
	"net/http"
	"time"

	"greenlight.agou-ops.cn/internal/data"
	"greenlight.agou-ops.cn/internal/validator"
)

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.Password)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrResponse(w, r, err)
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrResponse(w, r, err)
		return
	}

	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	token, err := app.models.Tokens.New(user.ID, time.Hour*24, data.ScopeAuthentication)
	if err != nil {
		app.serverErrResponse(w, r, err)
		return
	}

	err = app.writeJson(w, *r, http.StatusCreated, envelope{"authentication_token": token}, nil)
	if err != nil {
		app.serverErrResponse(w, r, err)
	}
}

func (app *application) createPasswordResetTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateEmail(v, input.Email); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("email", "no matching record found!")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrResponse(w, r, err)
		}
		return
	}

	if !user.Activated {
		v.AddError("user", "user account must be activated!")
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	token, err := app.models.Tokens.New(user.ID, time.Minute*39, data.ScopePasswordReset)
	if err != nil {
		app.serverErrResponse(w, r, err)
		return
	}

	app.background(func() {
		data := map[string]interface{}{
			"passwordResetToken": token.PlainText,
		}
		err = app.dialer.Send(user.Email, "../../internal/mailer/templates/token_password_reset.tmpl", data)
		if err != nil {
			app.logger.PrintError(err, nil)
		}
	})
	env := envelope{"message": "an email will be sent to you containing password reset instructions"}
	if err = app.writeJson(w, *r, http.StatusAccepted, env, nil); err != nil {
		app.serverErrResponse(w, r, err)
	}
}

func (app *application) createActivationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.New()
	if data.ValidateEmail(v, input.Email); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Try to retrieve the corresponding user record for the email address. If it can't // be found, return an error message to the client.
	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("email", "no matching email address found")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrResponse(w, r, err)
		}
		return
	}
	// Return an error if the user has already been activated.
	if user.Activated {
		v.AddError("email", "user has already been activated")
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Otherwise, create a new activation token.
	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrResponse(w, r, err)
		return
	}
	// Email the user with their additional activation token.
	app.background(func() {
		data := map[string]interface{}{
			"activationToken": token.PlainText,
		}
		// Since email addresses MAY be case sensitive, notice that we are sending this // email using the address stored in our database for the user --- not to the // input.Email address provided by the client in this request.
		err = app.dialer.Send(user.Email, "../../internal/mailer/templates/token_activation.tmpl", data)
		if err != nil {
			app.logger.PrintError(err, nil)
		}
	})
	// Send a 202 Accepted response and confirmation message to the client.
	env := envelope{"message": "an email will be sent to you containing activation instructions"}
	err = app.writeJson(w, *r, http.StatusAccepted, env, nil)
	if err != nil {
		app.serverErrResponse(w, r, err)
	}
}
