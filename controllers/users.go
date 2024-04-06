package controllers

import (
	"fmt"
	"net/http"

	"github.com/redvant/lenslocked/models"
)

type Users struct {
	Templates struct {
		New Template
	}
	UserService *models.UserService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, r.FormValue("email"))
	fmt.Fprintln(w, r.FormValue("password"))
}
