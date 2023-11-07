package controllers

import (
	"net/http"

	"github.com/devinmiller/web-dev-with-go/models"
	"github.com/devinmiller/web-dev-with-go/services"
	"github.com/devinmiller/web-dev-with-go/views"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

type HomeController struct {
	views          *views.TemplateManager
	userService    *services.UserService
	sessionService *services.SessionService
}

func NewHomeController(
	views *views.TemplateManager,
	userService *services.UserService,
	sessionService *services.SessionService) HomeController {

	c := HomeController{
		views:          views,
		userService:    userService,
		sessionService: sessionService,
	}

	return c
}

func (c *HomeController) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", c.GetIndex)

	r.Get("/signin", c.GetSignIn)
	r.Post("/signin", c.PostSignIn)

	r.Get("/signup", c.GetSignUp)
	r.Post("/signup", c.PostSignUp)

	return r
}

func (c *HomeController) GetIndex(w http.ResponseWriter, r *http.Request) {
	if err := c.views.RenderView(w, "home/index"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *HomeController) GetSignIn(w http.ResponseWriter, r *http.Request) {
	// TODO: Improve CSRF handling
	templateData := map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	}

	// TODO:
	if err := c.views.RenderPage(w, "home/signin", templateData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *HomeController) PostSignIn(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	form := models.SignInForm{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	user, err := c.userService.SignIn(r.Context(), form)

	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	cookie := http.Cookie{
		Name:     "flashy",
		Value:    user.Email,
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
}

func (c *HomeController) GetSignUp(w http.ResponseWriter, r *http.Request) {
	// TODO: Improve CSRF handling
	templateData := map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	}

	if err := c.views.RenderPage(w, "home/signup", templateData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *HomeController) PostSignUp(w http.ResponseWriter, r *http.Request) {
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

	err = c.userService.SignUp(r.Context(), form)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/signin", http.StatusFound)
}
