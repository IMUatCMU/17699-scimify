package processor

import (
	"encoding/json"
)

type simpleJsonSerializationProcessor struct{}

func (sjp *simpleJsonSerializationProcessor) Process(ctx *ProcessorContext) error {
	target := sjp.getTarget(ctx)
	bytes, err := json.Marshal(target)
	if len(bytes) > 0 {
		ctx.ResponseBody = bytes
	}
	return err
}

func (sjp *simpleJsonSerializationProcessor) getTarget(ctx *ProcessorContext) interface{} {
	target := ctx.SerializationTargetFunc()
	if target == nil {
		panic(&MissingContextValueError{"serialization target"})
	}
	return target
}
