package processor

import (
	"github.com/jeffail/tunny"
)

type Worker interface {
	Processor
	initialize(int)
	Close()
}

type WorkerWrapper struct {
	processor Processor
	pool      *tunny.WorkPool
}

func (ww *WorkerWrapper) initialize(num int) {
	if pool, err := tunny.CreatePool(num, func(arg interface{}) interface{} {
		return ww.processor.Process(arg.(*ProcessorContext))
	}).Open(); err != nil {
		panic(err)
	} else {
		ww.pool = pool
	}
}

func (ww *WorkerWrapper) Process(ctx *ProcessorContext) error {
	e0, e1 := ww.pool.SendWork(ctx)
	if e0 != nil {
		return e0.(error)
	} else if e1 != nil {
		return e1
	} else {
		return nil
	}
}

func (ww *WorkerWrapper) Close() {
	ww.pool.Close()
}
