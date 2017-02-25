package validation

import (
	"context"
	"github.com/go-scim/scimify/resource"
)

type Validator interface {
	Validate(*resource.Resource, ValidationOptions, context.Context) (bool, error)
}

type ValidationOptions struct {
	ReadOnlyIsMandatory          bool
	UnassignedImmutableIsIgnored bool
}

// Validation Error
type validationError struct {
	ViolationType string
	Message       string
	FullPath      string
}

func (e *validationError) Error() string {
	return e.Message
}

// Constant for violation type
const (
	typeCheck       = "type_check"
	requiredMissing = "required_missing"
	mutabilityCheck = "mutability_check"
	unknown         = "unknown"
)
