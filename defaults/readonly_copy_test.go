package defaults

import (
	"context"
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/resource"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestCopyReadOnlyValueDefaulter_Default(t *testing.T) {
	schema, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		t.Fatal(err)
	}

	ref := resource.NewResourceFromMap(map[string]interface{}{
		"id": "EA49451B-38DD-4F62-915F-CCA9E49A0332",
		"meta": map[string]interface{}{
			"resourceType": "User",
			"created":      "2016-01-23T04:56:22Z",
			"lastModified": "2016-05-13T04:42:34Z",
			"version":      "1",
			"location":     "https://example.com/v2/Users/EA49451B-38DD-4F62-915F-CCA9E49A0332",
		},
		"userName": "david@example.com",
		"name": map[string]interface{}{
			"formatted":       "David Qiu",
			"familyName":      "Qiu",
			"givenName":       "David",
			"middleName":      "",
			"honorificPrefix": "Mr.",
			"honorificSuffix": "",
		},
		"emails": []interface{}{
			map[string]interface{}{
				"value":   "david@example.com",
				"type":    "work",
				"primary": true,
			},
			map[string]interface{}{
				"value": "david@home.com",
				"type":  "home",
			},
		},
	})

	r := resource.NewResourceFromMap(map[string]interface{}{
		"id": "EA49451B-38DD-4F62-915F-CCA9E49A0332",
		"meta": map[string]interface{}{
			"resourceType": "changed",
			"created":      "changed",
			"lastModified": "changed",
			"version":      "changed",
			"location":     "changed",
		},
		"userName": "changed",
		"name": map[string]interface{}{
			"formatted":       "changed",
			"familyName":      "changed",
			"givenName":       "changed",
			"middleName":      "changed",
			"honorificPrefix": "changed",
			"honorificSuffix": "changed",
		},
		"emails": []interface{}{
			map[string]interface{}{
				"value":   "david@example.com",
				"type":    "changed",
				"primary": true,
			},
			map[string]interface{}{
				"value": "david@home.com",
				"type":  "changed",
			},
			map[string]interface{}{
				"value": "david@bar.com",
				"type":  "work",
			},
		},
	})

	ctx := context.Background()
	ctx = context.WithValue(ctx, resource.CK_Schema, schema)
	ctx = context.WithValue(ctx, resource.CK_Reference, ref)

	defaulter := &copyReadOnlyValueDefaulter{}
	defaulter.Default(r, ctx)

	assert.True(t, reflect.DeepEqual(r.Attributes, map[string]interface{}{
		"id": "EA49451B-38DD-4F62-915F-CCA9E49A0332",
		"meta": map[string]interface{}{
			"resourceType": "User",
			"created":      "2016-01-23T04:56:22Z",
			"lastModified": "2016-05-13T04:42:34Z",
			"version":      "1",
			"location":     "https://example.com/v2/Users/EA49451B-38DD-4F62-915F-CCA9E49A0332",
		},
		"userName": "changed",
		"name": map[string]interface{}{
			"formatted":       "changed",
			"familyName":      "changed",
			"givenName":       "changed",
			"middleName":      "changed",
			"honorificPrefix": "changed",
			"honorificSuffix": "changed",
		},
		"emails": []interface{}{
			map[string]interface{}{
				"value":   "david@example.com",
				"type":    "changed",
				"primary": true,
			},
			map[string]interface{}{
				"value": "david@home.com",
				"type":  "changed",
			},
			map[string]interface{}{
				"value": "david@bar.com",
				"type":  "work",
			},
		},
	}))
}
