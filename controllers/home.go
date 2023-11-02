package controllers

import (
	"net/http"
	"strings"

	models "github.com/devinmiller/web-dev-with-go/models/data"
	"github.com/devinmiller/web-dev-with-go/services"
	"github.com/devinmiller/web-dev-with-go/views"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

type HomeController struct {
	tm  *views.TemplateManager
	svc services.UserService
}

func NewHomeController(
	tm *views.TemplateManager,
	svc services.UserService) HomeController {

	c := HomeController{
		tm:  tm,
		svc: svc,
	}

	return c
}

func (c HomeController) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", c.GetIndex())

	r.Get("/signin", c.GetSignIn())
	r.Get("/signup", c.GetSignUp())

	r.Post("/signup", FormHandler(c.PostSignUp))

	r.Get("/contact", TemplateHandler(c.tm, "home/contact", nil))
	r.Get("/faq", c.FAQ(c.tm))

	return r
}

func (c HomeController) GetIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		return TemplateHandler(c.tm, "home/index", nil)
	}
}

func (c HomeController) GetSignIn() http.HandlerFunc {
	return TemplateHandler(c.tm, "home/signin", nil)
}

func (c HomeController) GetSignUp() http.HandlerFunc {
	return TemplateHandler(c.tm, "home/signup", nil)
}

func (c HomeController) PostSignUp(form map[string][]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hashedBytes, _ := bcrypt.GenerateFromPassword([]byte(r.PostForm.Get("password")), bcrypt.DefaultCost)

		user := models.User{
			FirstName:    r.PostForm.Get("firstName"),
			LastName:     r.PostForm.Get("lastName"),
			Email:        strings.ToLower(r.PostForm.Get("email")),
			PasswordHash: string(hashedBytes),
		}

		c.svc.SignUp(r.Context(), user)
	}
}

func (c HomeController) FAQ(tm *views.TemplateManager) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   string
	}{
		{
			Question: "Can you...",
			Answer:   "No",
		},
		{
			Question: "But...",
			Answer:   "Still no",
		},
	}

	return TemplateHandler(tm, "home/faq", questions)
}
