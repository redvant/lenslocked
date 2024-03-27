package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/redvant/lenslocked/middleware"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tplPath := filepath.Join("templates", "home.gohtml")
	t, err := template.ParseFiles(tplPath)
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

func contactHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<body style=\"font-family: courier, monospace;\"><h1>Contact Page</h1><p>To get in touch, email me at <a href=\"mailto:redvant@outlook.com\">redvant@outlook.com</a>.</p></body>")
}

func faqHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<body style="font-family: courier, monospace;">
	<h1>FAQ page</h1>
	<ul>
		<li>
			<p style="font-weight: bold;">Is there a free version</p>
			Yes! We offer a free trial for 30 days on any paid plans.
		</li>
		<li>
			<p style="font-weight: bold;">What are your support hours?</p>
			We have support staff answering emails 24/7, though response
			times may be a bit slower on weekends.
		</li>
		<li>
			<p style="font-weight: bold;">How do I contact support?</p>
			Email us - <a href="mailto:support@lenslocked.com">support@lenslocked.com</a>
		</li>
	</ul>
	</body>
	`)
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
