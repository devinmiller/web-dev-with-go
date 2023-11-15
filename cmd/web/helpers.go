package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/devinmiller/web-dev-with-go/context"
	"github.com/devinmiller/web-dev-with-go/models"
	"github.com/gorilla/csrf"
)

func (app *application) View(name string) http.HandlerFunc {
	return app.ViewHandler(name, func(r *http.Request) {})
}

func (app *application) ViewHandler(name string, f func(r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := app.views.Template(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		f(r)

		// include csrf token and user template funcs in every request
		tmpl = tmpl.Funcs(template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
			"currentUser": func() *models.User {
				return context.User(r.Context())
			},
		})

		buf, err := app.views.RenderTemplate(tmpl, nil)
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

func (app *application) CurrentUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := app.sessionService.CurrentUser(w, r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (app *application) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		fmt.Println(user)
		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
