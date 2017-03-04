package worker

import (
	"github.com/go-scim/scimify/defaults"
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

// Filter Worker =======================================================================================================

var (
	oneFilterWorker      sync.Once
	filterWorkerInstance *filterWorker
)

func GetFilterWorker() Worker {
	oneFilterWorker.Do(func() {
		filterWorkerInstance = &filterWorker{}
		filterWorkerInstance.initialize(2)
	})
	return filterWorkerInstance
}

// Repository Worker ===================================================================================================

var (
	oneRepoUserQueryWorker,
	oneRepoGroupQueryWorker,
	oneRepoUserCreateWorker,
	oneRepoGroupCreateWorker sync.Once

	repoUserQueryWorkerInstance,
	repoGroupQueryWorkerInstance *repoQueryWorker

	repoUserCreateWorkerInstance,
	repoGroupCreateWorkerInstance *repoCreateWorker
)

func GetRepoUserCreateWorker() Worker {
	oneRepoUserCreateWorker.Do(func() {
		repoUserCreateWorkerInstance = &repoCreateWorker{
			Repo: persistence.NewMongoRepository(
				viper.GetString("mongo.address"),
				viper.GetString("mongo.database"),
				viper.GetString("mongo.userCollectionName")),
		}
		repoUserCreateWorkerInstance.initialize(9)
	})
	return repoUserCreateWorkerInstance
}

func GetRepoGroupCreateWorker() Worker {
	oneRepoGroupCreateWorker.Do(func() {
		repoGroupCreateWorkerInstance = &repoCreateWorker{
			Repo: persistence.NewMongoRepository(
				viper.GetString("mongo.address"),
				viper.GetString("mongo.database"),
				viper.GetString("mongo.groupCollectionName")),
		}
	})
	return repoGroupCreateWorkerInstance
}

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

// JSON Serializer Worker ==============================================================================================

var (
	oneDefaultJsonSerializer, oneSchemaAssistedJsonSerializer           sync.Once
	defaultJsonSerializerInstance, schemaAssistedJsonSerializerInstance *jsonWorker
)

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
			Serializer: &serialize.DynamicJsonSerializer{},
		}
		schemaAssistedJsonSerializerInstance.initialize(1)
	})
	return schemaAssistedJsonSerializerInstance
}

// Validation Worker ===================================================================================================

var (
	oneCreationValidator,
	oneUpdateValidator sync.Once

	creationValidatorInstance,
	updateValidatorInstance *validateWorker
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

// Value Defaulter Worker ==============================================================================================

var (
	oneCreateValueDefaulter,
	oneUpdateValueDefaulter,
	oneSharedValueDefaulter sync.Once

	creationValueDefaulter,
	updateValueDefaulter,
	sharedValueDefaulter *valueDefaulterWorker
)

func GetCreationValueDefaulterWorker() Worker {
	oneCreateValueDefaulter.Do(func() {
		creationValueDefaulter = &valueDefaulterWorker{
			Worker: defaults.GetResourceCreationValueDefaulter(),
		}
		creationValueDefaulter.initialize(2)
	})
	return creationValueDefaulter
}

func GetUpdateValueDefaulterWorker() Worker {
	oneUpdateValueDefaulter.Do(func() {
		updateValueDefaulter = &valueDefaulterWorker{
			Worker: defaults.GetResourceUpdateValueDefaulter(),
		}
		updateValueDefaulter.initialize(2)
	})
	return updateValueDefaulter
}

func GetSharedValueDefaulterWorker() Worker {
	oneSharedValueDefaulter.Do(func() {
		sharedValueDefaulter = &valueDefaulterWorker{
			Worker: defaults.GetSharedValueDefaulter(),
		}
		sharedValueDefaulter.initialize(2)
	})
	return sharedValueDefaulter
}
