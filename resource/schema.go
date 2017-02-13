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
// The rule in general is core attributes (i.e. stock User and Group attributes) are indexed with their short names,
// which means the assist.fullName has their short names; Extension attributes must be defined with their full names
// to avoid collision
func (s *Schema) GetAttribute(path string) *Attribute {
	if strings.HasPrefix(path, s.Id+":") {
		return s.attrIndex[strings.ToLower(path[len(s.Id+":")+1:])]
	} else {
		return s.attrIndex[strings.ToLower(path)]
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
