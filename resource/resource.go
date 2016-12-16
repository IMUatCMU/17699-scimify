// Model and express methods for a SCIM resource
// schemas and common attributes are represented explicitly
// core attributes and extension attributes are modelled as map
package resource

type Resource struct {
	Schemas    []string `json:"schemas"`
	Id         string   `json:"id"`
	ExternalId string   `json:"externalId"`
	Meta       *Meta    `json:"meta"`
	Attributes map[string]interface{}
}

// Create a new resource
func NewResource(resourceBaseUrn string) *Resource {
	return &Resource{
		Schemas:    []string{resourceBaseUrn},
		Meta:       &Meta{},
		Attributes: make(map[string]interface{}, 0),
	}
}
