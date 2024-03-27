package main

import (
	"html/template"
	"os"
)

func main() {
	t := template.Must(template.ParseFiles("hello.gohtml"))

	user := struct {
		Name    string
		Bio     string
		Age     int
		IsAdmin bool
	}{
		Name:    "Roman",
		Bio:     `<script>alert("Haha, you have been h4x0r3d!");</script>`,
		Age:     30,
		IsAdmin: false,
	}

	err := t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}
}
