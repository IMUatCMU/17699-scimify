package persistence

import "github.com/go-scim/scimify/resource"

// An implementation of Repository interface that stores SCIM resource
// in a hash map. It only offers functionality to create, delete and get SCIM
// resource by id.
type SimpleRepository struct {
	repo map[string]*resource.Resource
}

// Put the provided resource into the hash map, overwrites existing id
func (r *SimpleRepository) Create(resource *resource.Resource, context resource.Context) error {
	r.repo[resource.Id] = resource
	return nil
}

// Get the resource indexed by id in the hash map
func (r *SimpleRepository) Get(id string, context resource.Context) (*resource.Resource, error) {
	return r.repo[id], nil
}

// Not implemented by design
func (r *SimpleRepository) Replace(id string, resource *resource.Resource, context resource.Context) error {
	return nil
}

// Overwrite the indexed id slot with nil
func (r *SimpleRepository) Delete(id string, context resource.Context) error {
	r.repo[id] = nil
	return nil
}

// Not implemented by design
func (r *SimpleRepository) Query(filter string, sortBy string, ascending bool, pageStart int32, pageSize int32, context resource.Context) ([]*resource.Resource, error) {
	return nil, nil
}
