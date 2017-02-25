package worker

import (
	"context"
	"github.com/go-scim/scimify/resource"
	"github.com/go-scim/scimify/serialize"
	"github.com/jeffail/tunny"
)

type JsonSerializeInput struct {
	Resource       *resource.Resource
	InclusionPaths []string
	ExclusionPaths []string
	Context        context.Context
}

type jsonWorker struct {
	Serializer serialize.JSONSerializer
	pool       *tunny.WorkPool
}

func (w *jsonWorker) initialize(numProcs int) {
	if pool, err := tunny.CreatePool(numProcs, func(input interface{}) interface{} {
		r := &wrappedReturn{}
		if bytes, err := w.Serializer.Serialize(
			input.(*JsonSerializeInput).Resource,
			input.(*JsonSerializeInput).InclusionPaths,
			input.(*JsonSerializeInput).ExclusionPaths,
			input.(*JsonSerializeInput).Context,
		); err != nil {
			r.Err = err
			return r
		} else {
			r.ReturnData = bytes
			return r
		}
	}).Open(); err != nil {
		panic("Failed to initialize json serializer worker pool")
	} else {
		w.pool = pool
	}
}

func (w *jsonWorker) Do(job interface{}) (interface{}, error) {
	if r, err := w.pool.SendWork(job); err != nil {
		return []byte(``), err
	} else if r.(*wrappedReturn).Err != nil {
		return []byte(``), r.(*wrappedReturn).Err
	} else {
		return r.(*wrappedReturn).ReturnData, nil
	}
}

func (w *jsonWorker) Close() {
	w.pool.Close()
}
