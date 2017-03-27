package processor

import (
	"fmt"
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/persistence"
	"github.com/spf13/viper"
	"net/http"
	"testing"
)

func benchmarkParseParamForDeleteEndpointProcessor(poolSize int, b *testing.B) {
	viper.Set("scim.internalSchemaId.user", "urn:ietf:params:scim:schemas:core:2.0:User")
	viper.Set("scim.api.userIdUrlParam", "userId")
	viper.Set(fmt.Sprintf("scim.threadPool.%s", ParamUserDelete), poolSize)

	internalSchemaRepo := persistence.GetInternalSchemaRepository()
	sch, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		b.Fatal(err)
	}
	err = internalSchemaRepo.Create(sch)
	if err != nil {
		b.Fatal(err)
	}

	processor := GetWorkerBean(ParamUserDelete)
	makeContext := func() *ProcessorContext {
		return &ProcessorContext{
			Request: &MockRequestSource{
				M:  http.MethodDelete,
				T:  "/Users/61C2FFDB-C967-4FF0-BF5F-B9A91E550323",
				UP: map[string]string{"userId": "61C2FFDB-C967-4FF0-BF5F-B9A91E550323"},
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

func BenchmarkParseParamForDeleteEndpointProcessorWithPoolSize1(b *testing.B) {
	benchmarkParseParamForDeleteEndpointProcessor(1, b)
}

func BenchmarkParseParamForDeleteEndpointProcessorWithPoolSize2(b *testing.B) {
	benchmarkParseParamForDeleteEndpointProcessor(2, b)
}

func BenchmarkParseParamForDeleteEndpointProcessorWithPoolSize3(b *testing.B) {
	benchmarkParseParamForDeleteEndpointProcessor(3, b)
}

func BenchmarkParseParamForDeleteEndpointProcessorWithPoolSize4(b *testing.B) {
	benchmarkParseParamForDeleteEndpointProcessor(4, b)
}

func BenchmarkParseParamForDeleteEndpointProcessorWithPoolSize5(b *testing.B) {
	benchmarkParseParamForDeleteEndpointProcessor(5, b)
}

func BenchmarkParseParamForDeleteEndpointProcessorWithPoolSize6(b *testing.B) {
	benchmarkParseParamForDeleteEndpointProcessor(6, b)
}

func BenchmarkParseParamForDeleteEndpointProcessorWithPoolSize7(b *testing.B) {
	benchmarkParseParamForDeleteEndpointProcessor(7, b)
}

func BenchmarkParseParamForDeleteEndpointProcessorWithPoolSize8(b *testing.B) {
	benchmarkParseParamForDeleteEndpointProcessor(8, b)
}

func BenchmarkParseParamForDeleteEndpointProcessorWithPoolSize9(b *testing.B) {
	benchmarkParseParamForDeleteEndpointProcessor(9, b)
}
