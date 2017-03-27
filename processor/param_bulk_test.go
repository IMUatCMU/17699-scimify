package processor

import (
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"testing"
)

func benchmarkParseParamForBulkEndpointProcessor(poolSize int, b *testing.B) {
	viper.Set("scim.internalSchemaId.user", "urn:ietf:params:scim:schemas:core:2.0:User")
	viper.Set(fmt.Sprintf("scim.threadPool.%s", ParamBulk), poolSize)

	processor := GetWorkerBean(ParamBulk)
	makeContext := func() *ProcessorContext {
		return &ProcessorContext{
			Request: &MockRequestSource{
				M: http.MethodPost,
				T: "/Bulk",
				B: []byte(`{"schemas": ["urn:ietf:params:scim:api:messages:2.0:BulkRequest"], "Operations": [{
					"method": "POST",
					"path": "/Users",
					"data": {
						"schemas": ["urn:ietf:params:scim:api:messages:2.0:User"],
						"userName": "david"
					}
				}, {
					"method": "PUT",
					"path": "/Users/0283874D-A169-4847-A71C-44BF3E4E4242",
					"data": {
						"schemas": ["urn:ietf:params:scim:api:messages:2.0:User"],
						"userName": "anne"
					}
				}]}`),
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

func BenchmarkParseParamForBulkEndpointProcessorWithPoolSize1(b *testing.B) {
	benchmarkParseParamForBulkEndpointProcessor(1, b)
}

func BenchmarkParseParamForBulkEndpointProcessorWithPoolSize2(b *testing.B) {
	benchmarkParseParamForBulkEndpointProcessor(2, b)
}

func BenchmarkParseParamForBulkEndpointProcessorWithPoolSize3(b *testing.B) {
	benchmarkParseParamForBulkEndpointProcessor(3, b)
}

func BenchmarkParseParamForBulkEndpointProcessorWithPoolSize4(b *testing.B) {
	benchmarkParseParamForBulkEndpointProcessor(4, b)
}

func BenchmarkParseParamForBulkEndpointProcessorWithPoolSize5(b *testing.B) {
	benchmarkParseParamForBulkEndpointProcessor(5, b)
}

func BenchmarkParseParamForBulkEndpointProcessorWithPoolSize6(b *testing.B) {
	benchmarkParseParamForBulkEndpointProcessor(6, b)
}

func BenchmarkParseParamForBulkEndpointProcessorWithPoolSize7(b *testing.B) {
	benchmarkParseParamForBulkEndpointProcessor(7, b)
}

func BenchmarkParseParamForBulkEndpointProcessorWithPoolSize8(b *testing.B) {
	benchmarkParseParamForBulkEndpointProcessor(8, b)
}

func BenchmarkParseParamForBulkEndpointProcessorWithPoolSize9(b *testing.B) {
	benchmarkParseParamForBulkEndpointProcessor(9, b)
}
