package processor

import "github.com/go-scim/scimify/persistence"

type dbGetProcessor struct {
	repo persistence.Repository
}

func (dgp *dbGetProcessor) Process(ctx *ProcessorContext) error {
	r, err := dgp.repo.Get(ctx.Identity)
	ctx.SingleResult = r
	return err
}
