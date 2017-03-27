package processor

import (
	"fmt"
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/persistence"
	"github.com/spf13/viper"
	"net/http"
	"testing"
)

func benchmarkParseParamForQueryEndpointProcessor(poolSize int, b *testing.B) {
	viper.Set("scim.internalSchemaId.user", "urn:ietf:params:scim:schemas:core:2.0:User")
	viper.Set(fmt.Sprintf("scim.threadPool.%s", ParamUserQuery), poolSize)

	internalSchemaRepo := persistence.GetInternalSchemaRepository()
	sch, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		b.Fatal(err)
	}
	err = internalSchemaRepo.Create(sch)
	if err != nil {
		b.Fatal(err)
	}

	processor := GetWorkerBean(ParamUserQuery)
	makeContext := func() *ProcessorContext {
		return &ProcessorContext{
			Request: &MockRequestSource{
				M: http.MethodGet,
				T: "/Users",
				P: map[string]string{
					"filter": "username eq \"david\" and x509Certificates.value sw \"abc\"",
				},
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

func BenchmarkParseParamForQueryEndpointProcessorWithPoolSize1(b *testing.B) {
	benchmarkParseParamForQueryEndpointProcessor(1, b)
}

func BenchmarkParseParamForQueryEndpointProcessorWithPoolSize2(b *testing.B) {
	benchmarkParseParamForQueryEndpointProcessor(2, b)
}

func BenchmarkParseParamForQueryEndpointProcessorWithPoolSize3(b *testing.B) {
	benchmarkParseParamForQueryEndpointProcessor(3, b)
}

func BenchmarkParseParamForQueryEndpointProcessorWithPoolSize4(b *testing.B) {
	benchmarkParseParamForQueryEndpointProcessor(4, b)
}

func BenchmarkParseParamForQueryEndpointProcessorWithPoolSize5(b *testing.B) {
	benchmarkParseParamForQueryEndpointProcessor(5, b)
}

func BenchmarkParseParamForQueryEndpointProcessorWithPoolSize6(b *testing.B) {
	benchmarkParseParamForQueryEndpointProcessor(6, b)
}

func BenchmarkParseParamForQueryEndpointProcessorWithPoolSize7(b *testing.B) {
	benchmarkParseParamForQueryEndpointProcessor(7, b)
}

func BenchmarkParseParamForQueryEndpointProcessorWithPoolSize8(b *testing.B) {
	benchmarkParseParamForQueryEndpointProcessor(8, b)
}

func BenchmarkParseParamForQueryEndpointProcessorWithPoolSize9(b *testing.B) {
	benchmarkParseParamForQueryEndpointProcessor(9, b)
}
