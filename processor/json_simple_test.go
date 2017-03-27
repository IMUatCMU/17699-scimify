package processor

import (
	"fmt"
	"github.com/go-scim/scimify/helper"
	"github.com/spf13/viper"
	"testing"
)

func benchmarkSimpleJsonSerializationProcessor(poolSize int, b *testing.B) {
	viper.Set(fmt.Sprintf("scim.threadPool.%s", JsonSimple), poolSize)

	processor := GetWorkerBean(JsonSimple)
	r, _, err := helper.LoadResource("../test_data/single_test_user_david.json")
	if err != nil {
		b.Fatal(err)
	}
	makeContext := func() *ProcessorContext {
		return &ProcessorContext{
			SerializationTargetFunc: func() interface{} {
				return r
			},
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

func BenchmarkSimpleJsonSerializationProcessorWithPoolSize1(b *testing.B) {
	benchmarkSimpleJsonSerializationProcessor(1, b)
}

func BenchmarkSimpleJsonSerializationProcessorWithPoolSize2(b *testing.B) {
	benchmarkSimpleJsonSerializationProcessor(2, b)
}

func BenchmarkSimpleJsonSerializationProcessorWithPoolSize3(b *testing.B) {
	benchmarkSimpleJsonSerializationProcessor(3, b)
}

func BenchmarkSimpleJsonSerializationProcessorWithPoolSize4(b *testing.B) {
	benchmarkSimpleJsonSerializationProcessor(4, b)
}

func BenchmarkSimpleJsonSerializationProcessorWithPoolSize5(b *testing.B) {
	benchmarkSimpleJsonSerializationProcessor(5, b)
}

func BenchmarkSimpleJsonSerializationProcessorWithPoolSize6(b *testing.B) {
	benchmarkSimpleJsonSerializationProcessor(6, b)
}

func BenchmarkSimpleJsonSerializationProcessorWithPoolSize7(b *testing.B) {
	benchmarkSimpleJsonSerializationProcessor(7, b)
}

func BenchmarkSimpleJsonSerializationProcessorWithPoolSize8(b *testing.B) {
	benchmarkSimpleJsonSerializationProcessor(8, b)
}

func BenchmarkSimpleJsonSerializationProcessorWithPoolSize9(b *testing.B) {
	benchmarkSimpleJsonSerializationProcessor(9, b)
}
