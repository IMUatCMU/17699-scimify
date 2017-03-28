package processor

import (
	"fmt"
	"github.com/go-scim/scimify/resource"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateIdProcessor_Process(t *testing.T) {
	processor := &generateIdProcessor{}
	r := resource.NewResource()
	err := processor.Process(&ProcessorContext{Resource: r})
	assert.Nil(t, err)
	assert.NotEqual(t, 0, len(r.Attributes["id"].(string)))
}

func benchmarkGenerateIdProcessor(poolSize int, b *testing.B) {
	viper.Set(fmt.Sprintf("scim.threadPool.%s", GenerateId), poolSize)

	processor := GetWorkerBean(GenerateId)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			context := &ProcessorContext{Resource: resource.NewResource()}

			err := processor.Process(context)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkGenerateIdProcessorWithPoolSize1(b *testing.B) {
	benchmarkGenerateIdProcessor(1, b)
}

func BenchmarkGenerateIdProcessorWithPoolSize2(b *testing.B) {
	benchmarkGenerateIdProcessor(2, b)
}

func BenchmarkGenerateIdProcessorWithPoolSize3(b *testing.B) {
	benchmarkGenerateIdProcessor(3, b)
}

func BenchmarkGenerateIdProcessorWithPoolSize4(b *testing.B) {
	benchmarkGenerateIdProcessor(4, b)
}

func BenchmarkGenerateIdProcessorWithPoolSize5(b *testing.B) {
	benchmarkGenerateIdProcessor(5, b)
}

func BenchmarkGenerateIdProcessorWithPoolSize6(b *testing.B) {
	benchmarkGenerateIdProcessor(6, b)
}

func BenchmarkGenerateIdProcessorWithPoolSize7(b *testing.B) {
	benchmarkGenerateIdProcessor(7, b)
}

func BenchmarkGenerateIdProcessorWithPoolSize8(b *testing.B) {
	benchmarkGenerateIdProcessor(8, b)
}

func BenchmarkGenerateIdProcessorWithPoolSize9(b *testing.B) {
	benchmarkGenerateIdProcessor(9, b)
}
