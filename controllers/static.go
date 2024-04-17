package controllers

import (
	"html/template"
	"net/http"
)

func StaticHandler(tpl Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, nil)
	}
}

func FAQ(tpl Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Anwser   template.HTML
	}{
		{
			Question: "Is there a free version?",
			Anwser:   "Yes! We offer a free trial for 30 days on any paid plans.",
		},
		{
			Question: "What are your support hours?",
			Anwser:   "We have support staff answering emails 24/7, though response times may be a bit slower on weekends.",
		},
		{
			Question: "How do I contact support?",
			Anwser:   `Email us - <a class="text-indigo-600 hover:text-indigo-800 visited:text-purple-400 underline" href="mailto:support@lenslocked.com">support@lenslocked.com</a>`,
		},
		{
			Question: "Where is your office located?",
			Anwser:   "Our entire team is remote!",
		},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, questions)
	}
}

func Home(tpl Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		tpl.Execute(w, r, nil)
	}
}
