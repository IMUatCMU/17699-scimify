package worker

import (
	"github.com/go-scim/scimify/filter"
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/resource"
	"github.com/jeffail/tunny"
	"gopkg.in/mgo.v2/bson"
)

type FilterWorkerInput struct {
	filterText string
	schema     *resource.Schema
}

type filterWorker struct {
	pool *tunny.WorkPool
}

func (w *filterWorker) initialize(numProcs int) {
	if pool, err := tunny.CreatePool(numProcs, func(input interface{}) interface{} {
		r := &wrappedReturn{}
		if tokens, err := filter.Tokenize(input.(*FilterWorkerInput).filterText); err != nil {
			r.Err = err
			return r
		} else if root, err := filter.Parse(tokens); err != nil {
			r.Err = err
			return r
		} else if bson, err := persistence.TranspileToMongoQuery(root, input.(*FilterWorkerInput).schema); err != nil {
			r.Err = err
			return r
		} else {
			r.ReturnData = bson
			return r
		}
	}).Open(); err != nil {
		panic("Failed to initialize filter worker pool")
	} else {
		w.pool = pool
	}
}

func (w *filterWorker) Do(job interface{}) (interface{}, error) {
	if r, err := w.pool.SendWork(job); err != nil {
		return bson.M{}, err
	} else if r.(*wrappedReturn).Err != nil {
		return bson.M{}, r.(*wrappedReturn).Err
	} else {
		return r.(*wrappedReturn).ReturnData, nil
	}
}

func (w *filterWorker) Close() {
	w.pool.Close()
}
