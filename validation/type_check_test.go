package validation

import (
	"github.com/go-scim/scimify/resource"
	"context"
	"testing"
	"github.com/go-scim/scimify/helper"
	"github.com/stretchr/testify/assert"
)

type typesRuleValidatorTest struct {
	name         string
	resourcePath string
	assertion    func(bool, error)
}

func TestTypeCheckValidator_Validate(t *testing.T) {
	validator := &typeCheckValidator{}
	for _, test := range []typesRuleValidatorTest{
		{
			"test valid resource",
			"../test_data/single_test_user_david.json",
			func(ok bool, err error) {
				assert.True(t, ok)
				assert.Nil(t, err)
			},
		},
		{
			"test string type has number",
			"../test_data/bad_string_type_user.json",
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, "urn:ietf:params:scim:schemas:core:2.0:User:displayName", err.(*UnsupportedTypeError).Attr.Assist.FullPath)
			},
		},
		// TODO
		//{
		//	"test invalid datetime format",
		//	"../test_data/bad_datetime_format_user.json",
		//	func(ok bool, err error) {
		//		assert.False(t, ok)
		//		assert.NotNil(t, err)
		//		assert.Equal(t, "meta.created", err.(*UnsupportedTypeError).Attr.Assist.FullPath)
		//	},
		//},
		{
			"test bool type has string",
			"../test_data/bad_bool_type_user.json",
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, "urn:ietf:params:scim:schemas:core:2.0:User:active", err.(*UnsupportedTypeError).Attr.Assist.FullPath)
			},
		},
		{
			"test array type has string",
			"../test_data/bad_array_type_user.json",
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, "urn:ietf:params:scim:schemas:core:2.0:User:emails", err.(*UnsupportedTypeError).Attr.Assist.FullPath)
			},
		},
		{
			"test complex type has string",
			"../test_data/bad_complex_type_user.json",
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, "urn:ietf:params:scim:schemas:core:2.0:User:name", err.(*UnsupportedTypeError).Attr.Assist.FullPath)
			},
		},
		{
			"test bad partial array type",
			"../test_data/bad_partial_array_type_user.json",
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, "urn:ietf:params:scim:schemas:core:2.0:User:emails", err.(*UnsupportedTypeError).Attr.Assist.FullPath)
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

		opt := ValidationOptions{UnassignedImmutableIsIgnored: false, ReadOnlyIsMandatory: false}
		ctx := context.WithValue(context.Background(), resource.CK_Schema, schema)

		ok, err := validator.Validate(r, opt, ctx)
		test.assertion(ok, err)
	}
}