package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/redvant/lenslocked/middleware"
)

func executeTemplate(w http.ResponseWriter, filepath string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.ParseFiles(filepath)
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(w, "There was an error parsing the template.", http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	tplPath := filepath.Join("templates", "home.gohtml")
	executeTemplate(w, tplPath)
}

func contactHandler(w http.ResponseWriter, _ *http.Request) {
	tplPath := filepath.Join("templates", "contact.gohtml")
	executeTemplate(w, tplPath)
}

func faqHandler(w http.ResponseWriter, _ *http.Request) {
	tplPath := filepath.Join("templates", "faq.gohtml")
	executeTemplate(w, tplPath)
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
