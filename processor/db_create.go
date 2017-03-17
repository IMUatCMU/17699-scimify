package processor

import (
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/resource"
)

type dbCreateProcessor struct {
	repo persistence.Repository
}

func (dcp *dbCreateProcessor) Process(ctx *ProcessorContext) error {
	r := dcp.getResource(ctx)
	return dcp.repo.Create(r)
}

func (dcp *dbCreateProcessor) getResource(ctx *ProcessorContext) *resource.Resource {
	if ctx.Resource == nil {
		panic(&MissingContextValueError{"resource"})
	}
	return ctx.Resource
}
