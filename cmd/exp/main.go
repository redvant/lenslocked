package main

import (
	"html/template"
	"os"
)

func main() {
	t := template.Must(template.ParseFiles("hello.gohtml"))

	user := struct {
		Name string
		Bio  string
		Age  int
	}{
		Name: "Roman",
		Bio:  `<script>alert("Haha, you have been h4x0r3d!");</script>`,
		Age:  123,
	}

	err := t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}
}
