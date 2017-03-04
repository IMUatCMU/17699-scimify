package persistence

import (
	"github.com/go-scim/scimify/resource"
	"sync"
)

var (
	oneSchemaRepo,
	oneResourceTypeRepo    sync.Once
	schemaRepository,
	resourceTypeRepository Repository
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
