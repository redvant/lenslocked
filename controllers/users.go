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
	ServerAddress        string
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
				"That email address is already associated with an account.")
		}
		u.Templates.New.Execute(w, r, data, err)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		// TODO: Errors session creation
		u.Templates.SignIn.Execute(w, r, data, err)
		return
	}
	cookies.SetCookie(w, cookies.CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
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
		if errors.Is(err, models.ErrEmailNotFound) {
			err = errors.Public(err,
				"An account could not be found with the email address provided.")
		} else if errors.Is(err, models.ErrBadPassword) {
			err = errors.Public(err, "The provided password is incorrect.")
		}
		u.Templates.SignIn.Execute(w, r, data, err)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		// TODO: Errors session creation
		u.Templates.SignIn.Execute(w, r, data, err)
		return
	}
	cookies.SetCookie(w, cookies.CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	fmt.Fprintf(w, "Current user: %s\n", user.Email)
}

func (u Users) SignOut(w http.ResponseWriter, r *http.Request) {
	token, err := cookies.ReadCookie(r, cookies.CookieSession)
	if err != nil {
		// TODO: Errors reading session cookie
		u.Templates.SignIn.Execute(w, r, nil, err)
		return
	}
	err = u.SessionService.Delete(token)
	if err != nil {
		// TODO: Errors deleting session
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
		// TODO: Errors password reset creation
		u.Templates.ForgotPassword.Execute(w, r, data, err)
		return
	}
	vals := url.Values{
		"token": {pwReset.Token},
	}
	resetURL := "http://" + u.ServerAddress + "/reset-pw?" + vals.Encode()
	err = u.EmailService.ForgotPassword(data.Email, resetURL)
	if err != nil {
		// TODO: Errors send forgot password email
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
		// TODO: Errors consume password reset
		u.Templates.ResetPassword.Execute(w, r, data, err)
		return
	}
	err = u.UserService.UpdatePassword(user.ID, data.Password)
	if err != nil {
		// TODO: Errors update password
		u.Templates.ResetPassword.Execute(w, r, data, err)
		return
	}

	// Sign the user in now that their password has been reset.
	// Any errors from this point should redirect the user to the sign in page.
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		// TODO: Errors session create
		u.Templates.SignIn.Execute(w, r, nil, err)
		return
	}
	cookies.SetCookie(w, cookies.CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}
