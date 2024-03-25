package main

import (
	"fmt"
	"net/http"
)

func homeHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<body style=\"font-family: courier, monospace;\"><h1>Welcome to my fantastic site!</h1></body>")
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
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", homeHandler)
	mux.HandleFunc("GET /contact", contactHandler)
	mux.HandleFunc("GET /faq", faqHandler)
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", mux)
}
