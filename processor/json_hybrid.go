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

	intermediate := hjp.f(ctx.ResponseBody, ctx)
	ctx.SerializationTargetFunc = func() interface{} {
		return intermediate
	}

	return hjp.sjp.Process(ctx)
}
