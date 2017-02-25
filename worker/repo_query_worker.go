package worker

import (
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/resource"
	"github.com/jeffail/tunny"
)

type RepoQueryWorkerInput struct {
	Filter    interface{}
	SortBy    string
	Ascending bool
	PageStart int
	PageSize  int
	Context   resource.Context
}

type repoQueryWorker struct {
	pool *tunny.WorkPool
	Repo persistence.Repository
}

func (w *repoQueryWorker) initialize(numProcs int) {
	if pool, err := tunny.CreatePool(numProcs, func(input interface{}) interface{} {
		args := input.(*RepoQueryWorkerInput)
		r := &wrappedReturn{}

		if results, err := w.Repo.Query(args.Filter,
			args.SortBy, args.Ascending,
			args.PageStart, args.PageSize,
			args.Context); err != nil {
			r.Err = err
			return r
		} else {
			r.ReturnData = results
			return r
		}
	}).Open(); err != nil {
		panic("Failed to initialize repository query worker pool")
	} else {
		w.pool = pool
	}
}

func (w *repoQueryWorker) Do(job interface{}) (interface{}, error) {
	if r, err := w.pool.SendWork(job); err != nil {
		return make([]*resource.Resource, 0), err
	} else if r.(*wrappedReturn).Err != nil {
		return make([]*resource.Resource, 0), r.(*wrappedReturn).Err
	} else {
		return r.(*wrappedReturn).ReturnData, nil
	}
}

func (w *repoQueryWorker) Close() {
	w.pool.Close()
}
