package processor

import "sync"

var (
	oneSetResource,
	oneSingleResult,
	oneMultipleResult,
	oneSetError sync.Once

	setResourceInstance,
	singleResultInstance,
	multipleResultInstance,
	setErrorInstance Processor
)

func ResourceAsJsonTargetProcessor() Processor {
	oneSetResource.Do(func() {
		setResourceInstance = &setAsJsonTargetProcessor{
			f: func(ctx *ProcessorContext) func() interface{} {
				ctx.SingleResult = ctx.Resource
				return func() interface{} {
					return ctx.SingleResult
				}
			},
		}
	})
	return setResourceInstance
}

func ErrorAsJsonTargetProcessor() Processor {
	oneSetError.Do(func() {
		setErrorInstance = &setAsJsonTargetProcessor{
			f: func(ctx *ProcessorContext) func() interface{} {
				return func() interface{} {
					return ctx.Err
				}
			},
		}
	})
	return setErrorInstance
}

func SingleResultAsJsonTargetProcessor() Processor {
	oneSingleResult.Do(func() {
		singleResultInstance = &setAsJsonTargetProcessor{
			f: func(ctx *ProcessorContext) func() interface{} {
				return func() interface{} {
					return ctx.SingleResult
				}
			},
		}
	})
	return singleResultInstance
}

func MultipleResultAsJsonTargetProcessor() Processor {
	oneMultipleResult.Do(func() {
		multipleResultInstance = &setAsJsonTargetProcessor{
			f: func(ctx *ProcessorContext) func() interface{} {
				return func() interface{} {
					return ctx.MultiResults
				}
			},
		}
	})
	return multipleResultInstance
}

type setAsJsonTargetProcessor struct {
	f func(ctx *ProcessorContext) func() interface{}
}

func (jtp *setAsJsonTargetProcessor) Process(ctx *ProcessorContext) error {
	ctx.SerializationTargetFunc = jtp.f(ctx)
	return nil
}
