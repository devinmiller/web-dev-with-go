package controllers

import (
	"net/http"

	"github.com/devinmiller/web-dev-with-go/services"
	"github.com/devinmiller/web-dev-with-go/views"
)

type Controller struct {
	views          *views.TemplateManager
	userService    *services.UserService
	sessionService *services.SessionService
}

func (controller Controller) ViewHandler(name string, f func(r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := controller.views.Template(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		f(r)

		buf, err := controller.views.RenderTemplate(tmpl, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(buf)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
