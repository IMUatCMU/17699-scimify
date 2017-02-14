package worker

import (
	"github.com/go-scim/scimify/persistence"
	"github.com/spf13/viper"
	"sync"
)

// Common interface for a worker that can do some work
type Worker interface {

	// initialize the thread pool (internally called)
	initialize(numProcs int)

	// Perform work
	Do(interface{}) (interface{}, error)

	// Destroy the thread pool
	Close()
}

type wrappedReturn struct {
	ReturnData interface{}
	Err        error
}

var once sync.Once

var filterWorkerInstance *filterWorker

func GetFilterWorker() Worker {
	once.Do(func() {
		filterWorkerInstance = &filterWorker{}
		filterWorkerInstance.initialize(2)
	})
	return filterWorkerInstance
}

var repoUserQueryWorkerInstance, repoGroupQueryWorkerInstance *repoQueryWorker

func GetRepoUserQueryWorker() Worker {
	once.Do(func() {
		repoUserQueryWorkerInstance = &repoQueryWorker{
			Repo: persistence.NewMongoRepository(
				viper.GetString("mongo.address"),
				viper.GetString("mongo.database"),
				viper.GetString("mongo.userCollectionName")),
		}
		repoUserQueryWorkerInstance.initialize(2)
	})
	return repoUserQueryWorkerInstance
}

func GetRepoGroupQueryWorker() Worker {
	once.Do(func() {
		repoGroupQueryWorkerInstance = &repoQueryWorker{
			Repo: persistence.NewMongoRepository(
				viper.GetString("mongo.address"),
				viper.GetString("mongo.database"),
				viper.GetString("mongo.groupCollectionName")),
		}
		repoGroupQueryWorkerInstance.initialize(2)
	})
	return repoGroupQueryWorkerInstance
}
