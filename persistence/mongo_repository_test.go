package persistence

import (
	"encoding/json"
	"github.com/go-scim/scimify/resource"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// This test case requires tester to have a local mongoDB instance
// - listening on port 32768
// - a test_db database
// - user collection
// The easiest way to setup is to use docker

func TestMongoRepository_Query(t *testing.T) {
	repo, cleanUp := prepareTestMongoConnection(t)
	defer cleanUp(repo)

	for _, test := range []struct {
		name      string
		filter    bson.M
		sortBy    string
		sortOrder bool
		skip      int
		size      int
		context   resource.Context
		assertion func([]*resource.Resource, error)
	}{
		{
			"simple equality filter",
			bson.M{
				"userName": bson.M{"$eq": "david@example.com"},
			},
			"", false,
			0, 0,
			nil,
			func(resources []*resource.Resource, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 1, len(resources))
				assert.Equal(t, "david@example.com", resources[0].Attributes["userName"])
			},
		},
		{
			"and filter",
			bson.M{
				"$and": []interface{}{
					bson.M{
						"addresses.locality": bson.M{"$eq": "Toronto"},
					},
					bson.M{
						"externalId": bson.M{"$eq": "996624032"},
					},
				},
			},
			"", false,
			0, 0,
			nil,
			func(resources []*resource.Resource, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 1, len(resources))
				assert.Equal(t, "david@example.com", resources[0].Attributes["userName"])
			},
		},
		{
			"test sort",
			bson.M{},
			"nickName", true,
			0, 0,
			nil,
			func(resources []*resource.Resource, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 3, len(resources))
				assert.Equal(t, "Babs", resources[0].Attributes["nickName"])
				assert.Equal(t, "Q", resources[1].Attributes["nickName"])
				assert.Equal(t, "Tom", resources[2].Attributes["nickName"])
			},
		},
		{
			"test sort reversed order",
			bson.M{},
			"nickName", false,
			0, 0,
			nil,
			func(resources []*resource.Resource, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 3, len(resources))
				assert.Equal(t, "Tom", resources[0].Attributes["nickName"])
				assert.Equal(t, "Q", resources[1].Attributes["nickName"])
				assert.Equal(t, "Babs", resources[2].Attributes["nickName"])
			},
		},
		{
			"test paging",
			bson.M{},
			"nickName", true,
			1, 1,
			nil,
			func(resources []*resource.Resource, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 1, len(resources))
				assert.Equal(t, "Q", resources[0].Attributes["nickName"])
			},
		},
	} {
		test.assertion(repo.Query(test.filter, test.sortBy, test.sortOrder, test.skip, test.size, test.context))
	}
}

func prepareTestMongoConnection(t *testing.T) (*MongoRepository, func(*MongoRepository)) {
	cleanUp := func(r *MongoRepository) {
		r.getCollection(r.getSession()).RemoveAll(nil)
	}

	testData := make(map[string]interface{})
	if path, err := filepath.Abs("../test_data/test_users.json"); err != nil {
		t.Fatal(err)
	} else {
		file, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			t.Fatal(err)
		}

		err = json.Unmarshal(fileBytes, &testData)
		if err != nil {
			t.Fatal(err)
		}
	}
	assert.NotEmpty(t, testData["data"], "Test data must not be empty")

	repo := NewMongoRepository("localhost:32768", "test_db", "users")
	c := repo.getCollection(repo.getSession())
	cleanUp(repo)

	for _, each := range testData["data"].([]interface{}) {
		c.Insert(each)
	}

	return repo, cleanUp
}
