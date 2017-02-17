package worker

import (
	"encoding/json"
	"github.com/go-scim/scimify/resource"
	"github.com/go-scim/scimify/serialize"
	"io/ioutil"
	"os"
	"testing"
)

func BenchmarkSchemaAssistedJsonSerializerWorker(b *testing.B) {
	worker := GetSchemaAssistedJsonSerializerWorker()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			b.StopTimer()
			resource, schema := prepareResourceAndSchema(b)
			b.StartTimer()

			worker.Do(&JsonSerializeInput{
				Resource: resource,
				Context: &serialize.SchemaJsonSerializerContext{
					InclusionPaths: []string{},
					ExclusionPaths: []string{},
					Schema:         schema,
				},
			})
		}
	})
}

func BenchmarkDefaultJsonSerializerWorker(b *testing.B) {
	resource, _ := prepareResourceAndSchema(b)
	worker := GetDefaultJsonSerializerWorker()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			worker.Do(&JsonSerializeInput{Resource: resource, Context: nil})
		}
	})
}

func prepareResourceAndSchema(b *testing.B) (*resource.Resource, *resource.Schema) {
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
		b.Fatal(err)
	}
	schema.MergeWith(coreSchema, userSchema)
	schema.ConstructAttributeIndex()

	// prepare data
	file, err := os.Open("../test_data/single_test_user_david.json")
	if err != nil {
		b.Fatal(err)
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		b.Fatal(err)
	}

	data := make(map[string]interface{}, 0)
	err = json.Unmarshal(fileBytes, &data)
	if err != nil {
		b.Fatal(err)
	}

	return resource.NewResourceFromMap(data), schema
}
