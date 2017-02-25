package worker

import (
	"context"
	"github.com/go-scim/scimify/resource"
	"github.com/go-scim/scimify/validation"
	"github.com/jeffail/tunny"
)

type ValidationInput struct {
	Resource *resource.Resource
	Option   validation.ValidationOptions
	Context  context.Context
}

type validateWorker struct {
	Validator validation.Validator
	pool      *tunny.WorkPool
}

func (w *validateWorker) initialize(numProcs int) {
	if pool, err := tunny.CreatePool(numProcs, func(input interface{}) interface{} {
		ok, err := w.Validator.Validate(
			input.(*ValidationInput).Resource,
			input.(*ValidationInput).Option,
			input.(*ValidationInput).Context)

		r := &wrappedReturn{}
		r.ReturnData = ok
		r.Err = err

		return r
	}).Open(); err != nil {
		panic("Failed to initialize validation worker pool")
	} else {
		w.pool = pool
	}
}

func (w *validateWorker) Do(job interface{}) (interface{}, error) {
	if r, err := w.pool.SendWork(job); err != nil {
		return false, err
	} else {
		return r.(*wrappedReturn).ReturnData, r.(*wrappedReturn).Err
	}
}

func (w *validateWorker) Close() {
	w.pool.Close()
}
