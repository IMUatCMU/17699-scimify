package processor

import (
	"fmt"
	"github.com/spf13/viper"
	"testing"
)

func benchmarkSetStatusToOKProcessor(poolSize int, b *testing.B) {
	viper.Set(fmt.Sprintf("scim.threadPool.%s", SetStatusToOk), poolSize)

	processor := GetWorkerBean(SetStatusToOk)
	makeContext := func() *ProcessorContext {
		return &ProcessorContext{}
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

func BenchmarkSetStatusToOKProcessorWithPoolSize1(b *testing.B) {
	benchmarkSetStatusToOKProcessor(1, b)
}

func BenchmarkSetStatusToOKProcessorWithPoolSize2(b *testing.B) {
	benchmarkSetStatusToOKProcessor(2, b)
}

func BenchmarkSetStatusToOKProcessorWithPoolSize3(b *testing.B) {
	benchmarkSetStatusToOKProcessor(3, b)
}

func BenchmarkSetStatusToOKProcessorWithPoolSize4(b *testing.B) {
	benchmarkSetStatusToOKProcessor(4, b)
}

func BenchmarkSetStatusToOKProcessorWithPoolSize5(b *testing.B) {
	benchmarkSetStatusToOKProcessor(5, b)
}

func BenchmarkSetStatusToOKProcessorWithPoolSize6(b *testing.B) {
	benchmarkSetStatusToOKProcessor(6, b)
}

func BenchmarkSetStatusToOKProcessorWithPoolSize7(b *testing.B) {
	benchmarkSetStatusToOKProcessor(7, b)
}

func BenchmarkSetStatusToOKProcessorWithPoolSize8(b *testing.B) {
	benchmarkSetStatusToOKProcessor(8, b)
}

func BenchmarkSetStatusToOKProcessorWithPoolSize9(b *testing.B) {
	benchmarkSetStatusToOKProcessor(9, b)
}
