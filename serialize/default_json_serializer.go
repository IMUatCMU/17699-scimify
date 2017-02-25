package serialize

import (
	"context"
	"encoding/json"
	"github.com/go-scim/scimify/resource"
)

// A simple/default implementation that uses Go's json marshal capability
type DefaultJsonSerializer struct{}

func (_ *DefaultJsonSerializer) Serialize(resource *resource.Resource, inclusionPaths, exclusionPaths []string, context context.Context) ([]byte, error) {
	return json.Marshal(resource)
}
