package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/redvant/lenslocked/controllers"
	"github.com/redvant/lenslocked/middleware"
	"github.com/redvant/lenslocked/views"
)

func main() {
	router := http.NewServeMux()

	tplHome := views.Must(views.Parse(filepath.Join("templates", "home.gohtml")))
	// Add notFound check to StaticHandler for "/"
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		controllers.StaticHandler(tplHome).ServeHTTP(w, r)
	})

	router.HandleFunc("GET /contact", controllers.StaticHandler(
		views.Must(views.Parse(filepath.Join("templates", "contact.gohtml")))))

	router.HandleFunc("GET /faq", controllers.StaticHandler(
		views.Must(views.Parse(filepath.Join("templates", "faq.gohtml")))))

	server := http.Server{
		Addr:    ":3000",
		Handler: middleware.Logging(router),
	}

	fmt.Println("Starting the server on :3000...")
	server.ListenAndServe()
}
