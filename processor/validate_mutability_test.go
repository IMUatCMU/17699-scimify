package processor

import (
	"context"
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/resource"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type mutabilityValidationProcessorTest struct {
	name         string
	resourcePath string
	assertion    func(err error, r *resource.Resource, ref *resource.Resource)
}

func BenchmarkMutabilityValidationProcessor_Process(b *testing.B) {
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

	processor := &mutabilityValidationProcessor{}

	ctx := context.Background()
	ctx = context.WithValue(context.Background(), resource.CK_Schema, schema)
	ctx = context.WithValue(ctx, resource.CK_Reference, ref)

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

func TestMutabilityValidationProcessor_Process(t *testing.T) {
	processor := &mutabilityValidationProcessor{}

	for _, test := range []mutabilityValidationProcessorTest{
		{
			"test same resource passes",
			"../test_data/single_test_user_david.json",
			func(err error, _ *resource.Resource, _ *resource.Resource) {
				assert.Nil(t, err)
			},
		},
		{
			"test changes readOnly string attribute",
			"../test_data/changes_readonly_string.json",
			func(err error, _ *resource.Resource, _ *resource.Resource) {
				assert.NotNil(t, err)
				assert.Equal(t, "id", err.(*ValueChangedError).Attr.Assist.FullPath)
			},
		},
		{
			"test changes readOnly complex attribute",
			"../test_data/changes_readonly_complex.json",
			func(err error, _ *resource.Resource, _ *resource.Resource) {
				assert.NotNil(t, err)
				assert.Equal(t, "meta", err.(*ValueChangedError).Attr.Assist.FullPath)
			},
		},
		{
			"meta was copied over",
			"../test_data/single_test_user_david_without_meta.json",
			func(err error, r *resource.Resource, ref *resource.Resource) {
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

		ctx := context.Background()
		ctx = context.WithValue(ctx, resource.CK_Schema, schema)
		ctx = context.WithValue(ctx, resource.CK_Reference, ref)

		err = processor.Process(r, ctx)
		test.assertion(err, r, ref)
	}
}
