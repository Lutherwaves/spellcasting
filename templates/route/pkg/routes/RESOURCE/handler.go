package RESOURCE

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"SERVICENAME/pkg/types"

	"github.com/tink3rlabs/magic/errors"
)

// List handles GET / — cursor-paginated list with optional Lucene filter.
func (router *Router) List(w http.ResponseWriter, r *http.Request) error {
	cursor := r.URL.Query().Get("next")
	filter := r.URL.Query().Get("filter")

	limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 1000 {
		limit = 1000
	}

	items, next, err := router.service.List(r.Context(), filter, int(limit), cursor)
	if err != nil {
		return &errors.BadRequest{Message: err.Error()}
	}

	render.JSON(w, r, map[string]any{"items": items, "next": next})
	return nil
}

// Get handles GET /{id}.
func (router *Router) Get(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	item, err := router.service.Get(r.Context(), id)
	if err != nil {
		return err
	}

	render.JSON(w, r, item)
	return nil
}

// Create handles POST /.
func (router *Router) Create(w http.ResponseWriter, r *http.Request) error {
	var in types.CreateResource
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return &errors.BadRequest{Message: err.Error()}
	}

	if err := ValidateCreate(in); err != nil {
		return err
	}

	item, err := router.service.Create(r.Context(), in)
	if err != nil {
		return err
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, item)
	return nil
}

// Update handles PATCH /{id}.
func (router *Router) Update(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	var in types.UpdateResource
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return &errors.BadRequest{Message: err.Error()}
	}

	if err := ValidateUpdate(in); err != nil {
		return err
	}

	item, err := router.service.Update(r.Context(), id, in)
	if err != nil {
		return err
	}

	render.JSON(w, r, item)
	return nil
}

// Delete handles DELETE /{id}.
func (router *Router) Delete(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	if err := router.service.Delete(r.Context(), id); err != nil {
		return err
	}

	render.NoContent(w, r)
	return nil
}
