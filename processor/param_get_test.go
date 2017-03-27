package processor

import (
	"fmt"
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/persistence"
	"github.com/spf13/viper"
	"net/http"
	"testing"
)

func benchmarkParseParamForGetEndpointProcessor(poolSize int, b *testing.B) {
	viper.Set("scim.internalSchemaId.user", "urn:ietf:params:scim:schemas:core:2.0:User")
	viper.Set("scim.api.userIdUrlParam", "userId")
	viper.Set(fmt.Sprintf("scim.threadPool.%s", ParamUserGet), poolSize)

	internalSchemaRepo := persistence.GetInternalSchemaRepository()
	sch, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		b.Fatal(err)
	}
	err = internalSchemaRepo.Create(sch)
	if err != nil {
		b.Fatal(err)
	}

	processor := GetWorkerBean(ParamUserGet)
	makeContext := func() *ProcessorContext {
		return &ProcessorContext{
			Request: &MockRequestSource{
				M:  http.MethodGet,
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

func BenchmarkParseParamForGetEndpointProcessorWithPoolSize1(b *testing.B) {
	benchmarkParseParamForGetEndpointProcessor(1, b)
}

func BenchmarkParseParamForGetEndpointProcessorWithPoolSize2(b *testing.B) {
	benchmarkParseParamForGetEndpointProcessor(2, b)
}

func BenchmarkParseParamForGetEndpointProcessorWithPoolSize3(b *testing.B) {
	benchmarkParseParamForGetEndpointProcessor(3, b)
}

func BenchmarkParseParamForGetEndpointProcessorWithPoolSize4(b *testing.B) {
	benchmarkParseParamForGetEndpointProcessor(4, b)
}

func BenchmarkParseParamForGetEndpointProcessorWithPoolSize5(b *testing.B) {
	benchmarkParseParamForGetEndpointProcessor(5, b)
}

func BenchmarkParseParamForGetEndpointProcessorWithPoolSize6(b *testing.B) {
	benchmarkParseParamForGetEndpointProcessor(6, b)
}

func BenchmarkParseParamForGetEndpointProcessorWithPoolSize7(b *testing.B) {
	benchmarkParseParamForGetEndpointProcessor(7, b)
}

func BenchmarkParseParamForGetEndpointProcessorWithPoolSize8(b *testing.B) {
	benchmarkParseParamForGetEndpointProcessor(8, b)
}

func BenchmarkParseParamForGetEndpointProcessorWithPoolSize9(b *testing.B) {
	benchmarkParseParamForGetEndpointProcessor(9, b)
}
