package main

import (
	"net/http"

	"github.com/devinmiller/web-dev-with-go/context"
	"github.com/go-chi/chi/v5"
)

func (app *application) ClassRoutes() chi.Router {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Use(app.RequireUser)
		r.Post("/", app.PostClass)
	})

	return r
}

func (app *application) PostClass(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := context.User(r.Context())

	err = app.classService.CreateClass(r.Context(), r.FormValue("className"), user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
