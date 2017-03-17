package processor

import (
	"fmt"
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/resource"
	"sync"
)

var (
	oneUserGetToSingleResult,
	oneGroupGetToSingleResult,
	oneUserGetToReference,
	oneGroupGetToReference sync.Once

	userGetToSingleResultProcessor,
	groupGetToSingleResultProcessor,
	userGetToReferenceProcessor,
	groupGetToReferenceProcessor Processor
)

func DBUserGetToSingleResultProcessor() Processor {
	oneUserGetToSingleResult.Do(func() {
		userGetToSingleResultProcessor = &dbGetProcessor{
			repo: persistence.GetUserRepository(),
			f:    putToSingleResult,
		}
	})
	return userGetToSingleResultProcessor
}

func DBGroupGetToSingleResultProcessor() Processor {
	oneGroupGetToSingleResult.Do(func() {
		groupGetToSingleResultProcessor = &dbGetProcessor{
			repo: persistence.GetGroupRepository(),
			f:    putToSingleResult,
		}
	})
	return groupGetToSingleResultProcessor
}

func DBUserGetToReferenceProcessor() Processor {
	oneUserGetToReference.Do(func() {
		userGetToReferenceProcessor = &dbGetProcessor{
			repo: persistence.GetUserRepository(),
			f:    putToReference,
		}
	})
	return userGetToReferenceProcessor
}

func DBGroupGetToReferenceProcessor() Processor {
	oneGroupGetToReference.Do(func() {
		groupGetToReferenceProcessor = &dbGetProcessor{
			repo: persistence.GetGroupRepository(),
			f:    putToReference,
		}
	})
	return groupGetToReferenceProcessor
}

type SingleResultCallback func(obj resource.ScimObject, ctx *ProcessorContext)

var putToSingleResult = func(obj resource.ScimObject, ctx *ProcessorContext) {
	ctx.SingleResult = obj
}

var putToReference = func(obj resource.ScimObject, ctx *ProcessorContext) {
	if r, ok := obj.(*resource.Resource); ok {
		ctx.Reference = r
	}
}

type dbGetProcessor struct {
	repo persistence.Repository
	f    SingleResultCallback
}

func (dgp *dbGetProcessor) Process(ctx *ProcessorContext) error {
	r, err := dgp.repo.Get(ctx.Identity)
	if err != nil {
		return err
	} else if r == nil {
		return resource.CreateError(resource.NotFound, fmt.Sprintf("resource by id %s is not found", ctx.Identity))
	}

	dgp.f(r, ctx)
	return nil
}
