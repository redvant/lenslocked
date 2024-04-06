package main

import (
	"fmt"
	"net/http"

	"github.com/redvant/lenslocked/controllers"
	"github.com/redvant/lenslocked/middleware"
	"github.com/redvant/lenslocked/models"
	"github.com/redvant/lenslocked/templates"
	"github.com/redvant/lenslocked/views"
)

func main() {
	router := http.NewServeMux()

	tplHome := views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "home.gohtml"))
	// Add notFound check to StaticHandler for "/"
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		controllers.StaticHandler(tplHome).ServeHTTP(w, r)
	})

	router.HandleFunc("GET /contact", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "contact.gohtml"))))

	router.HandleFunc("GET /faq", controllers.FAQ(
		views.Must(views.ParseFS(templates.FS, "tailwind.gohtml", "faq.gohtml"))))

	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userService := models.UserService{
		DB: db,
	}
	usersC := controllers.Users{
		UserService: &userService,
	}
	usersC.Templates.New = views.Must(views.ParseFS(
		templates.FS, "tailwind.gohtml", "signup.gohtml",
	))
	usersC.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS, "tailwind.gohtml", "signin.gohtml",
	))
	router.HandleFunc("GET /signup", usersC.New)
	router.HandleFunc("POST /users", usersC.Create)
	router.HandleFunc("GET /signin", usersC.SignIn)

	server := http.Server{
		Addr:    ":3000",
		Handler: middleware.Logging(router),
	}

	fmt.Println("Starting the server on :3000...")
	server.ListenAndServe()
}
