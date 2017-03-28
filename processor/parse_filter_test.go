package processor

import (
	"fmt"
	"github.com/go-scim/scimify/helper"
	"github.com/spf13/viper"
	"testing"
)

func benchmarkParseFilterProcessor(poolSize int, b *testing.B) {
	viper.Set(fmt.Sprintf("scim.threadPool.%s", ParseFilter), poolSize)

	sch, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		b.Fatal(err)
	}
	sch.ConstructAttributeIndex()

	processor := GetWorkerBean(ParseFilter)
	makeContext := func() *ProcessorContext {
		return &ProcessorContext{
			QueryFilter: "username eq \"david\" and x509Certificates.value eq \"dummy\"",
			Schema:      sch,
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

func BenchmarkParseFilterProcessorWithPoolSize1(b *testing.B) {
	benchmarkParseFilterProcessor(1, b)
}

func BenchmarkParseFilterProcessorWithPoolSize2(b *testing.B) {
	benchmarkParseFilterProcessor(2, b)
}

func BenchmarkParseFilterProcessorWithPoolSize3(b *testing.B) {
	benchmarkParseFilterProcessor(3, b)
}

func BenchmarkParseFilterProcessorWithPoolSize4(b *testing.B) {
	benchmarkParseFilterProcessor(4, b)
}

func BenchmarkParseFilterProcessorWithPoolSize5(b *testing.B) {
	benchmarkParseFilterProcessor(5, b)
}

func BenchmarkParseFilterProcessorWithPoolSize6(b *testing.B) {
	benchmarkParseFilterProcessor(6, b)
}

func BenchmarkParseFilterProcessorWithPoolSize7(b *testing.B) {
	benchmarkParseFilterProcessor(7, b)
}

func BenchmarkParseFilterProcessorWithPoolSize8(b *testing.B) {
	benchmarkParseFilterProcessor(8, b)
}

func BenchmarkParseFilterProcessorWithPoolSize9(b *testing.B) {
	benchmarkParseFilterProcessor(9, b)
}
