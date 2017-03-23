package processor

import (
	"github.com/go-scim/scimify/modify"
	"github.com/go-scim/scimify/resource"
	"sync"
)

var (
	oneMod      sync.Once
	modInstance Processor
)

func ModificationProcessor() Processor {
	oneMod.Do(func() {
		modInstance = &modProcessor{
			delegate: &modify.DefaultModifier{},
		}
	})
	return modInstance
}

type modProcessor struct {
	delegate modify.Modifier
}

func (mp *modProcessor) Process(ctx *ProcessorContext) error {
	r := mp.getResource(ctx)
	sch := mp.getSchema(ctx)
	mod := mp.getModification(ctx)

	return mp.delegate.Modify(r, sch, mod)
}

func (mp *modProcessor) getResource(ctx *ProcessorContext) *resource.Resource {
	if ctx.Resource == nil {
		panic(&MissingContextValueError{"resource"})
	}
	return ctx.Resource
}

func (mp *modProcessor) getSchema(ctx *ProcessorContext) *resource.Schema {
	if ctx.Schema == nil {
		panic(&MissingContextValueError{"schema"})
	}
	return ctx.Schema
}

func (mp *modProcessor) getModification(ctx *ProcessorContext) *modify.Modification {
	if ctx.Mod == nil {
		panic(&MissingContextValueError{"mod"})
	}
	return ctx.Mod
}
