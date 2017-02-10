package persistence

import "github.com/go-scim/scimify/resource"

type Repository interface {

	// Create a new SCIM resource in the repository
	Create(resource *resource.Resource, context resource.Context) error

	// Get a SCIM resource by id from the repository
	Get(id string, context resource.Context) (*resource.Resource, error)

	// Completely replace a SCIM resource by id in the repository
	Replace(id string, resource *resource.Resource, context resource.Context) error

	// Delete a SCIM resource from the repository
	Delete(id string, context resource.Context) error

	// Query a list of SCIM resource by the provided filter, sort and pagination parameters
	Query(filter string, sortBy string, ascending bool, pageStart int32, pageSize int32, context resource.Context) ([]*resource.Resource, error)
}
