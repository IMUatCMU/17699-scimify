package processor

import "github.com/go-scim/scimify/persistence"

type dbDeleteProcessor struct {
	repo persistence.Repository
}

func (ddp *dbDeleteProcessor) Process(ctx *ProcessorContext) error {
	return ddp.repo.Delete(ctx.Identity)
}
