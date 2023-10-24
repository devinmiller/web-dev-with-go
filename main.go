package main

import (
	"fmt"
	"net/http"

	"github.com/devinmiller/web-dev-with-go/controllers"
	"github.com/devinmiller/web-dev-with-go/templates"
	"github.com/devinmiller/web-dev-with-go/views"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	fmt.Println("Starting server on :3000...")

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	tm, err := views.NewTemplateManager(templates.FS, ".", "layouts", ".html")

	// error loading templates
	if err != nil {
		panic(err)
	}

	r.Mount("/", controllers.HomeController{}.Routes(tm))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	err = http.ListenAndServe(":3000", r)

	if err != nil {
		panic(err)
	}
}
