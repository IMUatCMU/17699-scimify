package processor

import (
	"fmt"
	"github.com/go-scim/scimify/helper"
	"github.com/spf13/viper"
	"testing"
)

func benchmarkSingleResultAsJsonTargetProcessor(poolSize int, b *testing.B) {
	viper.Set(fmt.Sprintf("scim.threadPool.%s", SetJsonToSingle), poolSize)

	r, _, err := helper.LoadResource("../test_data/single_test_user_david.json")
	if err != nil {
		b.Fatal(err)
	}

	processor := GetWorkerBean(SetJsonToSingle)
	makeContext := func() *ProcessorContext {
		return &ProcessorContext{SingleResult: r}
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

func BenchmarkSingleResultAsJsonTargetProcessorWithPoolSize1(b *testing.B) {
	benchmarkSingleResultAsJsonTargetProcessor(1, b)
}

func BenchmarkSingleResultAsJsonTargetProcessorWithPoolSize2(b *testing.B) {
	benchmarkSingleResultAsJsonTargetProcessor(2, b)
}

func BenchmarkSingleResultAsJsonTargetProcessorWithPoolSize3(b *testing.B) {
	benchmarkSingleResultAsJsonTargetProcessor(3, b)
}

func BenchmarkSingleResultAsJsonTargetProcessorWithPoolSize4(b *testing.B) {
	benchmarkSingleResultAsJsonTargetProcessor(4, b)
}

func BenchmarkSingleResultAsJsonTargetProcessorWithPoolSize5(b *testing.B) {
	benchmarkSingleResultAsJsonTargetProcessor(5, b)
}

func BenchmarkSingleResultAsJsonTargetProcessorWithPoolSize6(b *testing.B) {
	benchmarkSingleResultAsJsonTargetProcessor(6, b)
}

func BenchmarkSingleResultAsJsonTargetProcessorWithPoolSize7(b *testing.B) {
	benchmarkSingleResultAsJsonTargetProcessor(7, b)
}

func BenchmarkSingleResultAsJsonTargetProcessorWithPoolSize8(b *testing.B) {
	benchmarkSingleResultAsJsonTargetProcessor(8, b)
}

func BenchmarkSingleResultAsJsonTargetProcessorWithPoolSize9(b *testing.B) {
	benchmarkSingleResultAsJsonTargetProcessor(9, b)
}
