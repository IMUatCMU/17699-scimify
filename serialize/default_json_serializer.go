package serialize

import (
	"context"
	"encoding/json"
	"github.com/go-scim/scimify/resource"
)

// A simple/default implementation that uses Go's json marshal capability
type DefaultJsonSerializer struct{}

func (_ *DefaultJsonSerializer) Serialize(target resource.ScimObject, inclusionPaths, exclusionPaths []string, context context.Context) ([]byte, error) {
	return json.Marshal(target)
}

func (_ *DefaultJsonSerializer) SerializeArray(target []resource.ScimObject, inclusionPaths, exclusionPaths []string, context context.Context) ([]byte, error) {
	return json.Marshal(target)
}
