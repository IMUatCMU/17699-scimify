package processor

import (
	"fmt"
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/resource"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type mutabilityValidationProcessorTest struct {
	name         string
	resourcePath string
	assertion    func(err error, r *resource.Resource, ref *resource.Resource)
}

func benchmarkMutabilityValidationProcessor(poolSize int, b *testing.B) {
	viper.Set(fmt.Sprintf("scim.threadPool.%s", ValidateMutability), poolSize)

	sch, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		b.Fatal(err)
	}

	r, _, err := helper.LoadResource("../test_data/single_test_user_david.json")
	if err != nil {
		b.Fatal(err)
	}

	processor := GetWorkerBean(ValidateMutability)
	makeContext := func() *ProcessorContext {
		return &ProcessorContext{
			Schema:    sch,
			Resource:  r,
			Reference: r,
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

func BenchmarkMutabilityValidationProcessorWithPoolSize1(b *testing.B) {
	benchmarkMutabilityValidationProcessor(1, b)
}

func BenchmarkMutabilityValidationProcessorWithPoolSize2(b *testing.B) {
	benchmarkMutabilityValidationProcessor(2, b)
}

func BenchmarkMutabilityValidationProcessorWithPoolSize3(b *testing.B) {
	benchmarkMutabilityValidationProcessor(3, b)
}

func BenchmarkMutabilityValidationProcessorWithPoolSize4(b *testing.B) {
	benchmarkMutabilityValidationProcessor(4, b)
}

func BenchmarkMutabilityValidationProcessorWithPoolSize5(b *testing.B) {
	benchmarkMutabilityValidationProcessor(5, b)
}

func BenchmarkMutabilityValidationProcessorWithPoolSize6(b *testing.B) {
	benchmarkMutabilityValidationProcessor(6, b)
}

func BenchmarkMutabilityValidationProcessorWithPoolSize7(b *testing.B) {
	benchmarkMutabilityValidationProcessor(7, b)
}

func BenchmarkMutabilityValidationProcessorWithPoolSize8(b *testing.B) {
	benchmarkMutabilityValidationProcessor(8, b)
}

func BenchmarkMutabilityValidationProcessorWithPoolSize9(b *testing.B) {
	benchmarkMutabilityValidationProcessor(9, b)
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

	ctx := &ProcessorContext{Resource: r, Reference: ref, Schema: schema}

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

		ctx := &ProcessorContext{Resource: r, Reference: ref, Schema: schema}

		err = processor.Process(ctx)
		test.assertion(err, r, ref)
	}
}
