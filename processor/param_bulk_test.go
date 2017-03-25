package processor

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"strings"
	"github.com/stretchr/testify/assert"
	"github.com/spf13/viper"
)

var (
	testBulkBodyOk = `{
		"schemas": ["urn:ietf:params:scim:api:messages:2.0:BulkRequest"],
		"Operations": [
			{
				"method": "post",
				"bulkId": "qwerty",
				"path": "/Users",
				"data": {
					"schemas": ["urn:ietf:params:scim:api:messages:2.0:User"],
					"userName": "david"
				}
			},
			{
				"method": "put",
				"bulkId": "ytrewq",
				"path": "/Users/7867B2D1-D1B3-4CA7-A51B-D3B405D70B6B",
				"data": {
					"schemas": ["urn:ietf:params:scim:api:messages:2.0:User"],
					"userName": "anne"
				}
			}
		]
	}`
	testBulkInvalidMethod = `{
		"schemas": ["urn:ietf:params:scim:api:messages:2.0:BulkRequest"],
		"Operations": [
			{
				"method": "options",
				"bulkId": "qwerty",
				"path": "/Users",
				"data": {
					"schemas": ["urn:ietf:params:scim:api:messages:2.0:User"],
					"userName": "david"
				}
			}
		]
	}`
	testBulkInvalidPath = `{
		"schemas": ["urn:ietf:params:scim:api:messages:2.0:BulkRequest"],
		"Operations": [
			{
				"method": "post",
				"bulkId": "qwerty",
				"path": "/Foo",
				"data": {
					"schemas": ["urn:ietf:params:scim:api:messages:2.0:User"],
					"userName": "david"
				}
			}
		]
	}`
)

func TestParseParamForBulkEndpointProcessor_Process(t *testing.T) {
	viper.Set("scim.resourceTypeUri.user", "/Users")
	viper.Set("scim.resourceTypeUri.group", "/Groups")

	p := &parseParamForBulkEndpointProcessor{}
	for _, test := range []struct{
		name 		string
		req 		*http.Request
		assertion 	func(bulk *BulkRequest, err error)
	}{
		{
			"not serializable body",
			httptest.NewRequest(http.MethodPost, "/Bulk", strings.NewReader(`scribble`)),
			func(bulk *BulkRequest, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			"normal bulk",
			httptest.NewRequest(http.MethodPost, "/Bulk", strings.NewReader(testBulkBodyOk)),
			func(bulk *BulkRequest, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 2, len(bulk.Operations))
				assert.Equal(t, http.MethodPost, strings.ToUpper(bulk.Operations[0].Method))
				assert.NotEmpty(t, bulk.Operations[0].Data)
				assert.Equal(t, http.MethodPut, strings.ToUpper(bulk.Operations[1].Method))
				assert.NotEmpty(t, bulk.Operations[1].Data)
			},
		},
		{
			"invalid http method",
			httptest.NewRequest(http.MethodPost, "/Bulk", strings.NewReader(testBulkInvalidMethod)),
			func(bulk *BulkRequest, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			"invalid path",
			httptest.NewRequest(http.MethodPost, "/Bulk", strings.NewReader(testBulkInvalidPath)),
			func(bulk *BulkRequest, err error) {
				assert.NotNil(t, err)
			},
		},
	}{
		ctx := &ProcessorContext{
			Request:&HttpRequestSource{Req:test.req},
		}
		err := p.Process(ctx)
		test.assertion(ctx.Bulk, err)
	}
}