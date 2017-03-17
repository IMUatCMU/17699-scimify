package processor

import (
	"github.com/go-scim/scimify/resource"
	"github.com/satori/go.uuid"
	"strings"
	"sync"
)

var (
	oneGenerateId      sync.Once
	generateIdInstance Processor
)

func GenerateIdProcessor() Processor {
	oneGenerateId.Do(func() {
		generateIdInstance = &generateIdProcessor{}
	})
	return generateIdInstance
}

type generateIdProcessor struct{}

func (gip *generateIdProcessor) Process(ctx *ProcessorContext) error {
	r := gip.getResource(ctx)
	r.Attributes["id"] = strings.ToLower(uuid.NewV4().String())
	return nil
}

func (gip *generateIdProcessor) getResource(ctx *ProcessorContext) *resource.Resource {
	if ctx.Resource == nil {
		panic(&MissingContextValueError{"resource"})
	}
	return ctx.Resource
}
