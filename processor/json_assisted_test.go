package processor

import (
	"fmt"
	"github.com/go-scim/scimify/helper"
	"github.com/spf13/viper"
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
		SerializationTargetFunc: func() interface{} {
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
				SerializationTargetFunc: func() interface{} {
					return resource
				},
			})
			if nil != err {
				b.Fatal(err)
			}
		}
	})
}

func benchmarkAssistedJsonSerializationProcessor(poolSize int, b *testing.B) {
	viper.Set(fmt.Sprintf("scim.threadPool.%s", JsonAssisted), poolSize)

	schema, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		b.Fatal(err)
	}
	schema.ConstructAttributeIndex()

	r, _, err := helper.LoadResource("../test_data/single_test_user_david.json")
	if err != nil {
		b.Fatal(err)
	}

	makeContext := func() *ProcessorContext {
		return &ProcessorContext{
			Schema: schema,
			SerializationTargetFunc: func() interface{} {
				return r
			},
		}
	}

	// serializer
	b.ResetTimer()
	processor := GetWorkerBean(JsonAssisted)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			context := makeContext()
			err := processor.Process(context)
			if nil != err {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkAssistedJsonSerializationProcessorWithPoolSize1(b *testing.B) {
	benchmarkAssistedJsonSerializationProcessor(1, b)
}

func BenchmarkAssistedJsonSerializationProcessorWithPoolSize2(b *testing.B) {
	benchmarkAssistedJsonSerializationProcessor(2, b)
}

func BenchmarkAssistedJsonSerializationProcessorWithPoolSize3(b *testing.B) {
	benchmarkAssistedJsonSerializationProcessor(3, b)
}

func BenchmarkAssistedJsonSerializationProcessorWithPoolSize4(b *testing.B) {
	benchmarkAssistedJsonSerializationProcessor(4, b)
}

func BenchmarkAssistedJsonSerializationProcessorWithPoolSize5(b *testing.B) {
	benchmarkAssistedJsonSerializationProcessor(5, b)
}

func BenchmarkAssistedJsonSerializationProcessorWithPoolSize6(b *testing.B) {
	benchmarkAssistedJsonSerializationProcessor(6, b)
}

func BenchmarkAssistedJsonSerializationProcessorWithPoolSize7(b *testing.B) {
	benchmarkAssistedJsonSerializationProcessor(7, b)
}

func BenchmarkAssistedJsonSerializationProcessorWithPoolSize8(b *testing.B) {
	benchmarkAssistedJsonSerializationProcessor(8, b)
}

func BenchmarkAssistedJsonSerializationProcessorWithPoolSize9(b *testing.B) {
	benchmarkAssistedJsonSerializationProcessor(9, b)
}
