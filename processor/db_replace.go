package processor

import (
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/resource"
)

type dbReplaceProcessor struct {
	repo persistence.Repository
}

func (drp *dbReplaceProcessor) Process(ctx *ProcessorContext) error {
	r := drp.getResource(ctx)
	return drp.repo.Replace(ctx.Identity, r)
}

func (drp *dbReplaceProcessor) getResource(ctx *ProcessorContext) *resource.Resource {
	if ctx.Resource == nil {
		panic(&MissingContextValueError{"resource"})
	}
	return ctx.Resource
}
