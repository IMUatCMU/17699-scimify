// Model for the SCIM schema and express methods
package schema

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Schema struct {
	Schemas     []string     `json:"schemas"`
	Id          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Attributes  []*Attribute `json:"attributes"`
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
