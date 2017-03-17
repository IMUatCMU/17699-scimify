package processor

import "sync"

var (
	oneSingleResult,
	oneMultipleResult sync.Once

	singleResultInstance,
	multipleResultInstance Processor
)

func SingleResultAsJsonTargetProcessor() Processor {
	oneSingleResult.Do(func() {
		singleResultInstance = &setSingleResultAsJsonTargetProcessor{}
	})
	return singleResultInstance
}

func MultipleResultAsJsonTargetProcessor() Processor {
	oneMultipleResult.Do(func() {
		multipleResultInstance = &setMultipleResultAsJsonTargetProcessor{}
	})
	return multipleResultInstance
}

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
