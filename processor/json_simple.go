package processor

import (
	"encoding/json"
)

type simpleJsonSerializationProcessor struct {
	argSlot 	RName
}

func (sjp *simpleJsonSerializationProcessor) Process(ctx *ProcessorContext) error {
	target := getR(ctx, sjp.argSlot, true, nil)
	bytes, err := json.Marshal(target)
	if len(bytes) > 0 {
		ctx.Results[RBodyBytes] = bytes
	}
	return err
}