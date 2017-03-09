package validation

import (
	"context"
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/resource"
	"github.com/stretchr/testify/assert"
	"testing"
)

type requiredRuleValidatorTest struct {
	name         string
	schemaPath   string
	resourcePath string
	assertion    func(bool, error)
}

func BenchmarkRequiredRulesValidator_Validate(b *testing.B) {
	validator := &requiredRulesValidator{}
	schema, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		b.Fatal(err)
	}

	r, _, err := helper.LoadResource("../test_data/single_test_user_david.json")
	if err != nil {
		b.Fatal(err)
	}

	opt := ValidationOptions{UnassignedImmutableIsIgnored: false, ReadOnlyIsMandatory: false}
	ctx := context.WithValue(context.Background(), resource.CK_Schema, schema)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := validator.Validate(r, opt, ctx)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func TestRequiredRulesValidator_Validate(t *testing.T) {
	validator := &requiredRulesValidator{}
	for _, test := range []requiredRuleValidatorTest{
		{
			"test success",
			"../test_data/test_user_schema_all.json",
			"../test_data/single_test_user_david.json",
			func(ok bool, err error) {
				assert.True(t, ok)
				assert.Nil(t, err)
			},
		},
		{
			"test missing required string attribute",
			"../test_data/required_rule_test_schema.json",
			"../test_data/missing_required_string_resource.json",
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, "a", err.(*validationError).FullPath)
			},
		},
		{
			"test missing required string array attribute",
			"../test_data/required_rule_test_schema.json",
			"../test_data/missing_required_string_array_resource.json",
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, "b", err.(*validationError).FullPath)
			},
		},
		{
			"test empty complex attribute",
			"../test_data/required_rule_test_schema.json",
			"../test_data/empty_complex_attribute_resource.json",
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, "c", err.(*validationError).FullPath)
			},
		},
		{
			"test optional non-empty complex attribute",
			"../test_data/required_rule_test_schema.json",
			"../test_data/missing_sub_attribute.json",
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, "d.d1", err.(*validationError).FullPath)
			},
		},
		{
			"test missing required sub in optional array",
			"../test_data/required_rule_test_schema.json",
			"../test_data/missing_required_sub_in_optional_array.json",
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, "e.e1", err.(*validationError).FullPath)
			},
		},
	} {
		//prepare schema
		schema, err := loadSchema(test.schemaPath)
		if err != nil {
			t.Fatal(err)
		}

		//prepare resource
		resourceData := loadTestDataFromJson(t, test.resourcePath)
		r := resource.NewResourceFromMap(resourceData)
		opt := ValidationOptions{UnassignedImmutableIsIgnored: false, ReadOnlyIsMandatory: false}
		ctx := context.WithValue(context.Background(), resource.CK_Schema, schema)

		ok, err := validator.Validate(r, opt, ctx)
		test.assertion(ok, err)
	}
}
