package persistence

import (
	"context"
	"github.com/go-scim/scimify/resource"
)

type Repository interface {

	// Create a new SCIM resource in the repository
	Create(resource *resource.Resource, context context.Context) error

	// Get a SCIM resource by id from the repository
	Get(id string, context context.Context) (*resource.Resource, error)

	// Completely replace a SCIM resource by id in the repository
	Replace(id string, resource *resource.Resource, context context.Context) error

	// Delete a SCIM resource from the repository
	Delete(id string, context context.Context) error

	// Query a list of SCIM resource by the provided filter, sort and pagination parameters
	Query(filter interface{}, sortBy string, ascending bool, pageStart int, pageSize int, context context.Context) ([]*resource.Resource, error)
}
