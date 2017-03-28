package processor

import (
	"fmt"
	"github.com/go-scim/scimify/resource"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdateMetaProcessor_Process(t *testing.T) {
	r := resource.NewResourceFromMap(map[string]interface{}{
		"id": "4FD76312-456B-4233-A357-16EA035637E2",
		"meta": map[string]interface{}{
			"version": "W\\/\"a330bc54f0671c9\"",
		},
	})

	processor := &updateMetaProcessor{}
	err := processor.Process(&ProcessorContext{Resource: r})
	assert.Nil(t, err)
	assert.NotEmpty(t, r.Attributes["meta"].(map[string]interface{})["version"].(string))
	assert.True(t, len(r.Attributes["meta"].(map[string]interface{})["lastModified"].(string)) > 0)
}

func benchmarkUpdateMetaProcessor(poolSize int, b *testing.B) {
	viper.Set(fmt.Sprintf("scim.threadPool.%s", UpdateMeta), poolSize)

	processor := GetWorkerBean(UpdateMeta)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			context := &ProcessorContext{Resource: resource.NewResourceFromMap(map[string]interface{}{
				"id": "4FD76312-456B-4233-A357-16EA035637E2",
				"meta": map[string]interface{}{
					"version": "W\\/\"a330bc54f0671c9\"",
				},
			})}

			err := processor.Process(context)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkUpdateMetaProcessorWithPoolSize1(b *testing.B) {
	benchmarkUpdateMetaProcessor(1, b)
}

func BenchmarkUpdateMetaProcessorWithPoolSize2(b *testing.B) {
	benchmarkUpdateMetaProcessor(2, b)
}

func BenchmarkUpdateMetaProcessorWithPoolSize3(b *testing.B) {
	benchmarkUpdateMetaProcessor(3, b)
}

func BenchmarkUpdateMetaProcessorWithPoolSize4(b *testing.B) {
	benchmarkUpdateMetaProcessor(4, b)
}

func BenchmarkUpdateMetaProcessorWithPoolSize5(b *testing.B) {
	benchmarkUpdateMetaProcessor(5, b)
}

func BenchmarkUpdateMetaProcessorWithPoolSize6(b *testing.B) {
	benchmarkUpdateMetaProcessor(6, b)
}

func BenchmarkUpdateMetaProcessorWithPoolSize7(b *testing.B) {
	benchmarkUpdateMetaProcessor(7, b)
}

func BenchmarkUpdateMetaProcessorWithPoolSize8(b *testing.B) {
	benchmarkUpdateMetaProcessor(8, b)
}

func BenchmarkUpdateMetaProcessorWithPoolSize9(b *testing.B) {
	benchmarkUpdateMetaProcessor(9, b)
}
