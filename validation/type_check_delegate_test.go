package validation

import (
	"github.com/go-scim/scimify/resource"
	"context"
	"github.com/go-scim/scimify/helper"
	"testing"
	"github.com/stretchr/testify/assert"
)

type testUseTypeCheckValidator struct {}

func (tcv *testUseTypeCheckValidator) Validate(r *resource.Resource, _ ValidationOptions, ctx context.Context) (pass bool, err error) {
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
	delegate := &typeCheckDelegate{}
	helper.Traverse(r, schema, []helper.ResourceTraversalDelegate{delegate})

	pass, err = true, nil
	return
}

type typeCheckDelegateTest struct {
	name         string
	resourcePath string
	assertion    func(bool, error)
}

func TestTypeCheckDelegate(t *testing.T) {
	validator := &testUseTypeCheckValidator{}
	for _, test := range []typeCheckDelegateTest{
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
				assert.Equal(t, "urn:ietf:params:scim:schemas:core:2.0:User:displayName", err.(*TypeMismatchError).Attr.Assist.FullPath)
			},
		},
		{
			"test invalid datetime format",
			"../test_data/bad_datetime_format_user.json",
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, "meta.created", err.(*FormatError).Attr.Assist.FullPath)
			},
		},
		{
			"test bool type has string",
			"../test_data/bad_bool_type_user.json",
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, "urn:ietf:params:scim:schemas:core:2.0:User:active", err.(*TypeMismatchError).Attr.Assist.FullPath)
			},
		},
		{
			"test array type has string",
			"../test_data/bad_array_type_user.json",
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, "urn:ietf:params:scim:schemas:core:2.0:User:emails", err.(*TypeMismatchError).Attr.Assist.FullPath)
			},
		},
		{
			"test complex type has string",
			"../test_data/bad_complex_type_user.json",
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, "urn:ietf:params:scim:schemas:core:2.0:User:name", err.(*TypeMismatchError).Attr.Assist.FullPath)
			},
		},
		{
			"test bad partial array type",
			"../test_data/bad_partial_array_type_user.json",
			func(ok bool, err error) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, "urn:ietf:params:scim:schemas:core:2.0:User:emails", err.(*TypeMismatchError).Attr.Assist.FullPath)
			},
		},
	}{
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

func BenchmarkTypeCheckDelegate(b *testing.B) {
	validator := &testUseTypeCheckValidator{}
	schema, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		b.Fatal(err)
	}
	schema.ConstructAttributeIndex()
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