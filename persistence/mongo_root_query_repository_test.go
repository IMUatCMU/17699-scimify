package persistence

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"testing"
)

func TestMongoRootQueryRepository_Query(t *testing.T) {
	viper.Set("mongo.address", "localhost:32768")
	viper.Set("mongo.database", "test_db")
	viper.Set("mongo.userCollectionName", "users")
	viper.Set("mongo.groupCollectionName", "groups")

	filter := bson.M{"id": bson.M{"$ne": ""}}

	for _, test := range []struct {
		pageStart        int
		pageSize         int
		expectUserCount  int
		expectGroupCount int
	}{
		{0, 2, 2, 0},
		{3, 5, 3, 2},
		{3, 10, 3, 5},
		{6, 3, 0, 3},
		{6, 5, 0, 5},
		{6, 10, 0, 5},
	} {
		repo, cleanUp := prepareTestConnection(t)
		resources, err := repo.Query(filter, "id", true, test.pageStart, test.pageSize)

		assert.Nil(t, err)
		assert.Equal(t, test.expectUserCount+test.expectGroupCount, len(resources))

		userCount, groupCount := 0, 0
		t.Log("--------")
		for _, r := range resources {
			id := r.GetId()
			t.Log(id)
			if strings.HasPrefix(id, "u-") {
				userCount++
			} else if strings.HasPrefix(id, "g-") {
				groupCount++
			}
		}
		t.Log("--------")

		assert.Equal(t, test.expectUserCount, userCount)
		assert.Equal(t, test.expectGroupCount, groupCount)

		cleanUp()
	}
}

func prepareTestConnection(t *testing.T) (*MongoRootQueryRepository, func()) {
	userTestData := []map[string]interface{}{
		{"id": "u-5514A474-3A68-45AB-B8CC-6EFBE4BA2601"},
		{"id": "u-5BFBEEAF-7604-4892-AF4C-51E02545B255"},
		{"id": "u-DBC7D326-AC86-4B2A-BD1C-D7B39B62B7A6"},
		{"id": "u-2EEFAA08-71FC-4799-A621-0C7D8EC49D8B"},
		{"id": "u-36E8832A-B34B-4458-A043-55C1AE639E27"},
	}
	groupTestData := []map[string]interface{}{
		{"id": "g-D615D937-9CF2-4D02-BBF7-307C4C2488DA"},
		{"id": "g-C07FD999-7830-4A66-B250-5B34BE751D97"},
		{"id": "g-37F18850-79EF-4CDD-9C3D-0F793F3BDC8F"},
		{"id": "g-90E41E2A-62D3-41DF-A23B-02825D2DF6C5"},
		{"id": "g-B05D3875-AB79-45FE-B5CF-ECF841F5FE69"},
	}
	repo := &MongoRootQueryRepository{
		address:      viper.GetString("mongo.address"),
		databaseName: viper.GetString("mongo.database"),
		collectionNames: []string{
			viper.GetString("mongo.userCollectionName"),
			viper.GetString("mongo.groupCollectionName"),
		},
	}

	userC := repo.getCollection(repo.getSession(), viper.GetString("mongo.userCollectionName"))
	groupC := repo.getCollection(repo.getSession(), viper.GetString("mongo.groupCollectionName"))
	cleanUp := func() {
		userC.RemoveAll(nil)
		groupC.RemoveAll(nil)
	}

	cleanUp()
	for _, d := range userTestData {
		err := userC.Insert(d)
		if err != nil {
			t.Fatal(err)
		}
	}
	for _, d := range groupTestData {
		err := groupC.Insert(d)
		if err != nil {
			t.Fatal(err)
		}
	}

	return repo, cleanUp
}
