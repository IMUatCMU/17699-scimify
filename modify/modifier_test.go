package modify

import (
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/resource"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultModifier_Modify(t *testing.T) {
	modifier := &DefaultModifier{}

	sch, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		t.Fatal(err)
	}
	sch.ConstructAttributeIndex()

	for _, test := range []struct {
		name      string
		mod       *Modification
		assertion func(*resource.Resource, error)
	}{
		{
			"add simple entry",
			&Modification{
				Operations: []ModUnit{
					{
						Op:    "add",
						Path:  "userName",
						Value: "davidqiu",
					},
				},
			},
			func(r *resource.Resource, err error) {
				assert.Nil(t, err)
				assert.Equal(t, "davidqiu", r.Data()["userName"])
			},
		},
		{
			"add nested entry",
			&Modification{
				Operations: []ModUnit{
					{
						Op:    "add",
						Path:  "name.formatted",
						Value: "davidqiu",
					},
				},
			},
			func(r *resource.Resource, err error) {
				assert.Nil(t, err)
				assert.Equal(t, "davidqiu", r.Data()["name"].(map[string]interface{})["formatted"])
			},
		},
		{
			"replace nested entry",
			&Modification{
				Operations: []ModUnit{
					{
						Op:    "replace",
						Path:  "name.formatted",
						Value: "davidqiu",
					},
				},
			},
			func(r *resource.Resource, err error) {
				assert.Nil(t, err)
				assert.Equal(t, "davidqiu", r.Data()["name"].(map[string]interface{})["formatted"])
			},
		},
		{
			"replace all array elements",
			&Modification{
				Operations: []ModUnit{
					{
						Op:    "replace",
						Path:  "emails.value",
						Value: "foo@bar.com",
					},
				},
			},
			func(r *resource.Resource, err error) {
				assert.Nil(t, err)
				for _, elem := range r.Data()["emails"].([]interface{}) {
					assert.Equal(t, "foo@bar.com", elem.(map[string]interface{})["value"])
				}
			},
		},
		{
			"remove all array elements",
			&Modification{
				Operations: []ModUnit{
					{
						Op:   "remove",
						Path: "emails.value",
					},
				},
			},
			func(r *resource.Resource, err error) {
				assert.Nil(t, err)
				for _, elem := range r.Data()["emails"].([]interface{}) {
					assert.Nil(t, elem.(map[string]interface{})["value"])
				}
			},
		},
		{
			"remove some array elements",
			&Modification{
				Operations: []ModUnit{
					{
						Op:   "remove",
						Path: "emails[type eq \"work\"]",
					},
				},
			},
			func(r *resource.Resource, err error) {
				assert.Nil(t, err)
				for _, elem := range r.Data()["emails"].([]interface{}) {
					assert.NotEqual(t, "work", elem.(map[string]interface{})["type"])
				}
			},
		},
		{
			"remove contents in some array",
			&Modification{
				Operations: []ModUnit{
					{
						Op:   "remove",
						Path: "emails[type eq \"work\"].value",
					},
				},
			},
			func(r *resource.Resource, err error) {
				assert.Nil(t, err)
				for _, elem := range r.Data()["emails"].([]interface{}) {
					if elem.(map[string]interface{})["type"] == "work" {
						assert.Nil(t, elem.(map[string]interface{})["value"])
					}
				}
			},
		},
	} {
		r, _, err := helper.LoadResource("../test_data/single_test_user_david.json")
		if err != nil {
			t.Fatal(err)
		}

		err = modifier.Modify(r, sch, test.mod)
		test.assertion(r, err)
	}
}
