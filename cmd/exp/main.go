package main

import (
	"html/template"
	"os"
)

type User struct {
	Name string
}

func main() {
	t, err := template.ParseFiles("hello.html")

	if err != nil {
		panic(err)
	}

	user := User{Name: "Devin Miller"}

	err = t.Execute(os.Stdout, user)

	if err != nil {
		panic(err)
	}
}
