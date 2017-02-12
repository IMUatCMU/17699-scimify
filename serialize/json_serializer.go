package serialize

import "github.com/go-scim/scimify/resource"

type JSONSerializer interface {

	// Serialize the given resource, with or without the help of the context to JSON bytes
	Serialize(resource *resource.Resource, context interface{}) ([]byte, error)
}
