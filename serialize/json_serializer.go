package serialize

import (
	"context"
	"github.com/go-scim/scimify/resource"
)

type JSONSerializer interface {

	// Serialize the given resource, with or without the help of the context to JSON bytes
	Serialize(resource *resource.Resource, inclusionPaths, exclusionPaths []string, context context.Context) ([]byte, error)
}
