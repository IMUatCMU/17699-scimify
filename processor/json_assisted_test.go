package processor

import (
	"github.com/go-scim/scimify/helper"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestAssistedJsonSerializationProcessor_Process(t *testing.T) {
	processor := &assistedJsonSerializationProcessor{}

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

	ctx := &ProcessorContext{
		Schema: schema,
		SerializationTargetFunc: func()interface{} {
			return target
		},
	}

	err = processor.Process(ctx)
	assert.Nil(t, err)
	assert.JSONEq(t, json, string(ctx.ResponseBody))
}

func BenchmarkAssistedJsonSerializationProcessor_Process(b *testing.B) {
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
	processor := &assistedJsonSerializationProcessor{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := processor.Process(&ProcessorContext{
				Schema: schema,
				SerializationTargetFunc: func()interface{} {
					return resource
				},
			})
			if nil != err {
				b.Fatal(err)
			}
		}
	})
}
