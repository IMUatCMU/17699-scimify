package persistence

import (
	"github.com/go-scim/scimify/resource"
	"sync"
)

var (
	oneSchemaRepo,
	oneResourceTypeRepo,
	oneServiceProviderRepo sync.Once
	schemaRepository,
	resourceTypeRepository,
	serviceProviderRepository Repository
)

func GetSchemaRepository() Repository {
	oneSchemaRepo.Do(func() {
		schemaRepository = &SimpleRepository{
			repo: make(map[string]resource.ScimObject, 0),
		}
	})
	return schemaRepository
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
