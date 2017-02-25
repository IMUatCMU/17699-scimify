package serialize

import (
	"github.com/go-scim/scimify/resource"
	"testing"
)

func BenchmarkDefaultJsonSerializer_Serialize(b *testing.B) {
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
	resource, _, err := loadResource("../test_data/single_test_user_david.json")
	if err != nil {
		b.Fatal(err)
	}

	// serializer
	b.ResetTimer()
	serializer := &DefaultJsonSerializer{}
	for i := 0; i < b.N; i++ {
		serializer.Serialize(resource, nil, nil, nil)
	}
}
