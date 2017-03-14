package worker

import (
	"context"
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/processor"
	"github.com/go-scim/scimify/resource"
	"github.com/stretchr/testify/assert"
	"testing"
)

func BenchmarkGetUserCreationProcessorWorker(b *testing.B) {
	w := GetUserCreationProcessorWorker()

	sch, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		b.Fatal(err)
	}
	sch.ConstructAttributeIndex()

	ctx := context.Background()
	ctx = context.WithValue(ctx, resource.CK_Schema, sch)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			b.StopTimer()
			r, _, err := helper.LoadResource("../test_data/user_creation_correct.json")
			if err != nil {
				b.Fatal(err)
			}
			b.StartTimer()

			_, err = w.Do(&ProcessorInput{R: r, Ctx: ctx})
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkGetUserUpdateProcessorWorker(b *testing.B) {
	w := GetUserUpdateProcessorWorker()

	ref, _, err := helper.LoadResource("../test_data/single_test_user_david.json")
	if err != nil {
		b.Fatal(err)
	}

	sch, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		b.Fatal(err)
	}
	sch.ConstructAttributeIndex()

	ctx := context.Background()
	ctx = context.WithValue(ctx, resource.CK_Reference, ref)
	ctx = context.WithValue(ctx, resource.CK_Schema, sch)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			b.StopTimer()
			r, _, err := helper.LoadResource("../test_data/user_update_correct.json")
			if err != nil {
				b.Fatal(err)
			}
			b.StartTimer()

			_, err = w.Do(&ProcessorInput{R: r, Ctx: ctx})
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func TestGetUserCreationProcessorWorker(t *testing.T) {
	w := GetUserCreationProcessorWorker()

	for _, test := range []struct {
		name         string
		resourcePath string
		assertion    func(*resource.Resource, error)
	}{
		{
			"create correct",
			"../test_data/user_creation_correct.json",
			func(r *resource.Resource, err error) {
				assert.Nil(t, err)
				assert.NotEmpty(t, r.Data()["id"])
				assert.NotNil(t, r.Data()["meta"])
			},
		},
		{
			"create missing required attribute schema",
			"../test_data/user_creation_missing_schema.json",
			func(r *resource.Resource, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, "schemas", err.(*processor.RequiredMissingError).Attr.Assist.FullPath)
			},
		},
		{
			"create with wrong type",
			"../test_data/user_creation_wrong_type.json",
			func(r *resource.Resource, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, "externalId", err.(*processor.TypeMismatchError).Attr.Assist.FullPath)
			},
		},
	} {
		r, _, err := helper.LoadResource(test.resourcePath)
		if err != nil {
			t.Fatal(err)
		}

		sch, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
		if err != nil {
			t.Fatal(err)
		}
		sch.ConstructAttributeIndex()

		ctx := context.Background()
		ctx = context.WithValue(ctx, resource.CK_Schema, sch)

		_, err = w.Do(&ProcessorInput{R: r, Ctx: ctx})
		test.assertion(r, err)
	}
}

func TestGetUserUpdateProcessorWorker(t *testing.T) {
	w := GetUserUpdateProcessorWorker()

	for _, test := range []struct {
		name         string
		resourcePath string
		assertion    func(*resource.Resource, error)
	}{
		{
			"update correct",
			"../test_data/user_update_correct.json",
			func(r *resource.Resource, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, r.Data()["meta"])
			},
		},
		{
			"update change readonly",
			"../test_data/user_update_change_readonly.json",
			func(r *resource.Resource, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, "id", err.(*processor.ValueChangedError).Attr.Assist.FullPath)
			},
		},
	} {
		r, _, err := helper.LoadResource(test.resourcePath)
		if err != nil {
			t.Fatal(err)
		}

		ref, _, err := helper.LoadResource("../test_data/single_test_user_david.json")
		if err != nil {
			t.Fatal(err)
		}

		sch, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
		if err != nil {
			t.Fatal(err)
		}
		sch.ConstructAttributeIndex()

		ctx := context.Background()
		ctx = context.WithValue(ctx, resource.CK_Reference, ref)
		ctx = context.WithValue(ctx, resource.CK_Schema, sch)

		_, err = w.Do(&ProcessorInput{R: r, Ctx: ctx})
		test.assertion(r, err)
	}
}
