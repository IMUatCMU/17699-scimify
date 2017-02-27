package worker

import (
	"context"
	"github.com/go-scim/scimify/defaults"
	"github.com/go-scim/scimify/resource"
	"github.com/jeffail/tunny"
)

type ValueDefaulterInput struct {
	Resource *resource.Resource
	Context  context.Context
}

type valueDefaulterWorker struct {
	Worker defaults.ValueDefaulter
	pool   *tunny.WorkPool
}

func (w *valueDefaulterWorker) initialize(numProcs int) {
	if pool, err := tunny.CreatePool(numProcs, func(input interface{}) interface{} {
		ok, err := w.Worker.Default(
			input.(*ValueDefaulterInput).Resource,
			input.(*ValueDefaulterInput).Context)

		r := &wrappedReturn{}
		r.ReturnData = ok
		r.Err = err

		return r
	}).Open(); err != nil {
		panic("Failed to initialize value defaulter worker pool")
	} else {
		w.pool = pool
	}
}

func (w *valueDefaulterWorker) Do(job interface{}) (interface{}, error) {
	if r, err := w.pool.SendWork(job); err != nil {
		return false, err
	} else {
		return r.(*wrappedReturn).ReturnData, r.(*wrappedReturn).Err
	}
}

func (w *valueDefaulterWorker) Close() {
	w.pool.Close()
}
