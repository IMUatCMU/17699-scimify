package worker

import (
	"testing"
	"gopkg.in/mgo.v2/bson"
	"github.com/go-scim/scimify/resource"
	"github.com/stretchr/testify/assert"
	"reflect"
)

var filterWorkerTestMockSchema = &resource.Schema{
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
				FullPath: "userName",
			},
		},
		{
			Name:      "age",
			Type:      resource.Integer,
			CaseExact: false,
			Assist: &resource.Assist{
				JSONName: "age",
				Path:     "age",
				FullPath: "age",
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
						Path:     "city",
						FullPath: "address.city",
					},
				},
			},
			Assist: &resource.Assist{
				JSONName: "address",
				Path:     "address",
				FullPath: "address",
			},
		},
		{
			Name: "active",
			Type: resource.Boolean,
			Assist: &resource.Assist{
				JSONName: "active",
				Path:     "active",
				FullPath: "active",
			},
		},
	},
}

type filterWorkerTest struct {
	name		string
	filter 		string
	assertion	func(result bson.M, err error)
}

func TestFilterWorker(t *testing.T) {
	filterWorkerTestMockSchema.ConstructAttributeIndex()

	for _, test := range []filterWorkerTest{
		{
			"transpile and filter",
			"username eq \"david\" and age gt 17",
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
			func(result bson.M, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			"transpile comparison on complex",
			"address eq \"foo\"",
			func(result bson.M, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			"transpile boolean",
			"active eq false",
			func(result bson.M, err error) {
				assert.Nil(t, err)
				assert.True(t, reflect.DeepEqual(result, bson.M{
					"active": bson.M{
						"$eq": false,
					},
				}))
			},
		},
	}{
		result := FilterWorker(&FilterWorkerInput{
			schema:filterWorkerTestMockSchema,
			filterText:test.filter,
		}).(*WrappedReturn)

		if nil == result.ReturnData {
			test.assertion(nil, result.Err)
		} else {
			test.assertion(result.ReturnData.(bson.M), result.Err)
		}
	}
}
