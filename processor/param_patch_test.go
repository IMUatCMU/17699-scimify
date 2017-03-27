package processor

import (
	"fmt"
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/persistence"
	"github.com/spf13/viper"
	"net/http"
	"testing"
)

func benchmarkParseParamForPatchEndpointProcessor(poolSize int, b *testing.B) {
	viper.Set("scim.internalSchemaId.user", "urn:ietf:params:scim:schemas:core:2.0:User")
	viper.Set("scim.api.userIdUrlParam", "userId")
	viper.Set(fmt.Sprintf("scim.threadPool.%s", ParamUserPatch), poolSize)

	internalSchemaRepo := persistence.GetInternalSchemaRepository()
	sch, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		b.Fatal(err)
	}
	err = internalSchemaRepo.Create(sch)
	if err != nil {
		b.Fatal(err)
	}

	processor := GetWorkerBean(ParamUserPatch)
	makeContext := func() *ProcessorContext {
		return &ProcessorContext{
			Request: &MockRequestSource{
				M:  http.MethodPatch,
				T:  "/Users/F0254C21-8034-447E-8DD1-69D7232D67E5",
				UP: map[string]string{"userId": "F0254C21-8034-447E-8DD1-69D7232D67E5"},
				B:  []byte(`{"schemas": ["urn:ietf:params:scim:api:messages:2.0:PatchOp"], "Operations": [{"op": "add", "path": "userName", "value": "david"}]}`),
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

func BenchmarkParseParamForPatchEndpointProcessorWithPoolSize1(b *testing.B) {
	benchmarkParseParamForPatchEndpointProcessor(1, b)
}

func BenchmarkParseParamForPatchEndpointProcessorWithPoolSize2(b *testing.B) {
	benchmarkParseParamForPatchEndpointProcessor(2, b)
}

func BenchmarkParseParamForPatchEndpointProcessorWithPoolSize3(b *testing.B) {
	benchmarkParseParamForPatchEndpointProcessor(3, b)
}

func BenchmarkParseParamForPatchEndpointProcessorWithPoolSize4(b *testing.B) {
	benchmarkParseParamForPatchEndpointProcessor(4, b)
}

func BenchmarkParseParamForPatchEndpointProcessorWithPoolSize5(b *testing.B) {
	benchmarkParseParamForPatchEndpointProcessor(5, b)
}

func BenchmarkParseParamForPatchEndpointProcessorWithPoolSize6(b *testing.B) {
	benchmarkParseParamForPatchEndpointProcessor(6, b)
}

func BenchmarkParseParamForPatchEndpointProcessorWithPoolSize7(b *testing.B) {
	benchmarkParseParamForPatchEndpointProcessor(7, b)
}

func BenchmarkParseParamForPatchEndpointProcessorWithPoolSize8(b *testing.B) {
	benchmarkParseParamForPatchEndpointProcessor(8, b)
}

func BenchmarkParseParamForPatchEndpointProcessorWithPoolSize9(b *testing.B) {
	benchmarkParseParamForPatchEndpointProcessor(9, b)
}
