package RESOURCE

import (
	"github.com/go-chi/chi/v5"

	"github.com/tink3rlabs/magic/middlewares"
)

// Router wires the 5 standard CRUD endpoints for Resources.
type Router struct {
	Router  chi.Router
	service ServiceInterface
}

// NewResourceRouter constructs a chi router mounted with all RESOURCE endpoints.
func NewResourceRouter(service ServiceInterface) *Router {
	r := chi.NewRouter()
	router := &Router{Router: r, service: service}

	h := middlewares.ErrorHandler{}

	r.Get("/", h.Wrap(router.List))
	r.Post("/", h.Wrap(router.Create))
	r.Get("/{id}", h.Wrap(router.Get))
	r.Patch("/{id}", h.Wrap(router.Update))
	r.Delete("/{id}", h.Wrap(router.Delete))

	return router
}
