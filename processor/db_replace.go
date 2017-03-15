package processor

import "github.com/go-scim/scimify/persistence"

type dbReplaceProcessor struct {
	repo 	persistence.Repository
}

func (drp *dbReplaceProcessor) Process(ctx *ProcessorContext) error {
	r := getResource(ctx, true)
	return drp.repo.Replace(ctx.Identity, r)
}