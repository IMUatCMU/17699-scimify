package worker

import (
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/resource"
	"github.com/jeffail/tunny"
)

type RepoCreateWorkerInput struct {
	Resource *resource.Resource
	Context  resource.Context
}

type repoCreateWorker struct {
	pool *tunny.WorkPool
	Repo persistence.Repository
}

func (w *repoCreateWorker) initialize(numProcs int) {
	if pool, err := tunny.CreatePool(numProcs, func(input interface{}) interface{} {
		args := input.(*RepoCreateWorkerInput)
		r := &wrappedReturn{}
		r.Err = w.Repo.Create(args.Resource, args.Context)
		if nil == r.Err {
			r.ReturnData = input.(*RepoCreateWorkerInput).Resource
		}
		return r
	}).Open(); err != nil {
		panic("Failed to initialize repository create worker pool")
	} else {
		w.pool = pool
	}
}

func (w *repoCreateWorker) Do(job interface{}) (interface{}, error) {
	if r, err := w.pool.SendWork(job); err != nil {
		return nil, err
	} else {
		return r.(*wrappedReturn).ReturnData, r.(*wrappedReturn).Err
	}
}

func (w *repoCreateWorker) Close() {
	w.pool.Close()
}
