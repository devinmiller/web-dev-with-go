package main

import (
	"net/http"

	"github.com/devinmiller/web-dev-with-go/models"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

func (app *application) HomeRoutes() chi.Router {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Use(app.RequireUser)
		r.Get("/", app.View("home/index"))
	})

	//r.Get("/", app.View("home/index"))

	r.Get("/signin", app.GetSignIn)
	r.Post("/signin", app.PostSignIn)

	r.Get("/signup", app.GetSignUp)
	r.Post("/signup", app.PostSignUp)

	return r
}

func (app *application) GetSignIn(w http.ResponseWriter, r *http.Request) {
	// TODO: Improve CSRF handling
	templateData := map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	}

	// TODO:
	if err := app.views.RenderPage(w, "home/signin", templateData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

func (app *application) GetSignUp(w http.ResponseWriter, r *http.Request) {
	// TODO: Improve CSRF handling
	templateData := map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	}

	if err := app.views.RenderPage(w, "home/signup", templateData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
