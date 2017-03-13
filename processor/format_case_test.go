package processor

import (
	"context"
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/resource"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type formatCaseProcessorTest struct {
	name         string
	resourcePath string
	schemaPath   string
	assertion    func(r *resource.Resource, err error)
}

func BenchmarkFormatCaseProcessor_Process(b *testing.B) {
	processor := &formatCaseProcessor{}

	sch, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if nil != err {
		b.Fatal(err)
	}
	sch.ConstructAttributeIndex()

	ctx := context.Background()
	ctx = context.WithValue(ctx, resource.CK_Schema, sch)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			b.StopTimer()
			r, _, err := helper.LoadResource("../test_data/single_test_user_david_capitalized.json")
			if nil != err {
				b.Fatal(err)
			}
			b.StartTimer()

			err = processor.Process(r, ctx)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func TestFormatCaseProcessor_Process(t *testing.T) {
	processor := &formatCaseProcessor{}

	for _, test := range []formatCaseProcessorTest{
		{
			"test identical resource",
			"../test_data/single_test_user_david.json",
			"../test_data/test_user_schema_all.json",
			func(r *resource.Resource, err error) {
				assert.Nil(t, err)
				expected, _, _ := helper.LoadResource("../test_data/single_test_user_david.json")
				assert.True(t, reflect.DeepEqual(r.Data(), expected.Data()))
			},
		},
		{
			"test capitalized case",
			"../test_data/single_test_user_david_capitalized.json",
			"../test_data/test_user_schema_all.json",
			func(r *resource.Resource, err error) {
				assert.Nil(t, err)
				expected, _, _ := helper.LoadResource("../test_data/single_test_user_david.json")
				assert.True(t, reflect.DeepEqual(r.Data(), expected.Data()))
			},
		},
	} {
		r, _, err := helper.LoadResource(test.resourcePath)
		if nil != err {
			t.Fatal(err)
		}

		sch, _, err := helper.LoadSchema(test.schemaPath)
		if nil != err {
			t.Fatal(err)
		}
		sch.ConstructAttributeIndex()

		ctx := context.Background()
		ctx = context.WithValue(ctx, resource.CK_Schema, sch)

		err = processor.Process(r, ctx)
		test.assertion(r, err)
	}
}
