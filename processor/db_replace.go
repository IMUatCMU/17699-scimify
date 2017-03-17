package processor

import (
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/resource"
	"sync"
)

var (
	oneUserReplace,
	oneGroupReplace sync.Once

	userReplaceProcessor,
	groupReplaceProcessor Processor
)

func DBUserReplaceProcessor() Processor {
	oneUserReplace.Do(func() {
		userReplaceProcessor = &dbReplaceProcessor{
			repo: persistence.GetUserRepository(),
		}
	})
	return userReplaceProcessor
}

func DBGroupReplaceProcessor() Processor {
	oneGroupReplace.Do(func() {
		groupReplaceProcessor = &dbReplaceProcessor{
			repo: persistence.GetGroupRepository(),
		}
	})
	return groupReplaceProcessor
}

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
