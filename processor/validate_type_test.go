package processor

import (
	"fmt"
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/resource"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

type typeValidationProcessorTest struct {
	name         string
	resourcePath string
	assertion    func(*resource.Resource, error)
}

func benchmarkTypeValidationProcessor(poolSize int, b *testing.B) {
	viper.Set(fmt.Sprintf("scim.threadPool.%s", ValidateType), poolSize)

	sch, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		b.Fatal(err)
	}

	r, _, err := helper.LoadResource("../test_data/single_test_user_david.json")
	if err != nil {
		b.Fatal(err)
	}

	processor := GetWorkerBean(ValidateType)
	makeContext := func() *ProcessorContext {
		return &ProcessorContext{
			Schema:   sch,
			Resource: r,
		}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			context := makeContext()
			err := processor.Process(context)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkTypeValidationProcessorWithPoolSize1(b *testing.B) {
	benchmarkTypeValidationProcessor(1, b)
}

func BenchmarkTypeValidationProcessorWithPoolSize2(b *testing.B) {
	benchmarkTypeValidationProcessor(2, b)
}

func BenchmarkTypeValidationProcessorWithPoolSize3(b *testing.B) {
	benchmarkTypeValidationProcessor(3, b)
}

func BenchmarkTypeValidationProcessorWithPoolSize4(b *testing.B) {
	benchmarkTypeValidationProcessor(4, b)
}

func BenchmarkTypeValidationProcessorWithPoolSize5(b *testing.B) {
	benchmarkTypeValidationProcessor(5, b)
}

func BenchmarkTypeValidationProcessorWithPoolSize6(b *testing.B) {
	benchmarkTypeValidationProcessor(6, b)
}

func BenchmarkTypeValidationProcessorWithPoolSize7(b *testing.B) {
	benchmarkTypeValidationProcessor(7, b)
}

func BenchmarkTypeValidationProcessorWithPoolSize8(b *testing.B) {
	benchmarkTypeValidationProcessor(8, b)
}

func BenchmarkTypeValidationProcessorWithPoolSize9(b *testing.B) {
	benchmarkTypeValidationProcessor(9, b)
}

func BenchmarkTypeValidationProcessor_Process(b *testing.B) {
	processor := &typeValidationProcessor{}
	schema, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		b.Fatal(err)
	}
	schema.ConstructAttributeIndex()
	r, _, err := helper.LoadResource("../test_data/single_test_user_david.json")
	if err != nil {
		b.Fatal(err)
	}
	ctx := &ProcessorContext{Resource: r, Schema: schema}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := processor.Process(ctx)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func TestTypeValidationProcessor_Process(t *testing.T) {
	processor := &typeValidationProcessor{}
	for _, test := range []typeValidationProcessorTest{
		{
			"test valid resource",
			"../test_data/single_test_user_david.json",
			func(r *resource.Resource, err error) {
				assert.Nil(t, err)
			},
		},
		{
			"test string type has number",
			"../test_data/bad_string_type_user.json",
			func(r *resource.Resource, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, "urn:ietf:params:scim:schemas:core:2.0:User:displayName", err.(*TypeMismatchError).Attr.Assist.FullPath)
			},
		},
		{
			"test invalid datetime format",
			"../test_data/bad_datetime_format_user.json",
			func(r *resource.Resource, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, "meta.created", err.(*FormatError).Attr.Assist.FullPath)
			},
		},
		{
			"test bool type has string",
			"../test_data/bad_bool_type_user.json",
			func(r *resource.Resource, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, "urn:ietf:params:scim:schemas:core:2.0:User:active", err.(*TypeMismatchError).Attr.Assist.FullPath)
			},
		},
		{
			"test array type has string",
			"../test_data/bad_array_type_user.json",
			func(r *resource.Resource, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, "urn:ietf:params:scim:schemas:core:2.0:User:emails", err.(*TypeMismatchError).Attr.Assist.FullPath)
			},
		},
		{
			"test complex type has string",
			"../test_data/bad_complex_type_user.json",
			func(r *resource.Resource, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, "urn:ietf:params:scim:schemas:core:2.0:User:name", err.(*TypeMismatchError).Attr.Assist.FullPath)
			},
		},
		{
			"test bad partial array type",
			"../test_data/bad_partial_array_type_user.json",
			func(r *resource.Resource, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, "urn:ietf:params:scim:schemas:core:2.0:User:emails", err.(*TypeMismatchError).Attr.Assist.FullPath)
			},
		},
	} {
		// prepare schema
		schema, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
		if err != nil {
			t.Fatal(err)
		}
		schema.ConstructAttributeIndex()

		// prepare test resource
		r, _, err := helper.LoadResource(test.resourcePath)
		if err != nil {
			t.Fatal(err)
		}

		err = processor.Process(&ProcessorContext{Resource: r, Schema: schema})
		test.assertion(r, err)
	}
}
