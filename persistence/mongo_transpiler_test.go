package persistence

import (
	fltr "github.com/go-scim/scimify/filter"
	"github.com/go-scim/scimify/resource"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	_ "reflect"
	"testing"
)

type mongoTranspilerTest struct {
	name      string
	filter    string
	schema    *resource.Schema
	assertion func(bson.M, error)
}

var testSchema = &resource.Schema{
	Id:   resource.UserUrn,
	Name: "Test Schema",
	Attributes: []*resource.Attribute{
		{
			Name:      "userName",
			Type:      resource.String,
			CaseExact: false,
			Assist: &resource.Assist{
				JSONName: "userName",
				Path:     "userName",
				FullPath: "urn:ietf:params:scim:schemas:core:2.0:User:userName",
			},
		},
		{
			Name:      "age",
			Type:      resource.Integer,
			CaseExact: false,
			Assist: &resource.Assist{
				JSONName: "age",
				Path:     "age",
				FullPath: "urn:ietf:params:scim:schemas:core:2.0:User:age",
			},
		},
		{
			Name: "address",
			Type: resource.Complex,
			SubAttributes: []*resource.Attribute{
				{
					Name:      "city",
					Type:      resource.String,
					CaseExact: true,
					Assist: &resource.Assist{
						JSONName: "city",
						Path:     "address.city",
						FullPath: "urn:ietf:params:scim:schemas:core:2.0:User:address.city",
					},
				},
			},
			Assist: &resource.Assist{
				JSONName: "address",
				Path:     "address",
				FullPath: "urn:ietf:params:scim:schemas:core:2.0:User:address",
			},
		},
		{
			Name: "active",
			Type: resource.Boolean,
			Assist: &resource.Assist{
				JSONName: "active",
				Path:     "active",
				FullPath: "urn:ietf:params:scim:schemas:core:2.0:User:active",
			},
		},
	},
}

func TestTranspileToMongoQuery(t *testing.T) {
	testSchema.ConstructAttributeIndex()

	for _, test := range []mongoTranspilerTest{
		{
			"transpile and filter",
			"username eq \"david\" and age gt 17",
			testSchema,
			func(result bson.M, err error) {
				assert.Nil(t, err)
				assert.True(t, reflect.DeepEqual(result, bson.M{
					"$and": []interface{}{
						bson.M{
							"userName": bson.M{
								"$regex": bson.RegEx{
									Pattern: "^david$",
									Options: "i",
								},
							},
						},
						bson.M{
							"age": bson.M{
								"$gt": int64(17),
							},
						},
					},
				}))
			},
		},
		{
			"transpile nested filter",
			"address[city sw \"Sh\"]",
			testSchema,
			func(result bson.M, err error) {
				assert.Nil(t, err)
				assert.True(t, reflect.DeepEqual(result, bson.M{
					"address.city": bson.M{
						"$regex": bson.RegEx{
							Pattern: "^Sh",
							Options: "",
						},
					},
				}))
			},
		},
		{
			"transpile string attribute starts with number",
			"username sw 3",
			testSchema,
			func(result bson.M, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			"transpile comparison on complex",
			"address eq \"foo\"",
			testSchema,
			func(result bson.M, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			"transpile boolean",
			"active eq false",
			testSchema,
			func(result bson.M, err error) {
				assert.Nil(t, err)
				assert.True(t, reflect.DeepEqual(result, bson.M{
					"active": bson.M{
						"$eq": false,
					},
				}))
			},
		},
	} {
		if tokens, err := fltr.Tokenize(test.filter); err != nil {
			t.Fatal(err, test.name)
		} else if root, err := fltr.Parse(tokens); err != nil {
			t.Fatal(err, test.name)
		} else {
			test.assertion(TranspileToMongoQuery(root, test.schema))
		}
	}
}
