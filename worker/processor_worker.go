package worker

import (
	"context"
	"github.com/go-scim/scimify/processor"
	"github.com/go-scim/scimify/resource"
	"github.com/jeffail/tunny"
)

type ProcessorInput struct {
	R   *resource.Resource
	Ctx context.Context
}

type processorWorker struct {
	P    processor.Processor
	pool *tunny.WorkPool
}

func (w *processorWorker) initialize(numProcs int) {
	if pool, err := tunny.CreatePool(numProcs, func(input interface{}) interface{} {
		r := &wrappedReturn{}

		if err := w.P.Process(
			input.(*ProcessorInput).R,
			input.(*ProcessorInput).Ctx,
		); err != nil {
			r.Err = err
			r.ReturnData = false
		} else {
			r.Err = nil
			r.ReturnData = true
		}

		return r
	}).Open(); err != nil {
		panic("Failed to initialize processor worker pool")
	} else {
		w.pool = pool
	}
}

func (w *processorWorker) Do(job interface{}) (interface{}, error) {
	if r, err := w.pool.SendWork(job); err != nil {
		return false, err
	} else {
		return r.(*wrappedReturn).ReturnData, r.(*wrappedReturn).Err
	}
}

func (w *processorWorker) Close() {
	w.pool.Close()
}
