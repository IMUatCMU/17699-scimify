package processor

import "github.com/go-scim/scimify/persistence"

type dbCreateProcessor struct {
	repo	persistence.Repository
}

func (dcp *dbCreateProcessor) Process(ctx *ProcessorContext) error {
	r := getResource(ctx, true)
	return dcp.repo.Create(r)
}