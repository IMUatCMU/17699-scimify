package validation

import (
	"context"
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/resource"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type mutabilityCheckTest struct {
	name         string
	resourcePath string
	assertion    func(bool, error, *resource.Resource, *resource.Resource)
}

func BenchmarkMutabilityValidator_Validate(b *testing.B) {
	schema, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		b.Fatal(err)
	}

	ref, _, err := helper.LoadResource("../test_data/single_test_user_david.json")
	if err != nil {
		b.Fatal(err)
	}

	r, _, err := helper.LoadResource("../test_data/single_test_user_david.json")
	if err != nil {
		b.Fatal(err)
	}

	validator := &MutabilityValidator{}
	opt := ValidationOptions{UnassignedImmutableIsIgnored: false, ReadOnlyIsMandatory: false}

	ctx := context.Background()
	ctx = context.WithValue(context.Background(), resource.CK_Schema, schema)
	ctx = context.WithValue(ctx, resource.CK_Reference, ref)

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

func TestMutabilityValidator_Validate(t *testing.T) {
	validator := &MutabilityValidator{}

	for _, test := range []mutabilityCheckTest{
		{
			"test same resource passes",
			"../test_data/single_test_user_david.json",
			func(ok bool, err error, _ *resource.Resource, _ *resource.Resource) {
				assert.True(t, ok)
				assert.Nil(t, err)
			},
		},
		{
			"test changes readOnly string attribute",
			"../test_data/changes_readonly_string.json",
			func(ok bool, err error, _ *resource.Resource, _ *resource.Resource) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, "id", err.(*ValueChangedError).Attr.Assist.FullPath)
			},
		},
		{
			"test changes readOnly complex attribute",
			"../test_data/changes_readonly_complex.json",
			func(ok bool, err error, _ *resource.Resource, _ *resource.Resource) {
				assert.False(t, ok)
				assert.NotNil(t, err)
				assert.Equal(t, "meta", err.(*ValueChangedError).Attr.Assist.FullPath)
			},
		},
		{
			"meta was copied over",
			"../test_data/single_test_user_david_without_meta.json",
			func(ok bool, err error, r *resource.Resource, ref *resource.Resource) {
				assert.True(t, ok)
				assert.Nil(t, err)
				assert.True(t, reflect.DeepEqual(r.Attributes["meta"], ref.Attributes["meta"]))
			},
		},
	} {
		schema, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
		if err != nil {
			t.Fatal(err)
		}

		ref, _, err := helper.LoadResource("../test_data/single_test_user_david.json")
		if err != nil {
			t.Fatal(err)
		}

		r, _, err := helper.LoadResource(test.resourcePath)
		if err != nil {
			t.Fatal(err)
		}

		opt := ValidationOptions{UnassignedImmutableIsIgnored: false, ReadOnlyIsMandatory: false}

		ctx := context.Background()
		ctx = context.WithValue(ctx, resource.CK_Schema, schema)
		ctx = context.WithValue(ctx, resource.CK_Reference, ref)

		ok, err := validator.Validate(r, opt, ctx)
		test.assertion(ok, err, r, ref)
	}
}
