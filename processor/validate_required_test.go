package processor

import (
	"fmt"
	"github.com/go-scim/scimify/helper"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

type requiredValidationProcessorTest struct {
	name         string
	schemaPath   string
	resourcePath string
	assertion    func(error)
}

func benchmarkRequiredValidationProcessor(poolSize int, b *testing.B) {
	viper.Set(fmt.Sprintf("scim.threadPool.%s", ValidateRequired), poolSize)

	sch, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		b.Fatal(err)
	}

	r, _, err := helper.LoadResource("../test_data/single_test_user_david.json")
	if err != nil {
		b.Fatal(err)
	}

	processor := GetWorkerBean(ValidateRequired)
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

func BenchmarkRequiredValidationProcessorWithPoolSize1(b *testing.B) {
	benchmarkRequiredValidationProcessor(1, b)
}

func BenchmarkRequiredValidationProcessorWithPoolSize2(b *testing.B) {
	benchmarkRequiredValidationProcessor(2, b)
}

func BenchmarkRequiredValidationProcessorWithPoolSize3(b *testing.B) {
	benchmarkRequiredValidationProcessor(3, b)
}

func BenchmarkRequiredValidationProcessorWithPoolSize4(b *testing.B) {
	benchmarkRequiredValidationProcessor(4, b)
}

func BenchmarkRequiredValidationProcessorWithPoolSize5(b *testing.B) {
	benchmarkRequiredValidationProcessor(5, b)
}

func BenchmarkRequiredValidationProcessorWithPoolSize6(b *testing.B) {
	benchmarkRequiredValidationProcessor(6, b)
}

func BenchmarkRequiredValidationProcessorWithPoolSize7(b *testing.B) {
	benchmarkRequiredValidationProcessor(7, b)
}

func BenchmarkRequiredValidationProcessorWithPoolSize8(b *testing.B) {
	benchmarkRequiredValidationProcessor(8, b)
}

func BenchmarkRequiredValidationProcessorWithPoolSize9(b *testing.B) {
	benchmarkRequiredValidationProcessor(9, b)
}

func BenchmarkRequiredValidationProcessor_Process(b *testing.B) {
	processor := &requiredValidationProcessor{}
	schema, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		b.Fatal(err)
	}

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

func TestRequiredValidationProcessor_Process(t *testing.T) {
	processor := &requiredValidationProcessor{}
	for _, test := range []requiredValidationProcessorTest{
		{
			"test success",
			"../test_data/test_user_schema_all.json",
			"../test_data/single_test_user_david.json",
			func(err error) {
				assert.Nil(t, err)
			},
		},
		{
			"test missing required string attribute",
			"../test_data/required_rule_test_schema.json",
			"../test_data/missing_required_string_resource.json",
			func(err error) {
				assert.NotNil(t, err)
				assert.Equal(t, "a", err.(*RequiredMissingError).Attr.Assist.FullPath)
			},
		},
		{
			"test missing required string array attribute",
			"../test_data/required_rule_test_schema.json",
			"../test_data/missing_required_string_array_resource.json",
			func(err error) {
				assert.NotNil(t, err)
				assert.Equal(t, "b", err.(*RequiredMissingError).Attr.Assist.FullPath)
			},
		},
		{
			"test empty complex attribute",
			"../test_data/required_rule_test_schema.json",
			"../test_data/empty_complex_attribute_resource.json",
			func(err error) {
				assert.NotNil(t, err)
				assert.Equal(t, "c", err.(*RequiredUnassignedError).Attr.Assist.FullPath)
			},
		},
		{
			"test optional non-empty complex attribute",
			"../test_data/required_rule_test_schema.json",
			"../test_data/missing_sub_attribute.json",
			func(err error) {
				assert.NotNil(t, err)
				assert.Equal(t, "d.d1", err.(*RequiredMissingError).Attr.Assist.FullPath)
			},
		},
		{
			"test missing required sub in optional array",
			"../test_data/required_rule_test_schema.json",
			"../test_data/missing_required_sub_in_optional_array.json",
			func(err error) {
				assert.NotNil(t, err)
				assert.Equal(t, "e.e1", err.(*RequiredMissingError).Attr.Assist.FullPath)
			},
		},
	} {
		schema, _, err := helper.LoadSchema(test.schemaPath)
		if err != nil {
			t.Fatal(err)
		}

		r, _, err := helper.LoadResource(test.resourcePath)
		if err != nil {
			t.Fatal(err)
		}

		ctx := &ProcessorContext{Resource: r, Schema: schema}

		err = processor.Process(ctx)
		test.assertion(err)
	}
}
