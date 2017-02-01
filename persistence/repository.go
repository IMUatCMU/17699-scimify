package persistence

import "github.com/go-scim/scimify/resource"

type Repository interface {

	// Create a new SCIM resource in the repository
	Create(resource *resource.Resource) error

	// Get a SCIM resource by id from the repository
	Get(id string) (*resource.Resource, error)

	// Completely replace a SCIM resource by id in the repository
	Replace(id string, resource *resource.Resource) error

	// Delete a SCIM resource from the repository
	Delete(id string) error

	// Query a list of SCIM resource by the provided filter, sort and pagination parameters
	Query(filter string, sortBy string, ascending bool, pageStart int32, pageSize int32) ([]*resource.Resource, error)
}
