package processor

import (
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/resource"
)

type SingleResultCallback func(obj resource.ScimObject, ctx *ProcessorContext)

type dbGetProcessor struct {
	repo persistence.Repository
	f    SingleResultCallback
}

func (dgp *dbGetProcessor) Process(ctx *ProcessorContext) error {
	r, err := dgp.repo.Get(ctx.Identity)
	if err != nil {
		return err
	}

	dgp.f(r, ctx)
	return nil
}
