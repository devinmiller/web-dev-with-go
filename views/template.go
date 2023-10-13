package views

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type Template struct {
	htmlTmpl *template.Template
}

func Must(t *Template, err error) Template {
	if err != nil {
		panic(err)
	}

	return *t
}

func Parse(filename string) (*Template, error) {
	filepath := filepath.Join("templates", filename)

	tpl, err := template.ParseFiles(filepath)

	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}

	return &Template{
		htmlTmpl: tpl,
	}, nil
}

func (t Template) Execute(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := t.htmlTmpl.Execute(w, data)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing the template", http.StatusInternalServerError)
		return
	}
}
