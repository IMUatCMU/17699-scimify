package worker

import (
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/processor"
	"github.com/go-scim/scimify/serialize"
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
	oneRepoGroupCreateWorker,
	oneRepoUserGetWorker,
	oneRepoGroupGetWorker,
	oneRepoUserDeleteWorker,
	oneRepoGroupDeleteWorker sync.Once

	repoUserQueryWorkerInstance,
	repoGroupQueryWorkerInstance *repoQueryWorker

	repoUserCreateWorkerInstance,
	repoGroupCreateWorkerInstance *repoCreateWorker

	repoUserGetWorkerInstance,
	repoGroupGetWorkerInstance *repoGetWorker

	repoUserDeleteWorkerInstance,
	repoGroupDeleteWorkerInstance *repoDeleteWorker
)

func GetRepoUserDeleteWorker() Worker {
	oneRepoUserDeleteWorker.Do(func() {
		repoUserDeleteWorkerInstance = &repoDeleteWorker{
			Repo: persistence.NewMongoRepository(
				viper.GetString("mongo.address"),
				viper.GetString("mongo.database"),
				viper.GetString("mongo.userCollectionName")),
		}
		repoUserDeleteWorkerInstance.initialize(2)
	})
	return repoUserDeleteWorkerInstance
}

func GetRepoGroupDeleteWorker() Worker {
	oneRepoGroupDeleteWorker.Do(func() {
		repoGroupDeleteWorkerInstance = &repoDeleteWorker{
			Repo: persistence.NewMongoRepository(
				viper.GetString("mongo.address"),
				viper.GetString("mongo.database"),
				viper.GetString("mongo.groupCollectionName")),
		}
		repoGroupDeleteWorkerInstance.initialize(2)
	})
	return repoGroupDeleteWorkerInstance
}

func GetRepoUserGetWorker() Worker {
	oneRepoUserGetWorker.Do(func() {
		repoUserGetWorkerInstance = &repoGetWorker{
			Repo: persistence.NewMongoRepository(
				viper.GetString("mongo.address"),
				viper.GetString("mongo.database"),
				viper.GetString("mongo.userCollectionName")),
		}
		repoUserGetWorkerInstance.initialize(2)
	})
	return repoUserGetWorkerInstance
}

func GetRepoGroupGetWorker() Worker {
	oneRepoGroupGetWorker.Do(func() {
		repoGroupGetWorkerInstance = &repoGetWorker{
			Repo: persistence.NewMongoRepository(
				viper.GetString("mongo.address"),
				viper.GetString("mongo.database"),
				viper.GetString("mongo.groupCollectionName")),
		}
		repoGroupGetWorkerInstance.initialize(2)
	})
	return repoGroupGetWorkerInstance
}

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

// Processor Worker ====================================================================================================

var (
	oneUserCreationProcessor,
	oneUserUpdateProcessor,
	oneGroupCreationProcessor,
	oneGroupUpdateProcessor sync.Once

	userCreationProcessorWorkerInstance,
	userUpdateProcessorWorkerInstance,
	groupCreationProcessorWorkerInstance,
	groupUpdateProcessorWorkerInstance *processorWorker
)

func GetUserCreationProcessorWorker() Worker {
	oneUserCreationProcessor.Do(func() {
		userCreationProcessorWorkerInstance = &processorWorker{
			P: processor.NewSerialProcessor(
				processor.GetFormatCaseInstance(),
				processor.GetValidateTypeProcessor(),
				processor.GetValidateRequiredProcessor(),
				processor.GetGenerateIdInstance(),
				processor.GetGenerateUserMetaInstance(),
			),
		}
		userCreationProcessorWorkerInstance.initialize(2)
	})
	return userCreationProcessorWorkerInstance
}

func GetUserUpdateProcessorWorker() Worker {
	oneUserUpdateProcessor.Do(func() {
		userUpdateProcessorWorkerInstance = &processorWorker{
			P: processor.NewSerialProcessor(
				processor.GetFormatCaseInstance(),
				processor.GetValidateTypeProcessor(),
				processor.GetValidateRequiredProcessor(),
				processor.GetValidateMutabilityInstance(),
				processor.GetUpdateMetaInstance(),
			),
		}
		userUpdateProcessorWorkerInstance.initialize(2)
	})
	return userUpdateProcessorWorkerInstance
}

func GetGroupCreationProcessorWorker() Worker {
	oneGroupCreationProcessor.Do(func() {
		groupCreationProcessorWorkerInstance = &processorWorker{
			P: processor.NewSerialProcessor(
				processor.GetFormatCaseInstance(),
				processor.GetValidateTypeProcessor(),
				processor.GetValidateRequiredProcessor(),
				processor.GetGenerateIdInstance(),
				processor.GetGenerateGroupMetaInstance(),
			),
		}
	})
	return groupCreationProcessorWorkerInstance
}

func GetGroupUpdateProcessorWorker() Worker {
	oneGroupUpdateProcessor.Do(func() {
		groupUpdateProcessorWorkerInstance = &processorWorker{
			P: processor.NewSerialProcessor(
				processor.GetFormatCaseInstance(),
				processor.GetValidateTypeProcessor(),
				processor.GetValidateRequiredProcessor(),
				processor.GetValidateMutabilityInstance(),
				processor.GetUpdateMetaInstance(),
			),
		}
	})
	return groupUpdateProcessorWorkerInstance
}
