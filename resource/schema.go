// Model for the SCIM schema and express methods
package resource

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

type Schema struct {
	Schemas     []string              `json:"schemas"`
	Id          string                `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Attributes  []*Attribute          `json:"attributes"`
	attrIndex   map[string]*Attribute `json:"-"`
}

func (s *Schema) GetId() string {
	return s.Id
}

func (s *Schema) Data() map[string]interface{} {
	data := map[string]interface{}{
		"schemas":     s.Schemas,
		"id":          s.Id,
		"name":        s.Name,
		"description": s.Description,
		"attributes":  make([]map[string]interface{}, 0, len(s.Attributes)),
	}
	for _, attr := range s.Attributes {
		data["attributes"] = append(data["attributes"].([]map[string]interface{}), attr.ToMap())
	}
	return data

}

func (s *Schema) AsAttribute() *Attribute {
	return &Attribute{
		Type:          Complex,
		MultiValued:   false,
		Required:      false,
		CaseExact:     false,
		Returned:      Default,
		Uniqueness:    None,
		Mutability:    ReadWrite,
		SubAttributes: s.Attributes,
		Assist:        &Assist{JSONName: "", Path: "", FullPath: ""},
	}
}

func (s *Schema) MergeWith(schemas ...*Schema) {
	for _, schema := range schemas {
		s.Attributes = append(s.Attributes, schema.Attributes...)
	}
}

// Construct an index system from attribute path and full path to attribute itself
// Indexes are all lower cased to allow case insensitive search (attribute names in user queries can be case insensitive)
func (s *Schema) ConstructAttributeIndex() {
	collectIndex := func(attr *Attribute) {
		if attr.Assist != nil {
			if len(attr.Assist.FullPath) > 0 {
				s.attrIndex[strings.ToLower(attr.Assist.FullPath)] = attr
			}
		}
	}

	s.attrIndex = make(map[string]*Attribute, 0)
	for _, attr := range s.Attributes {
		collectIndex(attr)
		for _, subAttr := range attr.SubAttributes {
			collectIndex(subAttr)
		}
	}
}

// Get the attribute from the index constructed, prerequisite is calling ConstructAttributeIndex() method first.
// If $path starts with a valid resource URN, convert to lower case and get from index
// If $path does not start with a valid resource URN, append id (which is a resource URN) of this schema and get from index
// If $path are keys from the core schema (i.e. id, externalId, meta, and so on) convert to lower case and get from index
func (s *Schema) GetAttribute(path string) *Attribute {
	switch strings.ToLower(path) {
	case "schemas",
		"id", "externalid",
		"meta", "meta.resourcetype", "meta.created", "meta.lastmodified", "meta.location", "meta.version":
		return s.attrIndex[strings.ToLower(path)]
	default:
		for _, resourceUrn := range AllResourceUrns {
			if strings.HasPrefix(strings.ToLower(path), strings.ToLower(resourceUrn+":")) {
				return s.attrIndex[strings.ToLower(path)]
			}
		}
		return s.attrIndex[strings.ToLower(s.Id)+":"+strings.ToLower(path)]
	}
}

// Load schema from a designated file path
func LoadSchema(path string) (*Schema, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	schema := &Schema{}
	err = json.Unmarshal(fileBytes, schema)
	if err != nil {
		return nil, err
	}

	return schema, nil
}
