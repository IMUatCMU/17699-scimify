package validation

import (
	"github.com/go-scim/scimify/resource"
	"sync"
)

type Validator interface {
	Validate(*resource.Resource, *ValidatorContext) (bool, error)
}

type ValidatorContext struct {
	sync.RWMutex
	Data map[string]interface{}
}

// Reserved keys in ValidatorContext
const (
	ReferenceResource = "ReferenceResource"
	Schema            = "Schema"
)

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
	typeCheck = "type_check"
)
