package processor

import (
	"github.com/go-scim/scimify/resource"
	"sync"
)

var (
	oneSetAllHeader sync.Once
	setAllHeader    Processor
)

func SetAllHeaderProcessor() Processor {
	oneSetAllHeader.Do(func() {
		setAllHeader = &setHeaderProcessor{eTag: true, location: true}
	})
	return setAllHeader
}

type setHeaderProcessor struct {
	eTag     bool
	location bool
}

func (shp *setHeaderProcessor) Process(ctx *ProcessorContext) error {
	ctx.ResponseHeaders = map[string]string{}
	r := shp.getResource(ctx)

	if shp.eTag {
		ctx.ResponseHeaders["ETag"] = r.Attributes["meta"].(map[string]interface{})["version"].(string)
	}

	if shp.location {
		ctx.ResponseHeaders["Location"] = r.Attributes["meta"].(map[string]interface{})["location"].(string)
	}

	return nil
}

func (shp *setHeaderProcessor) getResource(ctx *ProcessorContext) *resource.Resource {
	if nil == ctx.SingleResult {
		panic(&MissingContextValueError{"resource"})
	}
	return ctx.SingleResult.(*resource.Resource)
}
