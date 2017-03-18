package processor

import (
	"fmt"
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/resource"
	"sync"
)

var (
	oneGetAllResourceType,
	oneGetAllSchemas sync.Once

	getAllResourceType,
	getAllSchemas Processor
)

func DbGetAllResourceTypesProcessor() Processor {
	oneGetAllResourceType.Do(func() {
		getAllResourceType = &dbGetAllProcessor{
			repo: persistence.GetResourceTypeRepository(),
		}
	})
	return getAllResourceType
}

func DbGetAllSchemasProcessor() Processor {
	oneGetAllSchemas.Do(func() {
		getAllSchemas = &dbGetAllProcessor{
			repo: persistence.GetSchemaRepository(),
		}
	})
	return getAllSchemas
}

type dbGetAllProcessor struct {
	repo persistence.Repository
}

func (gap *dbGetAllProcessor) Process(ctx *ProcessorContext) error {
	all, err := gap.repo.GetAll()
	if err != nil {
		return err
	} else if len(all) == 0 {
		return resource.CreateError(resource.NotFound, fmt.Sprintf("resource by id %s is not found", ctx.Identity))
	}

	ctx.MultiResults = all
	return nil
}
