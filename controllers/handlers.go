package controllers

import (
	"net/http"

	"github.com/devinmiller/web-dev-with-go/views"
)

func FormHandler(f func(map[string][]string) http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		handler := f(r.PostForm)

		handler(w, r)
	}
}

func TemplateHandler(tm *views.TemplateManager, name string, data interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := tm.RenderPage(w, name, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (app Application) ViewHandler(name string, f func(r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := app.tm.Template(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		f(r)
	}
}

type Application struct {
	tm *views.TemplateManager
}
