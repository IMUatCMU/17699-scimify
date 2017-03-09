package validation

import (
	"context"
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/resource"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testUseRequiredCheckValidator struct{}

func (rcv *testUseRequiredCheckValidator) Validate(r *resource.Resource, _ ValidationOptions, ctx context.Context) (pass bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case error:
				pass, err = false, r.(error)
				return
			default:
				panic(r)
			}
		}
	}()

	schema := ctx.Value(resource.CK_Schema).(*resource.Schema)
	delegate := &requiredCheckDelegate{enforceReadOnlyAttributes: false}
	helper.TraverseWithSchema(r, schema, []helper.ResourceTraversalDelegate{delegate})

	pass, err = true, nil
	return
}

type requiredCheckDelegateTest struct {
	name         string
	schemaPath   string
	resourcePath string
	assertion    func(bool, error)
}

func BenchmarkRequiredCheckDelegate(b *testing.B) {
	validator := &testUseRequiredCheckValidator{}
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

func TestRequiredCheckDelegate(t *testing.T) {
	validator := &testUseRequiredCheckValidator{}
	for _, test := range []requiredCheckDelegateTest{
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
				assert.Equal(t, "a", err.(*RequiredMissingError).Attr.Assist.FullPath)
			},
		},
		{
			"test missing required string array attribute",
			"../test_data/required_rule_test_schema.json",
			"../test_data/missing_required_string_array_resource.json",
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, "b", err.(*RequiredMissingError).Attr.Assist.FullPath)
			},
		},
		{
			"test empty complex attribute",
			"../test_data/required_rule_test_schema.json",
			"../test_data/empty_complex_attribute_resource.json",
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, "c", err.(*RequiredUnassignedError).Attr.Assist.FullPath)
			},
		},
		{
			"test optional non-empty complex attribute",
			"../test_data/required_rule_test_schema.json",
			"../test_data/missing_sub_attribute.json",
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, "d.d1", err.(*RequiredMissingError).Attr.Assist.FullPath)
			},
		},
		{
			"test missing required sub in optional array",
			"../test_data/required_rule_test_schema.json",
			"../test_data/missing_required_sub_in_optional_array.json",
			func(ok bool, err error) {
				assert.False(t, ok)
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

		opt := ValidationOptions{UnassignedImmutableIsIgnored: false, ReadOnlyIsMandatory: false}
		ctx := context.WithValue(context.Background(), resource.CK_Schema, schema)

		ok, err := validator.Validate(r, opt, ctx)
		test.assertion(ok, err)
	}
}
