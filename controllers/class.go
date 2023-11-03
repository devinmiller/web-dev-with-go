package controllers

import (
	"github.com/devinmiller/web-dev-with-go/services"
	"github.com/devinmiller/web-dev-with-go/views"
	"github.com/go-chi/chi/v5"
)

type ClassController struct {
	views        *views.TemplateManager
	classService services.ClassService
}

func NewClassController(
	views *views.TemplateManager,
	classService services.UserService) ClassController {

	c := ClassController{
		views:        views,
		classService: classService,
	}

	return c
}

func (c ClassController) Routes() chi.Router {
	r := chi.NewRouter()

	return r
}
