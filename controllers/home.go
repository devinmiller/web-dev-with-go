package controllers

import (
	"net/http"

	"github.com/devinmiller/web-dev-with-go/templates"
	"github.com/devinmiller/web-dev-with-go/views"
	"github.com/go-chi/chi/v5"
)

type HomeController struct {
	tm *views.TemplateManager
}

func NewHomeController() (controller HomeController) {
	tm, err := views.NewTemplateManager(templates.FS, ".", "layout", ".html")

	if err != nil {
		panic(err)
	}

	c := HomeController{
		tm: tm,
	}

	return c
}

func (c HomeController) Routes(tm *views.TemplateManager) chi.Router {
	r := chi.NewRouter()

	r.Get("/", TemplateHandler(tm, "home/index", nil))
	r.Get("/contact", TemplateHandler(tm, "home/contact", nil))
	r.Get("/faq", c.FAQ(tm))

	return r
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
