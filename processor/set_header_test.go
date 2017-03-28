package processor

import (
	"fmt"
	"github.com/go-scim/scimify/helper"
	"github.com/spf13/viper"
	"testing"
)

func benchmarkSetAllHeaderProcessor(poolSize int, b *testing.B) {
	viper.Set(fmt.Sprintf("scim.threadPool.%s", SetAllHeader), poolSize)

	r, _, err := helper.LoadResource("../test_data/single_test_user_david.json")
	if err != nil {
		b.Fatal(err)
	}

	processor := GetWorkerBean(SetAllHeader)
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

func BenchmarkSetAllHeaderProcessorWithPoolSize1(b *testing.B) {
	benchmarkSetAllHeaderProcessor(1, b)
}

func BenchmarkSetAllHeaderProcessorWithPoolSize2(b *testing.B) {
	benchmarkSetAllHeaderProcessor(2, b)
}

func BenchmarkSetAllHeaderProcessorWithPoolSize3(b *testing.B) {
	benchmarkSetAllHeaderProcessor(3, b)
}

func BenchmarkSetAllHeaderProcessorWithPoolSize4(b *testing.B) {
	benchmarkSetAllHeaderProcessor(4, b)
}

func BenchmarkSetAllHeaderProcessorWithPoolSize5(b *testing.B) {
	benchmarkSetAllHeaderProcessor(5, b)
}

func BenchmarkSetAllHeaderProcessorWithPoolSize6(b *testing.B) {
	benchmarkSetAllHeaderProcessor(6, b)
}

func BenchmarkSetAllHeaderProcessorWithPoolSize7(b *testing.B) {
	benchmarkSetAllHeaderProcessor(7, b)
}

func BenchmarkSetAllHeaderProcessorWithPoolSize8(b *testing.B) {
	benchmarkSetAllHeaderProcessor(8, b)
}

func BenchmarkSetAllHeaderProcessorWithPoolSize9(b *testing.B) {
	benchmarkSetAllHeaderProcessor(9, b)
}
