package processor

import (
	"fmt"
	"github.com/go-scim/scimify/resource"
	"github.com/spf13/viper"
	"testing"
)

func benchmarkErrorTranslatingProcessor(poolSize int, b *testing.B) {
	viper.Set(fmt.Sprintf("scim.threadPool.%s", TranslateError), poolSize)

	processor := GetWorkerBean(TranslateError)
	makeContext := func() *ProcessorContext {
		return &ProcessorContext{Err: resource.CreateError(resource.ServerError, "dummy")}
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

func BenchmarkErrorTranslatingProcessorWithPoolSize1(b *testing.B) {
	benchmarkErrorTranslatingProcessor(1, b)
}

func BenchmarkErrorTranslatingProcessorWithPoolSize2(b *testing.B) {
	benchmarkErrorTranslatingProcessor(2, b)
}

func BenchmarkErrorTranslatingProcessorWithPoolSize3(b *testing.B) {
	benchmarkErrorTranslatingProcessor(3, b)
}

func BenchmarkErrorTranslatingProcessorWithPoolSize4(b *testing.B) {
	benchmarkErrorTranslatingProcessor(4, b)
}

func BenchmarkErrorTranslatingProcessorWithPoolSize5(b *testing.B) {
	benchmarkErrorTranslatingProcessor(5, b)
}

func BenchmarkErrorTranslatingProcessorWithPoolSize6(b *testing.B) {
	benchmarkErrorTranslatingProcessor(6, b)
}

func BenchmarkErrorTranslatingProcessorWithPoolSize7(b *testing.B) {
	benchmarkErrorTranslatingProcessor(7, b)
}

func BenchmarkErrorTranslatingProcessorWithPoolSize8(b *testing.B) {
	benchmarkErrorTranslatingProcessor(8, b)
}

func BenchmarkErrorTranslatingProcessorWithPoolSize9(b *testing.B) {
	benchmarkErrorTranslatingProcessor(9, b)
}
