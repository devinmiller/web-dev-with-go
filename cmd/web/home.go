package main

import (
	"net/http"

	"github.com/devinmiller/web-dev-with-go/models"
	"github.com/go-chi/chi/v5"
)

func (app *application) HomeRoutes() chi.Router {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Use(app.RequireUser)
		r.Get("/", app.View("home/index"))
		r.Get("/dashboard", app.View("home/dashboard"))
	})

	r.Get("/signin", app.View("home/signin"))
	r.Post("/signin", app.PostSignIn)

	r.Get("/signup", app.View("home/signup"))
	r.Post("/signup", app.PostSignUp)

	r.Post("/signout", app.PostSignOut)

	return r
}

func (app *application) PostSignIn(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	form := models.SignInForm{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	user, err := app.userService.SignIn(r.Context(), form)

	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	err = app.sessionService.SignIn(w, r, user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (app *application) PostSignUp(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	form := models.SignUpForm{
		FirstName: r.PostForm.Get("firstName"),
		LastName:  r.PostForm.Get("lastName"),
		Email:     r.PostForm.Get("email"),
		Password:  r.PostForm.Get("password"),
	}

	err = app.userService.SignUp(r.Context(), form)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/signin", http.StatusFound)
}

func (app *application) PostSignOut(w http.ResponseWriter, r *http.Request) {
	err := app.sessionService.SignOut(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
