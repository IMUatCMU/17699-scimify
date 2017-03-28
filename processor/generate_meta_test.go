package processor

import (
	"fmt"
	"github.com/go-scim/scimify/resource"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateMetaProcessor_Process(t *testing.T) {
	viper.Set("server.rootPath", "http://foo.com/v2/")

	r := resource.NewResourceFromMap(map[string]interface{}{
		"id": "bar",
	})

	processor := &generateMetaProcessor{
		ResourceType:    "User",
		ResourceTypeUri: "/User",
	}
	err := processor.Process(&ProcessorContext{Resource: r})
	assert.Nil(t, err)

	meta, ok := r.Attributes["meta"].(map[string]interface{})
	assert.True(t, ok)
	assert.NotNil(t, meta)
	assert.Equal(t, "User", meta["resourceType"].(string))
	assert.Equal(t, "http://foo.com/v2/User/bar", meta["location"].(string))
	assert.NotEmpty(t, meta["version"].(string))
	assert.NotEmpty(t, meta["created"].(string))
	assert.NotEmpty(t, meta["lastModified"].(string))
}

func benchmarkGenerateMetaProcessor(poolSize int, b *testing.B) {
	viper.Set("server.rootPath", "http://foo.com/v2/")
	viper.Set(fmt.Sprintf("scim.threadPool.%s", GenerateUserMeta), poolSize)

	processor := GetWorkerBean(GenerateUserMeta)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			context := &ProcessorContext{Resource: resource.NewResourceFromMap(map[string]interface{}{
				"id": uuid.NewV4().String(),
			})}

			err := processor.Process(context)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkGenerateMetaProcessorWithPoolSize1(b *testing.B) {
	benchmarkGenerateMetaProcessor(1, b)
}

func BenchmarkGenerateMetaProcessorWithPoolSize2(b *testing.B) {
	benchmarkGenerateMetaProcessor(2, b)
}

func BenchmarkGenerateMetaProcessorWithPoolSize3(b *testing.B) {
	benchmarkGenerateMetaProcessor(3, b)
}

func BenchmarkGenerateMetaProcessorWithPoolSize4(b *testing.B) {
	benchmarkGenerateMetaProcessor(4, b)
}

func BenchmarkGenerateMetaProcessorWithPoolSize5(b *testing.B) {
	benchmarkGenerateMetaProcessor(5, b)
}

func BenchmarkGenerateMetaProcessorWithPoolSize6(b *testing.B) {
	benchmarkGenerateMetaProcessor(6, b)
}

func BenchmarkGenerateMetaProcessorWithPoolSize7(b *testing.B) {
	benchmarkGenerateMetaProcessor(7, b)
}

func BenchmarkGenerateMetaProcessorWithPoolSize8(b *testing.B) {
	benchmarkGenerateMetaProcessor(8, b)
}

func BenchmarkGenerateMetaProcessorWithPoolSize9(b *testing.B) {
	benchmarkGenerateMetaProcessor(9, b)
}
