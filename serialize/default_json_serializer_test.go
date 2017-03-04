package serialize

import (
	"github.com/go-scim/scimify/helper"
	"testing"
)

func BenchmarkDefaultJsonSerializer_Serialize(b *testing.B) {
	// prepare schema
	schema, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	schema.ConstructAttributeIndex()

	// prepare data
	resource, _, err := helper.LoadResource("../test_data/single_test_user_david.json")
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
