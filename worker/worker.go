package worker

import (
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/serialize"
	"github.com/go-scim/scimify/validation"
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

var oneFilterWorker sync.Once
var filterWorkerInstance *filterWorker

func GetFilterWorker() Worker {
	oneFilterWorker.Do(func() {
		filterWorkerInstance = &filterWorker{}
		filterWorkerInstance.initialize(2)
	})
	return filterWorkerInstance
}

var (
	oneRepoUserQueryWorker, oneRepoGroupQueryWorker sync.Once
)
var repoUserQueryWorkerInstance, repoGroupQueryWorkerInstance *repoQueryWorker

func GetRepoUserQueryWorker() Worker {
	oneRepoUserQueryWorker.Do(func() {
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
	oneRepoGroupQueryWorker.Do(func() {
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

var (
	oneDefaultJsonSerializer, oneSchemaAssistedJsonSerializer sync.Once
)
var defaultJsonSerializerInstance, schemaAssistedJsonSerializerInstance *jsonWorker

func GetDefaultJsonSerializerWorker() Worker {
	oneDefaultJsonSerializer.Do(func() {
		defaultJsonSerializerInstance = &jsonWorker{
			Serializer: &serialize.DefaultJsonSerializer{},
		}
		defaultJsonSerializerInstance.initialize(2)
	})
	return defaultJsonSerializerInstance
}

func GetSchemaAssistedJsonSerializerWorker() Worker {
	oneSchemaAssistedJsonSerializer.Do(func() {
		schemaAssistedJsonSerializerInstance = &jsonWorker{
			Serializer: &serialize.SchemaJsonSerializer{},
		}
		schemaAssistedJsonSerializerInstance.initialize(1)
	})
	return schemaAssistedJsonSerializerInstance
}

var (
	oneCreationValidator, oneUpdateValidator           sync.Once
	creationValidatorInstance, updateValidatorInstance *validateWorker
)

func GetCreationValidatorWorker() Worker {
	oneCreationValidator.Do(func() {
		creationValidatorInstance = &validateWorker{
			Validator: validation.GetResourceCreationValidator(),
		}
		creationValidatorInstance.initialize(9)
	})
	return creationValidatorInstance
}

func GetUpdateValidatorWorker() Worker {
	oneUpdateValidator.Do(func() {
		updateValidatorInstance = &validateWorker{
			Validator: validation.GetResourceUpdateValidator(),
		}
		updateValidatorInstance.initialize(1)
	})
	return updateValidatorInstance
}
