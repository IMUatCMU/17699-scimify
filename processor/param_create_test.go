package processor

import (
	"fmt"
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/persistence"
	"github.com/spf13/viper"
	"net/http"
	"testing"
)

func benchmarkParseParamForCreateEndpointProcessor(poolSize int, b *testing.B) {
	viper.Set("scim.internalSchemaId.user", "urn:ietf:params:scim:schemas:core:2.0:User")
	viper.Set(fmt.Sprintf("scim.threadPool.%s", ParamUserCreate), poolSize)

	internalSchemaRepo := persistence.GetInternalSchemaRepository()
	sch, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		b.Fatal(err)
	}
	err = internalSchemaRepo.Create(sch)
	if err != nil {
		b.Fatal(err)
	}

	processor := GetWorkerBean(ParamUserCreate)
	makeContext := func() *ProcessorContext {
		return &ProcessorContext{
			Request: &MockRequestSource{
				M: http.MethodPost,
				T: "/Users",
				B: []byte(`{"schemas": ["urn:ietf:params:scim:schemas:core:2.0:User"], "userName": "david"}`),
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

func BenchmarkParseParamForCreateEndpointProcessorWithPoolSize1(b *testing.B) {
	benchmarkParseParamForCreateEndpointProcessor(1, b)
}

func BenchmarkParseParamForCreateEndpointProcessorWithPoolSize2(b *testing.B) {
	benchmarkParseParamForCreateEndpointProcessor(2, b)
}

func BenchmarkParseParamForCreateEndpointProcessorWithPoolSize3(b *testing.B) {
	benchmarkParseParamForCreateEndpointProcessor(3, b)
}

func BenchmarkParseParamForCreateEndpointProcessorWithPoolSize4(b *testing.B) {
	benchmarkParseParamForCreateEndpointProcessor(4, b)
}

func BenchmarkParseParamForCreateEndpointProcessorWithPoolSize5(b *testing.B) {
	benchmarkParseParamForCreateEndpointProcessor(5, b)
}

func BenchmarkParseParamForCreateEndpointProcessorWithPoolSize6(b *testing.B) {
	benchmarkParseParamForCreateEndpointProcessor(6, b)
}

func BenchmarkParseParamForCreateEndpointProcessorWithPoolSize7(b *testing.B) {
	benchmarkParseParamForCreateEndpointProcessor(7, b)
}

func BenchmarkParseParamForCreateEndpointProcessorWithPoolSize8(b *testing.B) {
	benchmarkParseParamForCreateEndpointProcessor(8, b)
}

func BenchmarkParseParamForCreateEndpointProcessorWithPoolSize9(b *testing.B) {
	benchmarkParseParamForCreateEndpointProcessor(9, b)
}
