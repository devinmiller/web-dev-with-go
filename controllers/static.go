package controllers

import (
	"net/http"

	"github.com/devinmiller/web-dev-with-go/views"
)

func StaticHandler(tmpl views.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	}
}

func RenderHandler(tm *views.TemplateManager, name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tm.Render(w, name, nil)
	}
}
