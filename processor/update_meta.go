package processor

import (
	"github.com/go-scim/scimify/resource"
	"sync"
)

var (
	oneUpdateMeta sync.Once
	updateMeta    Processor
)

func UpdateMetaProcessor() Processor {
	oneUpdateMeta.Do(func() {
		updateMeta = &updateMetaProcessor{}
	})
	return updateMeta
}

type updateMetaProcessor struct{}

func (ump *updateMetaProcessor) Process(ctx *ProcessorContext) error {
	r := ump.getResource(ctx)

	if meta, ok := r.Attributes["meta"].(map[string]interface{}); !ok {
		panic(&PrerequisiteFailedError{reporter: "meta update", requirement: "meta"})
	} else if id, ok := r.Attributes["id"].(string); !ok {
		panic(&PrerequisiteFailedError{reporter: "meta update", requirement: "id"})
	} else {
		newMeta := make(map[string]interface{})
		for k, v := range meta {
			newMeta[k] = v
		}

		newMeta["version"] = generateNewVersion(id)
		newMeta["lastModified"] = getCurrentTime()
		r.Attributes["meta"] = newMeta

		return nil
	}
}

func (ump *updateMetaProcessor) getResource(ctx *ProcessorContext) *resource.Resource {
	if ctx.Resource == nil {
		panic(&MissingContextValueError{"resource"})
	}
	return ctx.Resource
}
