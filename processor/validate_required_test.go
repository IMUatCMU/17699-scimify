package processor

import (
	"context"
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/resource"
	"github.com/stretchr/testify/assert"
	"testing"
)

type requiredValidationProcessorTest struct {
	name         string
	schemaPath   string
	resourcePath string
	assertion    func(error)
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

	ctx := context.WithValue(context.Background(), resource.CK_Schema, schema)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := processor.Process(r, ctx)
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

		ctx := context.WithValue(context.Background(), resource.CK_Schema, schema)

		err = processor.Process(r, ctx)
		test.assertion(err)
	}
}
