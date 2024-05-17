package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/redvant/lenslocked/context"
	"github.com/redvant/lenslocked/cookies"
	"github.com/redvant/lenslocked/errors"
	"github.com/redvant/lenslocked/models"
)

type Users struct {
	Templates struct {
		New            Template
		SignIn         Template
		ForgotPassword Template
		CheckYourEmail Template
		ResetPassword  Template
	}
	UserService          *models.UserService
	SessionService       *models.SessionService
	PasswordResetService *models.PasswordResetService
	EmailService         *models.EmailService
	HostAddress          string
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, r, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")
	user, err := u.UserService.Create(data.Email, data.Password)
	if err != nil {
		if errors.Is(err, models.ErrEmailTaken) {
			err = errors.Public(err,
				"That email address is already associated with an account.",
				http.StatusConflict,
			)
		}
		u.Templates.New.Execute(w, r, data, err)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		u.Templates.SignIn.Execute(w, r, data, err)
		return
	}
	cookies.SetCookie(w, cookies.CookieSession, session.Token)
	http.Redirect(w, r, "/galleries/", http.StatusFound)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.SignIn.Execute(w, r, data)
}

func (u Users) Authenticate(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")
	user, err := u.UserService.Authenticate(data.Email, data.Password)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			err = errors.Public(err,
				"An account could not be found with the email address provided.",
				http.StatusNotFound,
			)
		} else if errors.Is(err, models.ErrBadPassword) {
			err = errors.Public(err, "The provided password is incorrect.",
				http.StatusBadRequest,
			)
		}
		u.Templates.SignIn.Execute(w, r, data, err)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		u.Templates.SignIn.Execute(w, r, data, err)
		return
	}
	cookies.SetCookie(w, cookies.CookieSession, session.Token)
	http.Redirect(w, r, "/galleries/", http.StatusFound)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	fmt.Fprintf(w, "Current user: %s\n", user.Email)
}

func (u Users) SignOut(w http.ResponseWriter, r *http.Request) {
	token, err := cookies.ReadCookie(r, cookies.CookieSession)
	if err != nil {
		u.Templates.SignIn.Execute(w, r, nil, err)
		return
	}
	err = u.SessionService.Delete(token)
	if err != nil {
		u.Templates.SignIn.Execute(w, r, nil, err)
		return
	}
	cookies.DeleteCookie(w, cookies.CookieSession)
	http.Redirect(w, r, "/signin", http.StatusFound)
}

func (u Users) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.ForgotPassword.Execute(w, r, data)
}

func (u Users) ProcessForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	pwReset, err := u.PasswordResetService.Create(data.Email)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			err = errors.Public(err,
				"An account could not be found with the email address provided.",
				http.StatusNotFound,
			)
		}
		u.Templates.ForgotPassword.Execute(w, r, data, err)
		return
	}
	vals := url.Values{
		"token": {pwReset.Token},
	}
	resetURL := "http://" + u.HostAddress + "/reset-pw?" + vals.Encode()
	err = u.EmailService.ForgotPassword(data.Email, resetURL)
	if err != nil {
		u.Templates.ForgotPassword.Execute(w, r, data, err)
		return
	}
	u.Templates.CheckYourEmail.Execute(w, r, data)
}

func (u Users) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token string
	}
	data.Token = r.FormValue("token")
	u.Templates.ResetPassword.Execute(w, r, data)
}

func (u Users) ProcessResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token    string
		Password string
	}
	data.Token = r.FormValue("token")
	data.Password = r.FormValue("password")
	user, err := u.PasswordResetService.Consume(data.Token)
	if err != nil {
		if errors.Is(err, models.ErrInvalidPwResetToken) {
			err = errors.Public(err,
				`The provided password reset token is invalid.
				Please try generating a new token through the "Forgot your password?" page.`,
				http.StatusBadRequest,
			)
		} else if errors.Is(err, models.ErrExpiredPwResetToken) {
			err = errors.Public(err,
				`The provided password reset token has expired.
				Please try generating a new token through the "Forgot your password?" page.`,
				http.StatusBadRequest,
			)
		}
		u.Templates.ResetPassword.Execute(w, r, data, err)
		return
	}
	err = u.UserService.UpdatePassword(user.ID, data.Password)
	if err != nil {
		u.Templates.ResetPassword.Execute(w, r, data, err)
		return
	}

	// Sign the user in now that their password has been reset.
	// Any errors from this point should redirect the user to the sign in page.
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		u.Templates.SignIn.Execute(w, r, nil, err)
		return
	}
	cookies.SetCookie(w, cookies.CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}
