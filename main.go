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

	tpl, err := views.Parse(filepath.Join("templates", "home.gohtml"))
	if err != nil {
		panic(err)
	}

	// Add notFound check to StaticHandler for "/"
	router.HandleFunc("GET /", func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				http.NotFound(w, r)
				return
			}
			next.ServeHTTP(w, r)
		}
	}(controllers.StaticHandler(tpl)))

	tpl, err = views.Parse(filepath.Join("templates", "contact.gohtml"))
	if err != nil {
		panic(err)
	}
	router.HandleFunc("GET /contact", controllers.StaticHandler(tpl))

	tpl, err = views.Parse(filepath.Join("templates", "faq.gohtml"))
	if err != nil {
		panic(err)
	}
	router.HandleFunc("GET /faq", controllers.StaticHandler(tpl))

	server := http.Server{
		Addr:    ":3000",
		Handler: middleware.Logging(router),
	}

	fmt.Println("Starting the server on :3000...")
	server.ListenAndServe()
}
