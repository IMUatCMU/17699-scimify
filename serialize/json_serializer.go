package serialize

import (
	"context"
	"github.com/go-scim/scimify/resource"
)

type JSONSerializer interface {

	// Serialize the given resource, with or without the help of the context to JSON bytes
	Serialize(target resource.ScimObject, inclusionPaths, exclusionPaths []string, context context.Context) ([]byte, error)

	SerializeArray(target []resource.ScimObject, inclusionPaths, exclusionPaths []string, context context.Context) ([]byte, error)
}
