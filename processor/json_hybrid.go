package processor

type IntermediateJsonProcessFunc func([]byte, *ProcessorContext) interface{}

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

	ctx.Results[hjp.sjp.argSlot] = hjp.f(ctx.Results[RBodyBytes].([]byte), ctx)

	return hjp.sjp.Process(ctx)
}
