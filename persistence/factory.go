package persistence

import (
	"github.com/go-scim/scimify/resource"
	"github.com/spf13/viper"
	"sync"
)

var (
	oneSchemaRepo,
	oneInternalSchemaRepo,
	oneResourceTypeRepo,
	oneServiceProviderRepo,
	oneUserRepo,
	oneGroupRepo,
	oneRootQueryRepo sync.Once

	schemaRepository,
	internalSchemaRepository,
	resourceTypeRepository,
	serviceProviderRepository,
	userRepository,
	groupRepository,
	rootQueryRepository Repository
)

func GetRootQueryRepository() Repository {
	oneRootQueryRepo.Do(func() {
		rootQueryRepository = &MongoRootQueryRepository{
			address:      viper.GetString("mongo.address"),
			databaseName: viper.GetString("mongo.database"),
			collectionNames: []string{
				viper.GetString("mongo.userCollectionName"),
				viper.GetString("mongo.groupCollectionName"),
			},
		}
	})
	return rootQueryRepository
}

func GetUserRepository() Repository {
	oneUserRepo.Do(func() {
		userRepository = NewMongoRepository(
			viper.GetString("mongo.address"),
			viper.GetString("mongo.database"),
			viper.GetString("mongo.userCollectionName"),
		)
	})
	return userRepository
}

func GetGroupRepository() Repository {
	oneGroupRepo.Do(func() {
		groupRepository = NewMongoRepository(
			viper.GetString("mongo.address"),
			viper.GetString("mongo.database"),
			viper.GetString("mongo.groupCollectionName"),
		)
	})
	return groupRepository
}

func GetSchemaRepository() Repository {
	oneSchemaRepo.Do(func() {
		schemaRepository = &SimpleRepository{
			repo: make(map[string]resource.ScimObject, 0),
		}
	})
	return schemaRepository
}

func GetInternalSchemaRepository() Repository {
	oneInternalSchemaRepo.Do(func() {
		internalSchemaRepository = &SimpleRepository{
			repo: make(map[string]resource.ScimObject, 0),
		}
	})
	return internalSchemaRepository
}

func GetResourceTypeRepository() Repository {
	oneResourceTypeRepo.Do(func() {
		resourceTypeRepository = &SimpleRepository{
			repo: make(map[string]resource.ScimObject, 0),
		}
	})
	return resourceTypeRepository
}

func GetServiceProviderConfigRepository() Repository {
	oneServiceProviderRepo.Do(func() {
		serviceProviderRepository = &SimpleRepository{
			repo: make(map[string]resource.ScimObject, 0),
		}
	})
	return serviceProviderRepository
}
