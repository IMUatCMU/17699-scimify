package worker

import (
	"encoding/json"
	"github.com/go-scim/scimify/resource"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type repoQueryWorkerTest struct {
	name      string
	input     *RepoQueryWorkerInput
	assertion func([]*resource.Resource, error)
}

var repoQueryWorkerTestParams = []*RepoQueryWorkerInput{
	{
		filter: bson.M{
			"userName": bson.M{
				"$regex": bson.RegEx{
					Pattern: "^david",
					Options: "i",
				},
			},
		},
		sortBy: "", ascending: true,
		pageStart: 0, pageSize: 0,
		context: nil,
	},
	{
		filter: bson.M{
			"name.familyName": bson.M{
				"$regex": bson.RegEx{
					Pattern: "u",
					Options: "i",
				},
			},
		},
		sortBy: "nickName", ascending: true,
		pageStart: 0, pageSize: 0,
		context: nil,
	},
	{
		filter: bson.M{
			"userName": bson.M{
				"$regex": bson.RegEx{
					Pattern: "^D",
					Options: "i",
				},
			},
		},
		sortBy: "", ascending: true,
		pageStart: 0, pageSize: 0,
		context: nil,
	},
	{
		filter: bson.M{
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
		sortBy: "", ascending: true,
		pageStart: 0, pageSize: 0,
		context: nil,
	},
	{
		filter: bson.M{
			"meta.lastModified": bson.M{
				"$gt": "2016-05-13T04:42:34Z",
			},
		},
		sortBy: "", ascending: true,
		pageStart: 0, pageSize: 0,
		context: nil,
	},
	{
		filter: bson.M{
			"meta.lastModified": bson.M{
				"$gte": "2016-05-13T04:42:34Z",
			},
		},
		sortBy: "nickName", ascending: true,
		pageStart: 0, pageSize: 0,
		context: nil,
	},
	{
		filter: bson.M{
			"meta.lastModified": bson.M{
				"$lt": "2016-05-13T04:42:34Z",
			},
		},
		sortBy: "nickName", ascending: true,
		pageStart: 0, pageSize: 0,
		context: nil,
	},
	{
		filter: bson.M{
			"meta.lastModified": bson.M{
				"$lte": "2016-05-13T04:42:34Z",
			},
		},
		sortBy: "nickName", ascending: true,
		pageStart: 0, pageSize: 0,
		context: nil,
	},
	{
		filter: bson.M{
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
		},
		sortBy: "nickName", ascending: true,
		pageStart: 0, pageSize: 0,
		context: nil,
	},
	{
		filter: bson.M{
			"$or": []interface{}{
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
							Pattern: "^Intern$",
							Options: "i",
						},
					},
				},
			},
		},
		sortBy: "nickName", ascending: true,
		pageStart: 0, pageSize: 0,
		context: nil,
	},
}

func TestRepoQueryWorker(t *testing.T) {
	viper.Set("mongo.address", "localhost:32768")
	viper.Set("mongo.database", "test_db")
	viper.Set("mongo.userCollectionName", "users")

	PrepareTestMongoConnection(t, "../test_data/test_users_2.json")

	worker := GetRepoUserQueryWorker()
	defer worker.Close()
	for _, test := range []repoQueryWorkerTest{
		{
			"title pr or userType eq \"Intern\"",
			repoQueryWorkerTestParams[9],
			func(results []*resource.Resource, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 2, len(results))
				assert.Equal(t, "A", results[0].Attributes["nickName"])
				assert.Equal(t, "Q", results[1].Attributes["nickName"])
			},
		},
		{
			"title pr and userType eq \"Employee\"",
			repoQueryWorkerTestParams[8],
			func(results []*resource.Resource, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 1, len(results))
				assert.Equal(t, "Q", results[0].Attributes["nickName"])
			},
		},
		{
			"meta.lastModified le \"2016-05-13T04:42:34Z\"",
			repoQueryWorkerTestParams[7],
			func(results []*resource.Resource, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 1, len(results))
				assert.Equal(t, "Q", results[0].Attributes["nickName"])
			},
		},
		{
			"meta.lastModified lt \"2016-05-13T04:42:34Z\"",
			repoQueryWorkerTestParams[6],
			func(results []*resource.Resource, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 0, len(results))
			},
		},
		{
			"meta.lastModified ge \"2016-05-13T04:42:34Z\"",
			repoQueryWorkerTestParams[5],
			func(results []*resource.Resource, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 2, len(results))
				assert.Equal(t, "A", results[0].Attributes["nickName"])
				assert.Equal(t, "Q", results[1].Attributes["nickName"])
			},
		},
		{
			"meta.lastModified gt \"2016-05-13T04:42:34Z\"",
			repoQueryWorkerTestParams[4],
			func(results []*resource.Resource, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 1, len(results))
				assert.Equal(t, "A", results[0].Attributes["nickName"])
			},
		},
		{
			"title pr",
			repoQueryWorkerTestParams[3],
			func(results []*resource.Resource, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 1, len(results))
				assert.Equal(t, "Q", results[0].Attributes["nickName"])
			},
		},
		{
			"userName sw \"D\"",
			repoQueryWorkerTestParams[2],
			func(results []*resource.Resource, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 1, len(results))
				assert.Equal(t, "Q", results[0].Attributes["nickName"])
			},
		},
		{
			"name.familyName co \"u\"",
			repoQueryWorkerTestParams[1],
			func(results []*resource.Resource, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 2, len(results))
				assert.Equal(t, "A", results[0].Attributes["nickName"])
				assert.Equal(t, "Q", results[1].Attributes["nickName"])
			},
		},
		{
			"userName eq \"david\"",
			repoQueryWorkerTestParams[0],
			func(results []*resource.Resource, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 1, len(results))
				assert.Equal(t, "Q", results[0].Attributes["nickName"])
			},
		},
	} {
		r, err := worker.Do(test.input)
		test.assertion(r.([]*resource.Resource), err)
	}
}

func BenchmarkRepoQueryWorker(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	viper.Set("mongo.address", "localhost:32768")
	viper.Set("mongo.database", "test_db")
	viper.Set("mongo.userCollectionName", "users")

	PrepareTestMongoConnection(b, "../test_data/test_users_2.json")
	worker := GetRepoUserQueryWorker()
	defer worker.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := r.Intn(len(repoQueryWorkerTestParams))
			b.Logf("%d\n", n)
			worker.Do(repoQueryWorkerTestParams[n])
		}
	})
}

func PrepareTestMongoConnection(t interface{}, dataPath string) {
	testData := make(map[string]interface{})
	fatal := func(err error) {
		if t0, ok := t.(*testing.T); ok {
			t0.Fatal(err)
		} else if b0, ok := t.(*testing.B); ok {
			b0.Fatal(err)
		}
	}

	if path, err := filepath.Abs(dataPath); err != nil {
		fatal(err)
	} else {
		file, err := os.Open(path)
		if err != nil {
			fatal(err)
		}
		defer file.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fatal(err)
		}

		err = json.Unmarshal(fileBytes, &testData)
		if err != nil {
			fatal(err)
		}
	}

	session, err := mgo.Dial(viper.GetString("mongo.address"))
	if err != nil {
		fatal(err)
	}

	collection := session.DB(viper.GetString("mongo.database")).C(viper.GetString("mongo.userCollectionName"))
	collection.RemoveAll(nil)

	for _, each := range testData["data"].([]interface{}) {
		collection.Insert(each)
	}
}
