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

type formatCaseProcessorTest struct {
	name         string
	resourcePath string
	schemaPath   string
	assertion    func(r *resource.Resource, err error)
}

func benchmarkFormatCaseProcessor(poolSize int, b *testing.B) {
	viper.Set(fmt.Sprintf("scim.threadPool.%s", FormatCase), poolSize)

	sch, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if nil != err {
		b.Fatal(err)
	}
	sch.ConstructAttributeIndex()

	processor := GetWorkerBean(FormatCase)
	makeContext := func() *ProcessorContext {
		return &ProcessorContext{Schema: sch}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r, _, err := helper.LoadResource("../test_data/single_test_user_david_capitalized.json")
			if nil != err {
				b.Fatal(err)
			}
			context := makeContext()
			context.Resource = r

			err = processor.Process(context)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkFormatCaseProcessorWithPoolSize1(b *testing.B) {
	benchmarkFormatCaseProcessor(1, b)
}

func BenchmarkFormatCaseProcessorWithPoolSize2(b *testing.B) {
	benchmarkFormatCaseProcessor(2, b)
}

func BenchmarkFormatCaseProcessorWithPoolSize3(b *testing.B) {
	benchmarkFormatCaseProcessor(3, b)
}

func BenchmarkFormatCaseProcessorWithPoolSize4(b *testing.B) {
	benchmarkFormatCaseProcessor(4, b)
}

func BenchmarkFormatCaseProcessorWithPoolSize5(b *testing.B) {
	benchmarkFormatCaseProcessor(5, b)
}

func BenchmarkFormatCaseProcessorWithPoolSize6(b *testing.B) {
	benchmarkFormatCaseProcessor(6, b)
}

func BenchmarkFormatCaseProcessorWithPoolSize7(b *testing.B) {
	benchmarkFormatCaseProcessor(7, b)
}

func BenchmarkFormatCaseProcessorWithPoolSize8(b *testing.B) {
	benchmarkFormatCaseProcessor(8, b)
}

func BenchmarkFormatCaseProcessorWithPoolSize9(b *testing.B) {
	benchmarkFormatCaseProcessor(9, b)
}

func BenchmarkFormatCaseProcessor_Process(b *testing.B) {
	processor := &formatCaseProcessor{}

	sch, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if nil != err {
		b.Fatal(err)
	}
	sch.ConstructAttributeIndex()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			b.StopTimer()
			r, _, err := helper.LoadResource("../test_data/single_test_user_david_capitalized.json")
			if nil != err {
				b.Fatal(err)
			}
			b.StartTimer()

			err = processor.Process(&ProcessorContext{Resource: r, Schema: sch})
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

		err = processor.Process(&ProcessorContext{Resource: r, Schema: sch})
		test.assertion(r, err)
	}
}
