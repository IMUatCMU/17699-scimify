package processor

import (
	"github.com/go-scim/scimify/persistence"
	"sync"
)

var (
	oneUserQuery,
	oneGroupQuery,
	oneRootQuery sync.Once

	userQueryProcessor,
	groupQueryProcessor,
	rootQueryProcessor Processor
)

func DBRootQueryProcessor() Processor {
	oneRootQuery.Do(func() {
		rootQueryProcessor = &dbQueryProcessor{
			repo: persistence.GetRootQueryRepository(),
		}
	})
	return rootQueryProcessor
}

func DBUserQueryProcessor() Processor {
	oneUserQuery.Do(func() {
		userQueryProcessor = &dbQueryProcessor{
			repo: persistence.GetUserRepository(),
		}
	})
	return userQueryProcessor
}

func DBGroupQueryProcessor() Processor {
	oneGroupQuery.Do(func() {
		groupQueryProcessor = &dbQueryProcessor{
			repo: persistence.GetGroupRepository(),
		}
	})
	return groupQueryProcessor
}

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
