package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/redvant/lenslocked/controllers"
	"github.com/redvant/lenslocked/middleware"
	"github.com/redvant/lenslocked/migrations"
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

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}
	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
	}
	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
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
	router.HandleFunc("POST /signin", usersC.Authenticate)
	router.HandleFunc("POST /signout", usersC.SignOut)
	router.HandleFunc("GET /users/me", usersC.CurrentUser)

	csrfKey := "qaKGjjr8CPhMUqTjLXU6oJ8PsS45UcgQ"
	csrfMw := csrf.Protect([]byte(csrfKey),
		csrf.Secure(false), // TODO: Remove this for PROD
	)

	server := http.Server{
		Addr:    ":3000",
		Handler: middleware.Logging(csrfMw(router)),
	}

	fmt.Println("Starting the server on :3000...")
	server.ListenAndServe()
}
