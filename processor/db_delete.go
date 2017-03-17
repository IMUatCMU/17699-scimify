package processor

import (
	"github.com/go-scim/scimify/persistence"
	"sync"
)

var (
	oneUserDeleteProcessor,
	oneGroupDeleteProcessor sync.Once

	userDeleteProcessor,
	groupDeleteProcessor Processor
)

func DBUserDeleteProcessor() Processor {
	oneUserDeleteProcessor.Do(func() {
		userDeleteProcessor = &dbDeleteProcessor{
			repo: persistence.GetUserRepository(),
		}
	})
	return userDeleteProcessor
}

func DBGroupDeleteProcessor() Processor {
	oneGroupDeleteProcessor.Do(func() {
		groupDeleteProcessor = &dbDeleteProcessor{
			repo: persistence.GetGroupRepository(),
		}
	})
	return groupDeleteProcessor
}

type dbDeleteProcessor struct {
	repo persistence.Repository
}

func (ddp *dbDeleteProcessor) Process(ctx *ProcessorContext) error {
	return ddp.repo.Delete(ctx.Identity)
}
