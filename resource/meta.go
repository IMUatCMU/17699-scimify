// Meta section for the SCIM resource
package resource

type Meta struct {
	ResourceType string `json:"resourceType"`
	Created      string `json:"created"`
	LastModified string `json:"lastModified"`
	Location     string `json:"location"`
	Version      string `json:"version"`
}

func (m *Meta) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"resourceType": m.ResourceType,
		"created":      m.Created,
		"lastModified": m.LastModified,
		"location":     m.Location,
		"version":      m.Version,
	}
}
