package worker

import (
	"github.com/go-scim/scimify/helper"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepoCreateWorker(t *testing.T) {
	viper.Set("mongo.address", "localhost:32768")
	viper.Set("mongo.database", "test_db")
	viper.Set("mongo.userCollectionName", "users")
	PrepareTestMongoConnection(t, "")

	resource, _, err := helper.LoadResource("../test_data/single_test_user_david.json")
	if err != nil {
		t.Fatal(err)
	}

	worker := GetRepoUserCreateWorker()
	input := &RepoCreateWorkerInput{Resource: resource, Context: nil}
	_, err = worker.Do(input)
	assert.Nil(t, err)
}

func BenchmarkRepoCreateWorker(b *testing.B) {
	viper.Set("mongo.address", "localhost:32768")
	viper.Set("mongo.database", "test_db")
	viper.Set("mongo.userCollectionName", "users")
	PrepareTestMongoConnection(b, "")

	resource, _, err := helper.LoadResource("../test_data/single_test_user_david.json")
	if err != nil {
		b.Fatal(err)
	}

	worker := GetRepoUserCreateWorker()
	input := &RepoCreateWorkerInput{Resource: resource, Context: nil}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			worker.Do(input)
		}
	})
}
