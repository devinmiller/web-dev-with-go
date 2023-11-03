package controllers

import (
	"net/http"

	"github.com/devinmiller/web-dev-with-go/models"
	models "github.com/devinmiller/web-dev-with-go/models/data"
	"github.com/devinmiller/web-dev-with-go/services"
	"github.com/devinmiller/web-dev-with-go/views"
	"github.com/go-chi/chi/v5"
)

type HomeController struct {
	views       *views.TemplateManager
	userService services.UserService
}

func NewHomeController(
	tm *views.TemplateManager,
	svc services.UserService) HomeController {

	c := HomeController{
		views:       tm,
		userService: svc,
	}

	return c
}

func (c HomeController) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", c.GetIndex)

	r.Get("/signin", c.GetSignIn)

	r.Get("/signup", c.GetSignUp)
	r.Post("/signup", c.PostSignUp)

	return r
}

func (c HomeController) GetIndex(w http.ResponseWriter, r *http.Request) {
	if err := c.views.RenderView(w, "home/index"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c HomeController) GetSignIn(w http.ResponseWriter, r *http.Request) {
	if err := c.views.RenderView(w, "home/signin"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c HomeController) PostSignIn(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (c HomeController) GetSignUp(w http.ResponseWriter, r *http.Request) {
	if err := c.views.RenderView(w, "home/signup"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c HomeController) PostSignUp(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := models.SignUpForm{
		FirstName: r.PostForm.Get("firstName"),
		LastName:  r.PostForm.Get("lastName"),
		Email:     r.PostForm.Get("email"),
		Password:  r.PostForm.Get("password"),
	}

	err = c.userService.SignUp(r.Context(), user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/signin", http.StatusFound)
}
