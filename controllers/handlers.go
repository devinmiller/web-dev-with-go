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

func RenderHandler(tm *views.TemplateManager, name string, data interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := tm.Render(w, name, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func FAQ(tm *views.TemplateManager) http.HandlerFunc {
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

	return RenderHandler(tm, "faq", questions)
}
