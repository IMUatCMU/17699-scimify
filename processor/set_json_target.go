package processor

type setSingleResultAsJsonTargetProcessor struct{}

func (_ *setSingleResultAsJsonTargetProcessor) Process(ctx *ProcessorContext) error {
	ctx.SerializationTargetFunc = func() interface{} {
		return ctx.SingleResult
	}
	return nil
}

type setMultipleResultAsJsonTargetProcessor struct{}

func (_ *setMultipleResultAsJsonTargetProcessor) Process(ctx *ProcessorContext) error {
	ctx.SerializationTargetFunc = func() interface{} {
		return ctx.MultiResults
	}
	return nil
}
