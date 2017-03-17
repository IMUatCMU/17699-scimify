package processor

import (
	"encoding/json"
	"github.com/go-scim/scimify/resource"
	"sync"
)

var (
	oneListResponseJson       sync.Once
	listResponseJsonProcessor Processor
)

func ListResponseJsonSerializationProcessor() Processor {
	oneListResponseJson.Do(func() {
		listResponseJsonProcessor = &hybridJsonSerializationProcessor{
			sjp: SimpleJsonSerializationProcessor().(*simpleJsonSerializationProcessor),
			ajp: AssistedJsonSerializationProcessor().(*assistedJsonSerializationProcessor),
			f:   listResponsePacker,
		}
	})
	return listResponseJsonProcessor
}

type IntermediateJsonProcessFunc func([]byte, *ProcessorContext) interface{}

var listResponsePacker = func(bytes []byte, ctx *ProcessorContext) interface{} {
	raw := json.RawMessage(bytes)
	return resource.NewListResponse(&raw, ctx.QueryPageStart, ctx.QueryPageSize, len(ctx.MultiResults))
}

type hybridJsonSerializationProcessor struct {
	sjp *simpleJsonSerializationProcessor
	ajp *assistedJsonSerializationProcessor
	f   IntermediateJsonProcessFunc
}

func (hjp *hybridJsonSerializationProcessor) Process(ctx *ProcessorContext) error {
	err := hjp.ajp.Process(ctx)
	if nil != err {
		return err
	}

	intermediate := hjp.f(ctx.ResponseBody, ctx)
	ctx.SerializationTargetFunc = func() interface{} {
		return intermediate
	}

	return hjp.sjp.Process(ctx)
}
