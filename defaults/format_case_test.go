package defaults

import (
	"context"
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/resource"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestCaseFormatValueDefaulter_Default(t *testing.T) {
	schema, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		t.Fatal(err)
	}

	r := resource.NewResourceFromMap(map[string]interface{}{
		"ID": "35621E54-E96D-41E1-8E6A-4BA787EB071A",
		"MeTA": map[string]interface{}{
			"VersiON":      "12345",
			"resourcetype": "User",
		},
		"Emails": []interface{}{
			map[string]interface{}{
				"VALUE": "foo@bar.com",
			},
			map[string]interface{}{
				"VALUE": "bar@foo.com",
			},
		},
	})

	defaulter := &caseFormatValueDefaulter{}

	ctx := context.WithValue(context.Background(), resource.CK_Schema, schema)
	defaulter.Default(r, ctx)

	assert.True(t, reflect.DeepEqual(r.Attributes, map[string]interface{}{
		"id": "35621E54-E96D-41E1-8E6A-4BA787EB071A",
		"meta": map[string]interface{}{
			"version":      "12345",
			"resourceType": "User",
		},
		"emails": []interface{}{
			map[string]interface{}{
				"value": "foo@bar.com",
			},
			map[string]interface{}{
				"value": "bar@foo.com",
			},
		},
	}))
}
