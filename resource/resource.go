// Model and express methods for a SCIM resource
// schemas and common attributes are represented explicitly
// core attributes and extension attributes are modelled as map
package resource

import (
	"sync"
)

type Resource struct {
	sync.RWMutex
	Attributes map[string]interface{}
}

func (r *Resource) ToMap() map[string]interface{} {
	data := make(map[string]interface{}, len(r.Attributes))
	r.RLock()
	for k, v := range r.Attributes {
		data[k] = v
	}
	r.RUnlock()
	return data
}

// Create a new resource
func NewResource() *Resource {
	return &Resource{
		Attributes: make(map[string]interface{}, 0),
	}
}

// Create a new resource from map data
func NewResourceFromMap(data map[string]interface{}) *Resource {
	resource := NewResource()
	resource.Attributes = data
	return resource
}
