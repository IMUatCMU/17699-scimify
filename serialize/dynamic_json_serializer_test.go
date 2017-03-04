package serialize

import (
	"context"
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/resource"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func BenchmarkDynamicJsonSerializer_Serialize(b *testing.B) {
	// prepare schema
	schema, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		b.Fatal(err)
	}
	schema.ConstructAttributeIndex()

	// prepare data
	target, _, err := helper.LoadResource("../test_data/single_test_user_david.json")
	if err != nil {
		b.Fatal(err)
	}

	serializer := &dynamicJsonSerializer{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx := context.Background()
			ctx = context.WithValue(ctx, resource.CK_Schema, schema)
			serializer.Serialize(target, nil, nil, ctx)
		}
	})
}

func TestDynamicJsonSerializer_Serialize(t *testing.T) {
	// prepare schema
	schema, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		t.Fatal(err)
	}
	schema.ConstructAttributeIndex()

	// prepare data
	target, json, err := helper.LoadResource("../test_data/single_test_user_david.json")
	if err != nil {
		t.Fatal(err)
	}

	// remove password from comparison since it won't be in it.
	json = strings.Replace(json, "\"password\": \"t1meMa$heen\",", "", 1)

	serializer := &dynamicJsonSerializer{}

	ctx := context.Background()
	ctx = context.WithValue(ctx, resource.CK_Schema, schema)
	bytes, err := serializer.Serialize(target, nil, nil, ctx)
	assert.Nil(t, err)
	assert.JSONEq(t, json, string(bytes))
}
