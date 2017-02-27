package worker

import (
	"context"
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/resource"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func BenchmarkUpdateValueDefaulterWorker(b *testing.B) {
	worker := GetUpdateValueDefaulterWorker()

	ref, _, err := helper.LoadResource("../test_data/single_test_user_david.json")
	if err != nil {
		b.Fatal(err)
	}

	schema, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, resource.CK_Schema, schema)
	ctx = context.WithValue(ctx, resource.CK_Reference, ref)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		r, _, err := helper.LoadResource("../test_data/update_test_user_david.json")
		if err != nil {
			b.Fatal(err)
		}
		for pb.Next() {
			worker.Do(&ValueDefaulterInput{Context: ctx, Resource: r})
		}
	})
}

func TestUpdateValueDefaulterWorker(t *testing.T) {
	worker := GetUpdateValueDefaulterWorker()

	r, _, err := helper.LoadResource("../test_data/update_test_user_david.json")
	if err != nil {
		t.Fatal(err)
	}

	ref, _, err := helper.LoadResource("../test_data/single_test_user_david.json")
	if err != nil {
		t.Fatal(err)
	}

	schema, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, resource.CK_Schema, schema)
	ctx = context.WithValue(ctx, resource.CK_Reference, ref)

	_, err = worker.Do(&ValueDefaulterInput{Context: ctx, Resource: r})
	assert.Nil(t, err)

	assert.True(t, reflect.DeepEqual(r.Attributes["schemas"], []interface{}{"urn:ietf:params:scim:schemas:core:2.0:User"}))
	assert.True(t, reflect.DeepEqual(r.Attributes["id"], "6B69753B-4E38-444E-8AC6-9D0E4D644D80"))
	assert.True(t, reflect.DeepEqual(r.Attributes["externalId"], "changed"))
	assert.True(t, reflect.DeepEqual(r.Attributes["userName"], "changed"))
	assert.True(t, reflect.DeepEqual(r.Attributes["name"], map[string]interface{}{
		"formatted":       "changed",
		"familyName":      "changed",
		"givenName":       "changed",
		"middleName":      "changed",
		"honorificPrefix": "changed",
		"honorificSuffix": "changed",
	}))
	assert.True(t, reflect.DeepEqual(r.Attributes["emails"], []interface{}{
		map[string]interface{}{
			"value":   "david@example.com",
			"type":    "changed",
			"primary": true,
		},
		map[string]interface{}{
			"value": "david@home.com",
			"type":  "changed",
		},
	}))
	assert.NotNil(t, r.Attributes["meta"])
}
