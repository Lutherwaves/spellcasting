package types

import "time"

// Resource is the database and API representation of a RESOURCE.
type Resource struct {
	ID        string    `json:"id" db:"id"`
	TenantID  string    `json:"tenant_id" db:"tenant_id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateResource is the request body for POST /RESOURCES.
type CreateResource struct {
	Name string `json:"name"`
}

// UpdateResource is the request body for PATCH /RESOURCES/{id}.
// All fields are optional; only non-nil values are applied.
type UpdateResource struct {
	Name *string `json:"name,omitempty"`
}
