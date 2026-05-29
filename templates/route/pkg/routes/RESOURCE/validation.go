package RESOURCE

import (
	"fmt"

	"SERVICENAME/pkg/types"

	"github.com/tink3rlabs/magic/errors"
)

const maxNameLen = 1000

// ValidateCreate checks required fields and length bounds on a CreateResource payload.
// Returns *errors.BadRequest when validation fails, nil on success.
func ValidateCreate(in types.CreateResource) error {
	if in.Name == "" {
		return &errors.BadRequest{Message: "name is required"}
	}
	if len(in.Name) > maxNameLen {
		return &errors.BadRequest{Message: fmt.Sprintf("name must be at most %d characters", maxNameLen)}
	}
	return nil
}

// ValidateUpdate checks length bounds on an UpdateResource payload.
// All fields are optional on update; only non-nil fields are validated.
func ValidateUpdate(in types.UpdateResource) error {
	if in.Name != nil {
		if *in.Name == "" {
			return &errors.BadRequest{Message: "name must not be empty"}
		}
		if len(*in.Name) > maxNameLen {
			return &errors.BadRequest{Message: fmt.Sprintf("name must be at most %d characters", maxNameLen)}
		}
	}
	return nil
}
