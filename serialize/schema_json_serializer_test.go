package serialize

import (
	"encoding/json"
	"github.com/go-scim/scimify/resource"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestSchemaJsonSerializer_Serialize(t *testing.T) {
	// prepare schema
	schema := &resource.Schema{
		Schemas:    []string{resource.SchemaUrn},
		Id:         resource.UserUrn,
		Name:       "User schema",
		Attributes: make([]*resource.Attribute, 0),
	}
	coreSchema, err := loadSchema("../schemas/common_schema.json")
	userSchema, err := loadSchema("../schemas/user_schema.json")
	if err != nil {
		t.Fatal(err)
	}
	schema.MergeWith(coreSchema, userSchema)
	schema.ConstructAttributeIndex()

	// prepare data
	resource, origJson, err := loadResource("../test_data/single_test_user_david.json")
	if err != nil {
		t.Fatal(err)
	}

	// serializer
	serializer := &SchemaJsonSerializer{}
	json, err := serializer.Serialize(resource, &SchemaJsonSerializerContext{
		InclusionPaths: make([]string, 0),
		ExclusionPaths: make([]string, 0),
		Schema:         schema,
	})
	assert.Nil(t, err)
	assert.JSONEq(t, origJson, string(json))
}

func TestSchemaJsonSerializer_Serialize_Error(t *testing.T) {
	// prepare schema
	schema := &resource.Schema{
		Schemas:    []string{resource.SchemaUrn},
		Id:         resource.UserUrn,
		Name:       "User schema",
		Attributes: make([]*resource.Attribute, 0),
	}
	coreSchema, err := loadSchema("../schemas/common_schema.json")
	userSchema, err := loadSchema("../schemas/user_schema.json")
	if err != nil {
		t.Fatal(err)
	}
	schema.MergeWith(coreSchema, userSchema)
	schema.ConstructAttributeIndex()

	// prepare data
	resource, _, err := loadResource("../test_data/single_test_user_david.json")
	if err != nil {
		t.Fatal(err)
	}

	// slightly alter resource in a wrong way
	resource.Attributes["addresses"] = make(map[string]interface{}, 0)

	// serializer
	serializer := &SchemaJsonSerializer{}
	result, err := serializer.Serialize(resource, &SchemaJsonSerializerContext{
		InclusionPaths: make([]string, 0),
		ExclusionPaths: make([]string, 0),
		Schema:         schema,
	})
	assert.NotNil(t, err)
	t.Log(string(result))
}

func loadSchema(filePath string) (*resource.Schema, error) {
	if path, err := filepath.Abs(filePath); err != nil {
		return nil, err
	} else if schema, err := resource.LoadSchema(path); err != nil {
		return nil, err
	} else {
		return schema, nil
	}
}

func loadResource(filePath string) (*resource.Resource, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, "", err
	}

	data := make(map[string]interface{}, 0)
	err = json.Unmarshal(fileBytes, &data)
	if err != nil {
		return nil, "", err
	}

	return resource.NewResourceFromMap(data), string(fileBytes), nil
}
