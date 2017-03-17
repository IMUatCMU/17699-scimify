package processor

import "github.com/go-scim/scimify/persistence"

type dbQueryProcessor struct {
	repo persistence.Repository
}

func (dqp *dbQueryProcessor) Process(ctx *ProcessorContext) error {
	filter := dqp.getFilter(ctx)

	all, err := dqp.repo.Query(filter, ctx.QuerySortBy, ctx.QuerySortOrder, ctx.QueryPageStart, ctx.QueryPageSize)

	ctx.MultiResults = all

	return err
}

func (dqp *dbQueryProcessor) getFilter(ctx *ProcessorContext) interface{} {
	if ctx.ParsedFilter == nil {
		panic(&MissingContextValueError{"parsed filter"})
	}
	return ctx.ParsedFilter
}
