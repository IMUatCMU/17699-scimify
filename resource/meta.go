// Meta section for the SCIM resource
package resource

type Meta struct {
	ResourceType string    `json:"resourceType"`
	Created      string 	`json:"created"`
	LastModified string 	`json:"lastModified"`
	Location     string    `json:"location"`
	Version      string    `json:"version"`
}
