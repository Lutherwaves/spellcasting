package RESOURCE

import (
	"context"
	"time"

	"github.com/google/uuid"

	"SERVICENAME/pkg/types"

	"github.com/tink3rlabs/magic/errors"
	"github.com/tink3rlabs/magic/storage"
)

// ServiceInterface is the contract the route layer depends on.
// Swap the concrete Service for a mock in tests.
type ServiceInterface interface {
	List(ctx context.Context, filter string, limit int, cursor string) ([]types.Resource, string, error)
	Get(ctx context.Context, id string) (*types.Resource, error)
	Create(ctx context.Context, in types.CreateResource) (*types.Resource, error)
	Update(ctx context.Context, id string, in types.UpdateResource) (*types.Resource, error)
	Delete(ctx context.Context, id string) error
}

// Service is the concrete implementation backed by a magic StorageAdapter.
type Service struct {
	storage storage.StorageAdapter
}

// NewService constructs a Service with the given storage adapter.
// The adapter is typically constructed in cmd/server.go and shared across services.
func NewService(s storage.StorageAdapter) *Service {
	return &Service{storage: s}
}

// List returns a cursor-paginated slice of Resources.
// When filter is non-empty it is treated as a Lucene query string and passed to
// storage.Search; otherwise storage.List is used (no filter predicate).
func (s *Service) List(_ context.Context, filter string, limit int, cursor string) ([]types.Resource, string, error) {
	items := []types.Resource{}

	var next string
	var err error

	if filter != "" {
		// storage.Search accepts a raw Lucene string; the adapter parses and evaluates it.
		next, err = s.storage.Search(&items, "id", filter, limit, cursor)
	} else {
		next, err = s.storage.List(&items, "id", map[string]any{}, limit, cursor)
	}

	return items, next, err
}

// Get returns a single Resource by ID, or *errors.NotFound if it doesn't exist.
func (s *Service) Get(_ context.Context, id string) (*types.Resource, error) {
	item := &types.Resource{}
	if err := s.storage.Get(item, map[string]any{"id": id}); err != nil {
		return nil, &errors.NotFound{Message: "Resource not found"}
	}
	return item, nil
}

// Create inserts a new Resource and returns it with a generated ID and timestamps.
// Uses UUIDv7 for time-ordered IDs that support cursor pagination without extra columns.
func (s *Service) Create(_ context.Context, in types.CreateResource) (*types.Resource, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	item := types.Resource{
		ID:        id.String(),
		Name:      in.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.storage.Create(item); err != nil {
		return nil, err
	}

	return &item, nil
}

// Update applies partial fields from UpdateResource to the existing record.
// Returns *errors.NotFound if the record doesn't exist.
func (s *Service) Update(ctx context.Context, id string, in types.UpdateResource) (*types.Resource, error) {
	item, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if in.Name != nil {
		item.Name = *in.Name
	}
	item.UpdatedAt = time.Now().UTC()

	if err := s.storage.Update(*item, map[string]any{"id": id}); err != nil {
		return nil, err
	}

	return item, nil
}

// Delete removes a Resource by ID.
func (s *Service) Delete(_ context.Context, id string) error {
	return s.storage.Delete(&types.Resource{}, map[string]any{"id": id})
}
