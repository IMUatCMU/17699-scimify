package persistence

import (
	"github.com/go-scim/scimify/resource"
)

type Repository interface {

	// Create a new SCIM resource in the repository
	Create(resource.ScimObject) error

	// Get all SCIM resource in this repository (optional)
	GetAll() ([]resource.ScimObject, error)

	// Get a SCIM resource by id from the repository
	Get(string) (resource.ScimObject, error)

	// Completely replace a SCIM resource by id in the repository
	Replace(string, resource.ScimObject) error

	// Delete a SCIM resource from the repository
	Delete(string) error

	// Query a list of SCIM resource by the provided filter, sort and pagination parameters
	Query(filter interface{}, sortBy string, ascending bool, pageStart int, pageSize int) ([]resource.ScimObject, error)
}
