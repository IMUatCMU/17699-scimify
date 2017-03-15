package persistence

import (
	"github.com/go-scim/scimify/resource"
)

// An implementation of Repository interface that stores SCIM resource
// in a hash map. It only offers functionality to create, delete and get SCIM
// resource by id.
type SimpleRepository struct {
	repo map[string]resource.ScimObject
}

// Put the provided resource into the hash map, overwrites existing id
func (r *SimpleRepository) Create(resource resource.ScimObject) error {
	r.repo[resource.GetId()] = resource
	return nil
}

func (r *SimpleRepository) GetAll() ([]resource.ScimObject, error) {
	resources := make([]resource.ScimObject, 0, len(r.repo))
	for _, v := range r.repo {
		resources = append(resources, v)
	}
	return resources, nil
}

// Get the resource indexed by id in the hash map
func (r *SimpleRepository) Get(id string) (resource.ScimObject, error) {
	return r.repo[id], nil
}

// Not implemented by design
func (r *SimpleRepository) Replace(id string, resource resource.ScimObject) error {
	return nil
}

// Overwrite the indexed id slot with nil
func (r *SimpleRepository) Delete(id string) error {
	r.repo[id] = nil
	return nil
}

// Not implemented by design
func (r *SimpleRepository) Query(filter interface{}, sortBy string, ascending bool, pageStart int, pageSize int) ([]resource.ScimObject, error) {
	return nil, nil
}
