// Model and express methods for a SCIM resource
// schemas and common attributes are represented explicitly
// core attributes and extension attributes are modelled as map
package resource

import (
	"encoding/json"
	"sync"
)

type Resource struct {
	sync.RWMutex
	Attributes map[string]interface{}
}

func (r *Resource) GetId() string {
	if id, ok := r.Attributes["id"].(string); ok {
		return id
	} else {
		return ""
	}
}

func (r *Resource) Data() map[string]interface{} {
	return r.ToMap()
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

// Create a new resource from bytes of json data
func NewResourceFromBytes(bytes []byte) (*Resource, error) {
	data := make(map[string]interface{}, 0)
	err := json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return NewResourceFromMap(data), nil
}
