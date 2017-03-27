package processor

import (
	"fmt"
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/persistence"
	"github.com/spf13/viper"
	"net/http"
	"testing"
)

func benchmarkParseParamForReplaceEndpointProcessor(poolSize int, b *testing.B) {
	viper.Set("scim.internalSchemaId.user", "urn:ietf:params:scim:schemas:core:2.0:User")
	viper.Set("scim.api.userIdUrlParam", "userId")
	viper.Set(fmt.Sprintf("scim.threadPool.%s", ParamUserReplace), poolSize)

	internalSchemaRepo := persistence.GetInternalSchemaRepository()
	sch, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		b.Fatal(err)
	}
	err = internalSchemaRepo.Create(sch)
	if err != nil {
		b.Fatal(err)
	}

	processor := GetWorkerBean(ParamUserReplace)
	makeContext := func() *ProcessorContext {
		return &ProcessorContext{
			Request: &MockRequestSource{
				M:  http.MethodPut,
				T:  "/Users/F0254C21-8034-447E-8DD1-69D7232D67E5",
				UP: map[string]string{"userId": "F0254C21-8034-447E-8DD1-69D7232D67E5"},
				B:  []byte(`{"schemas": ["urn:ietf:params:scim:api:messages:2.0:PatchOp"], "userName": "david2"}`),
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

func BenchmarkParseParamForReplaceEndpointProcessorWithPoolSize1(b *testing.B) {
	benchmarkParseParamForReplaceEndpointProcessor(1, b)
}

func BenchmarkParseParamForReplaceEndpointProcessorWithPoolSize2(b *testing.B) {
	benchmarkParseParamForReplaceEndpointProcessor(2, b)
}

func BenchmarkParseParamForReplaceEndpointProcessorWithPoolSize3(b *testing.B) {
	benchmarkParseParamForReplaceEndpointProcessor(3, b)
}

func BenchmarkParseParamForReplaceEndpointProcessorWithPoolSize4(b *testing.B) {
	benchmarkParseParamForReplaceEndpointProcessor(4, b)
}

func BenchmarkParseParamForReplaceEndpointProcessorWithPoolSize5(b *testing.B) {
	benchmarkParseParamForReplaceEndpointProcessor(5, b)
}

func BenchmarkParseParamForReplaceEndpointProcessorWithPoolSize6(b *testing.B) {
	benchmarkParseParamForReplaceEndpointProcessor(6, b)
}

func BenchmarkParseParamForReplaceEndpointProcessorWithPoolSize7(b *testing.B) {
	benchmarkParseParamForReplaceEndpointProcessor(7, b)
}

func BenchmarkParseParamForReplaceEndpointProcessorWithPoolSize8(b *testing.B) {
	benchmarkParseParamForReplaceEndpointProcessor(8, b)
}

func BenchmarkParseParamForReplaceEndpointProcessorWithPoolSize9(b *testing.B) {
	benchmarkParseParamForReplaceEndpointProcessor(9, b)
}
