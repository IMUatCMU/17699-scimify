package serialize

import (
	"encoding/json"
	"github.com/go-scim/scimify/resource"
)

// A simple/default implementation that uses Go's json marshal capability
type DefaultJsonSerializer struct{}

func (_ *DefaultJsonSerializer) Serialize(resource *resource.Resource, context interface{}) ([]byte, error) {
	return json.Marshal(resource)
}
