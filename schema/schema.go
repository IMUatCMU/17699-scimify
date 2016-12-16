// Model for the SCIM schema
package schema

type Schema struct {
	Schemas     []string     `json:"schemas"`
	Id          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Attributes  []*Attribute `json:"attributes"`
}
