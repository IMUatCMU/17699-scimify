package processor

import "github.com/go-scim/scimify/persistence"

type dbQueryProcessor struct {
	repo persistence.Repository
}

func (dqp *dbQueryProcessor) Process(ctx *ProcessorContext) error {
	filter := getA(ctx, ArgFilter, true, nil)
	sortBy := getString(ctx, ArgSortBy, true, "")
	sortOrder := getBool(ctx, ArgSortOrder, true, false)
	pageStart := getInt(ctx, ArgPageStart, true, 0)
	pageSize := getInt(ctx, ArgPageSize, true, 0)

	all, err := dqp.repo.Query(filter, sortBy, sortOrder, pageStart, pageSize)
	ctx.Results[RAllResources] = all
	return err
}
