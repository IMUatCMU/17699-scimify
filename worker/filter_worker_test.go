package worker

import (
	"github.com/go-scim/scimify/resource"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

type filterWorkerTest struct {
	name      string
	filter    string
	assertion func(result bson.M, err error)
}

func TestFilterWorker(t *testing.T) {
	schema := &resource.Schema{
		Schemas:    []string{resource.SchemaUrn},
		Id:         resource.UserUrn,
		Name:       "User schema",
		Attributes: make([]*resource.Attribute, 0),
	}
	coreSchema, err := loadSchema("../schemas/common_schema.json")
	userSchema, err := loadSchema("../schemas/user_schema.json")
	if err != nil {
		t.Fatal(err)
	}
	schema.MergeWith(coreSchema, userSchema)
	schema.ConstructAttributeIndex()

	for _, test := range []filterWorkerTest{
		{
			"14",
			"emails[type eq \"work\" and value co \"@example.com\"] or ims[type eq \"xmpp\" and value co \"@foo.com\"]",
			func(result bson.M, err error) {
				assert.Nil(t, err)
				assert.True(t, reflect.DeepEqual(result, bson.M{
					"$or": []interface{}{
						bson.M{
							"$and": []interface{}{
								bson.M{
									"emails.type": bson.M{
										"$regex": bson.RegEx{
											Pattern: "^work$",
											Options: "i",
										},
									},
								},
								bson.M{
									"emails.value": bson.M{
										"$regex": bson.RegEx{
											Pattern: "@example.com",
											Options: "i",
										},
									},
								},
							},
						},
						bson.M{
							"$and": []interface{}{
								bson.M{
									"ims.type": bson.M{
										"$regex": bson.RegEx{
											Pattern: "^xmpp$",
											Options: "i",
										},
									},
								},
								bson.M{
									"ims.value": bson.M{
										"$regex": bson.RegEx{
											Pattern: "@foo.com",
											Options: "i",
										},
									},
								},
							},
						},
					},
				}))
			},
		},
		{
			"13",
			"userType eq \"Employee\" and (emails.type eq \"work\")",
			func(result bson.M, err error) {
				assert.Nil(t, err)
				assert.True(t, reflect.DeepEqual(result, bson.M{
					"$and": []interface{}{
						bson.M{
							"userType": bson.M{
								"$regex": bson.RegEx{
									Pattern: "^Employee$",
									Options: "i",
								},
							},
						},
						bson.M{
							"emails.type": bson.M{
								"$regex": bson.RegEx{
									Pattern: "^work$",
									Options: "i",
								},
							},
						},
					},
				}))
			},
		},
		{
			"12",
			"userType eq \"Employee\" and not emails.value co \"example.org\"",
			func(result bson.M, err error) {
				assert.Nil(t, err)
				assert.True(t, reflect.DeepEqual(result, bson.M{
					"$and": []interface{}{
						bson.M{
							"userType": bson.M{
								"$regex": bson.RegEx{
									Pattern: "^Employee$",
									Options: "i",
								},
							},
						},
						bson.M{
							"$nor": []interface{}{
								bson.M{
									"emails.value": bson.M{
										"$regex": bson.RegEx{
											Pattern: "example.org",
											Options: "i",
										},
									},
								},
							},
						},
					},
				}))
			},
		},
		{
			"11",
			"schemas eq \"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User\"",
			func(result bson.M, err error) {
				assert.Nil(t, err)
				assert.True(t, reflect.DeepEqual(result, bson.M{
					"schemas": bson.M{
						"$eq": "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User",
					},
				}))
			},
		},
		{
			"10",
			"title pr and userType eq \"Employee\"",
			func(result bson.M, err error) {
				assert.Nil(t, err)
				assert.True(t, reflect.DeepEqual(result, bson.M{
					"$and": []interface{}{
						bson.M{
							"$and": []interface{}{
								bson.M{
									"title": bson.M{"$exists": true},
								},
								bson.M{
									"title": bson.M{"$ne": nil},
								},
								bson.M{
									"title": bson.M{"$ne": ""},
								},
							},
						},
						bson.M{
							"userType": bson.M{
								"$regex": bson.RegEx{
									Pattern: "^Employee$",
									Options: "i",
								},
							},
						},
					},
				}))
			},
		},
		{
			"9",
			"meta.lastModified le \"2011-05-13T04:42:34Z\"",
			func(result bson.M, err error) {
				assert.Nil(t, err)
				assert.True(t, reflect.DeepEqual(result, bson.M{
					"meta.lastModified": bson.M{
						"$lte": "2011-05-13T04:42:34Z",
					},
				}))
			},
		},
		{
			"8",
			"meta.lastModified lt \"2011-05-13T04:42:34Z\"",
			func(result bson.M, err error) {
				assert.Nil(t, err)
				assert.True(t, reflect.DeepEqual(result, bson.M{
					"meta.lastModified": bson.M{
						"$lt": "2011-05-13T04:42:34Z",
					},
				}))
			},
		},
		{
			"7",
			"meta.lastModified ge \"2011-05-13T04:42:34Z\"",
			func(result bson.M, err error) {
				assert.Nil(t, err)
				assert.True(t, reflect.DeepEqual(result, bson.M{
					"meta.lastModified": bson.M{
						"$gte": "2011-05-13T04:42:34Z",
					},
				}))
			},
		},
		{
			"6",
			"meta.lastModified gt \"2011-05-13T04:42:34Z\"",
			func(result bson.M, err error) {
				assert.Nil(t, err)
				assert.True(t, reflect.DeepEqual(result, bson.M{
					"meta.lastModified": bson.M{
						"$gt": "2011-05-13T04:42:34Z",
					},
				}))
			},
		},
		{
			"5",
			"title pr",
			func(result bson.M, err error) {
				assert.Nil(t, err)
				assert.True(t, reflect.DeepEqual(result, bson.M{
					"$and": []interface{}{
						bson.M{
							"title": bson.M{"$exists": true},
						},
						bson.M{
							"title": bson.M{"$ne": nil},
						},
						bson.M{
							"title": bson.M{"$ne": ""},
						},
					},
				}))
			},
		},
		{
			"4",
			"urn:ietf:params:scim:schemas:core:2.0:User:userName sw \"J\"",
			func(result bson.M, err error) {
				assert.Nil(t, err)
				assert.True(t, reflect.DeepEqual(result, bson.M{
					"userName": bson.M{
						"$regex": bson.RegEx{
							Pattern: "^J",
							Options: "i",
						},
					},
				}))
			},
		},
		{
			"3",
			"userName sw \"J\"",
			func(result bson.M, err error) {
				assert.Nil(t, err)
				assert.True(t, reflect.DeepEqual(result, bson.M{
					"userName": bson.M{
						"$regex": bson.RegEx{
							Pattern: "^J",
							Options: "i",
						},
					},
				}))
			},
		},
		{
			"2",
			"name.familyName co \"O'Malley\"",
			func(result bson.M, err error) {
				assert.Nil(t, err)
				assert.True(t, reflect.DeepEqual(result, bson.M{
					"name.familyName": bson.M{
						"$regex": bson.RegEx{
							Pattern: "O'Malley",
							Options: "i",
						},
					},
				}))
			},
		},
		{
			"1",
			"userName eq \"john\"",
			func(result bson.M, err error) {
				assert.Nil(t, err)
				assert.True(t, reflect.DeepEqual(result, bson.M{
					"userName": bson.M{
						"$regex": bson.RegEx{
							Pattern: "^john$",
							Options: "i",
						},
					},
				}))
			},
		},
	} {
		worker := GetFilterWorker()
		result, err := worker.Do(&FilterWorkerInput{
			schema:     schema,
			filterText: test.filter,
		})
		test.assertion(result.(bson.M), err)
	}
}

func BenchmarkFilterWorker(b *testing.B) {
	testFilters := []string{
		"userName eq \"john\"",
		"name.familyName co \"O'Malley\"",
		"userName sw \"J\"",
		"urn:ietf:params:scim:schemas:core:2.0:User:userName sw \"J\"",
		"title pr",
		"meta.lastModified gt \"2011-05-13T04:42:34Z\"",
		"meta.lastModified ge \"2011-05-13T04:42:34Z\"",
		"meta.lastModified lt \"2011-05-13T04:42:34Z\"",
		"meta.lastModified le \"2011-05-13T04:42:34Z\"",
		"title pr and userType eq \"Employee\"",
		"title pr or userType eq \"Intern\"",
		"schemas eq \"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User\"",
		"userType eq \"Employee\" and (emails co \"example.com\" or emails.value co \"example.org\")",
		"userType ne \"Employee\" and not (emails co \"example.com\" or emails.value co \"example.org\")",
		"userType eq \"Employee\" and (emails.type eq \"work\")",
		"userType eq \"Employee\" and emails[type eq \"work\" and value co \"@example.com\"]",
		"emails[type eq \"work\" and value co \"@example.com\"] or ims[type eq \"xmpp\" and value co \"@foo.com\"]",
	}
	r := rand.New(rand.NewSource(time.Now().Unix()))

	schema := &resource.Schema{
		Schemas:    []string{resource.SchemaUrn},
		Id:         resource.UserUrn,
		Name:       "User schema",
		Attributes: make([]*resource.Attribute, 0),
	}
	coreSchema, err := loadSchema("../schemas/common_schema.json")
	userSchema, err := loadSchema("../schemas/user_schema.json")
	if err != nil {
		b.Fatal(err)
	}
	schema.MergeWith(coreSchema, userSchema)
	schema.ConstructAttributeIndex()

	worker := GetFilterWorker()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			worker.Do(&FilterWorkerInput{
				schema:     schema,
				filterText: testFilters[r.Intn(len(testFilters)-1)],
			})
		}

	})
}

func loadSchema(filePath string) (*resource.Schema, error) {
	if path, err := filepath.Abs(filePath); err != nil {
		return nil, err
	} else if schema, err := resource.LoadSchema(path); err != nil {
		return nil, err
	} else {
		return schema, nil
	}
}
