// Model for the SCIM schema and express methods
package resource

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Schema struct {
	Schemas     []string              `json:"schemas"`
	Id          string                `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Attributes  []*Attribute          `json:"attributes"`
	attrIndex   map[string]*Attribute `json:"-"`
}

// Construct an index system from attribute path and full path to attribute itself
func (s *Schema) ConstructAttributeIndex() {
	collectIndex := func(attr *Attribute) {
		if attr.Assist != nil {
			if len(attr.Assist.FullPath) > 0 {
				s.attrIndex[attr.Assist.FullPath] = attr
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
func (s *Schema) GetAttribute(path string) *Attribute {
	return s.attrIndex[path]
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
