package processor

import (
	"fmt"
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/modify"
	"github.com/go-scim/scimify/resource"
	"github.com/spf13/viper"
	"testing"
)

func benchmarkModificationProcessor(poolSize int, b *testing.B) {
	viper.Set(fmt.Sprintf("scim.threadPool.%s", Modification), poolSize)

	sch, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if nil != err {
		b.Fatal(err)
	}
	sch.ConstructAttributeIndex()

	mod := &modify.Modification{
		Schemas: []string{resource.PathOpUrn},
		Operations: []modify.ModUnit{
			{Op: "add", Path: "userName", Value: "david"},
			{Op: "remove", Path: "active"},
		},
	}

	processor := GetWorkerBean(Modification)
	makeContext := func() *ProcessorContext {
		return &ProcessorContext{Schema: sch, Mod: mod}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			context := makeContext()
			context.Resource = resource.NewResourceFromMap(map[string]interface{}{
				"id": "4FD76312-456B-4233-A357-16EA035637E2",
				"meta": map[string]interface{}{
					"version": "W\\/\"a330bc54f0671c9\"",
				},
			})

			err = processor.Process(context)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkModificationProcessorWithPoolSize1(b *testing.B) {
	benchmarkModificationProcessor(1, b)
}

func BenchmarkModificationProcessorWithPoolSize2(b *testing.B) {
	benchmarkModificationProcessor(2, b)
}

func BenchmarkModificationProcessorWithPoolSize3(b *testing.B) {
	benchmarkModificationProcessor(3, b)
}

func BenchmarkModificationProcessorWithPoolSize4(b *testing.B) {
	benchmarkModificationProcessor(4, b)
}

func BenchmarkModificationProcessorWithPoolSize5(b *testing.B) {
	benchmarkModificationProcessor(5, b)
}

func BenchmarkModificationProcessorWithPoolSize6(b *testing.B) {
	benchmarkModificationProcessor(6, b)
}

func BenchmarkModificationProcessorWithPoolSize7(b *testing.B) {
	benchmarkModificationProcessor(7, b)
}

func BenchmarkModificationProcessorWithPoolSize8(b *testing.B) {
	benchmarkModificationProcessor(8, b)
}

func BenchmarkModificationProcessorWithPoolSize9(b *testing.B) {
	benchmarkModificationProcessor(9, b)
}
