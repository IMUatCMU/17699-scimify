package persistence

import (
	"github.com/go-scim/scimify/resource"
	"sync"
)

var (
	oneSchemaRepo    sync.Once
	schemaRepository Repository
)

func GetSchemaRepository() Repository {
	oneSchemaRepo.Do(func() {
		schemaRepository = &SimpleRepository{
			repo: make(map[string]resource.ScimObject, 0),
		}
	})
	return schemaRepository
}
