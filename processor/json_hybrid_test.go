package processor

import (
	"encoding/json"
	"fmt"
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/resource"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestHybridJsonSerializationProcessor_Process(t *testing.T) {
	processor := &hybridJsonSerializationProcessor{
		sjp: &simpleJsonSerializationProcessor{},
		ajp: &assistedJsonSerializationProcessor{},
		f: func(bytes []byte, _ *ProcessorContext) interface{} {
			raw := json.RawMessage(bytes)
			return resource.NewListResponse(&raw, 0, 10, 2)
		},
	}

	schema, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		t.Fatal(err)
	}
	schema.ConstructAttributeIndex()

	david, json, err := helper.LoadResource("../test_data/single_test_user_david.json")
	if err != nil {
		t.Fatal(err)
	}
	json = strings.Replace(json, "\"password\": \"t1meMa$heen\",", "", 1)

	ctx := &ProcessorContext{
		Schema: schema,
		SerializationTargetFunc: func() interface{} {
			return []resource.ScimObject{david, david}
		},
	}

	err = processor.Process(ctx)
	assert.Nil(t, err)
	assert.JSONEq(t, fmt.Sprintf(`{
		"schemas":["urn:ietf:params:scim:api:messages:2.0:ListResponse"],
		"totalResults":2,
		"itemsPerPage":10,
		"startIndex":0,
		"Resources":[%s, %s]
	}`, json, json), string(ctx.ResponseBody))
}

func benchmarkHybridJsonSerializationProcessor(poolSize int, b *testing.B) {
	viper.Set(fmt.Sprintf("scim.threadPool.%s", JsonHybridList), poolSize)

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
			Schema:         schema,
			MultiResults:   []resource.ScimObject{r, r},
			QueryPageStart: 1,
			QueryPageSize:  10,
			SerializationTargetFunc: func() interface{} {
				return []resource.ScimObject{r, r}
			},
		}
	}

	// serializer
	b.ResetTimer()
	processor := GetWorkerBean(JsonHybridList)
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

func BenchmarkHybridJsonSerializationProcessorWithPoolSize1(b *testing.B) {
	benchmarkHybridJsonSerializationProcessor(1, b)
}

func BenchmarkHybridJsonSerializationProcessorWithPoolSize2(b *testing.B) {
	benchmarkHybridJsonSerializationProcessor(2, b)
}

func BenchmarkHybridJsonSerializationProcessorWithPoolSize3(b *testing.B) {
	benchmarkHybridJsonSerializationProcessor(3, b)
}

func BenchmarkHybridJsonSerializationProcessorWithPoolSize4(b *testing.B) {
	benchmarkHybridJsonSerializationProcessor(4, b)
}

func BenchmarkHybridJsonSerializationProcessorWithPoolSize5(b *testing.B) {
	benchmarkHybridJsonSerializationProcessor(5, b)
}

func BenchmarkHybridJsonSerializationProcessorWithPoolSize6(b *testing.B) {
	benchmarkHybridJsonSerializationProcessor(6, b)
}

func BenchmarkHybridJsonSerializationProcessorWithPoolSize7(b *testing.B) {
	benchmarkHybridJsonSerializationProcessor(7, b)
}

func BenchmarkHybridJsonSerializationProcessorWithPoolSize8(b *testing.B) {
	benchmarkHybridJsonSerializationProcessor(8, b)
}

func BenchmarkHybridJsonSerializationProcessorWithPoolSize9(b *testing.B) {
	benchmarkHybridJsonSerializationProcessor(9, b)
}
