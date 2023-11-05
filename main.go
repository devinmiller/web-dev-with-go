package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/devinmiller/web-dev-with-go/controllers"
	"github.com/devinmiller/web-dev-with-go/services"
	"github.com/devinmiller/web-dev-with-go/templates"
	"github.com/devinmiller/web-dev-with-go/views"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initDatabase(uri string) (client *mongo.Client, err error) {
	fmt.Println("Connecting to database...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return
	}

	return
}

func termDatabase(client *mongo.Client) {
	fmt.Println("Disconnecting from database...")

	if err := client.Disconnect(context.Background()); err != nil {
		panic(err)
	}
}

type App struct {
	client *mongo.Client
}

func main() {
	fmt.Println("Starting server on :3000...")

	client, err := initDatabase(os.Getenv("MONGODB_URI"))
	if err != nil {
		panic(err)
	}

	defer termDatabase(client)

	tm, err := views.NewTemplateManager(templates.FS, ".", "layouts", ".html")
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	csrfMiddleware := csrf.Protect(
		// TODO: Store session key literally anywhere else
		[]byte("3F5G6H78J9KLMN0P1QR2ST3UVW4XYZ5A"),
		// TODO: Make this jank website secure
		csrf.Secure(false),
	)

	r.Use(csrfMiddleware)
	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		})
	})

	userService := services.NewUserService(client)
	userController := controllers.NewHomeController(tm, userService)

	r.Mount("/", userController.Routes())

	// Set up static file server for assets
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("dist"))))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	err = http.ListenAndServe(":3000", r)

	if err != nil {
		panic(err)
	}
}
