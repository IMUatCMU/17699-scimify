package processor

import (
	"encoding/json"
	"fmt"
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/resource"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestHybridJsonSerializationProcessor_Process(t *testing.T) {
	processor := &hybridJsonSerializationProcessor{
		sjp: &simpleJsonSerializationProcessor{},
		ajp: &assistedJsonSerializationProcessor{},
		f: func(bytes []byte, _ *ProcessorContext) interface{} {
			raw := json.RawMessage(bytes)
			return resource.NewListResponse(&raw, 0, 10, 2)
		},
	}

	schema, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		t.Fatal(err)
	}
	schema.ConstructAttributeIndex()

	david, json, err := helper.LoadResource("../test_data/single_test_user_david.json")
	if err != nil {
		t.Fatal(err)
	}
	json = strings.Replace(json, "\"password\": \"t1meMa$heen\",", "", 1)

	ctx := &ProcessorContext{
		Schema: schema,
		SerializationTargetFunc:func() interface{} {
			return []resource.ScimObject{david, david}
		},
	}

	err = processor.Process(ctx)
	assert.Nil(t, err)
	assert.JSONEq(t, fmt.Sprintf(`{
		"schemas":["urn:ietf:params:scim:api:messages:2.0:ListResponse"],
		"totalResults":2,
		"itemsPerPage":10,
		"startIndex":0,
		"Resources":[%s, %s]
	}`, json, json), string(ctx.ResponseBody))
}
