package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/redvant/lenslocked/middleware"
	"github.com/redvant/lenslocked/views"
)

func executeTemplate(w http.ResponseWriter, filepath string, data interface{}) {
	t, err := views.Parse(filepath)
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(w, "There was an error parsing the template.", http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	tplPath := filepath.Join("templates", "home.gohtml")
	executeTemplate(w, tplPath, nil)
}

func contactHandler(w http.ResponseWriter, _ *http.Request) {
	tplPath := filepath.Join("templates", "contact.gohtml")
	executeTemplate(w, tplPath, nil)
}

func faqHandler(w http.ResponseWriter, _ *http.Request) {
	tplPath := filepath.Join("templates", "faq.gohtml")
	executeTemplate(w, tplPath, nil)
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("GET /", homeHandler)
	router.HandleFunc("GET /contact", contactHandler)
	router.HandleFunc("GET /faq", faqHandler)

	server := http.Server{
		Addr:    ":3000",
		Handler: middleware.Logging(router),
	}

	fmt.Println("Starting the server on :3000...")
	server.ListenAndServe()
}
