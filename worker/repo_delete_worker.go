package worker

import (
	"context"
	"github.com/go-scim/scimify/persistence"
	"github.com/jeffail/tunny"
)

type RepoDeleteWorkerInput struct {
	Id  string
	Ctx context.Context
}

type repoDeleteWorker struct {
	pool *tunny.WorkPool
	Repo persistence.Repository
}

func (w *repoDeleteWorker) initialize(numProcs int) {
	if pool, err := tunny.CreatePool(numProcs, func(input interface{}) interface{} {
		args := input.(*RepoDeleteWorkerInput)
		r := &wrappedReturn{}
		r.Err = w.Repo.Delete(args.Id, args.Ctx)
		if r.Err != nil {
			r.ReturnData = false
		} else {
			r.ReturnData = true
		}
		return r
	}).Open(); err != nil {
		panic("Failed to initialize repository delete worker pool")
	} else {
		w.pool = pool
	}
}

func (w *repoDeleteWorker) Do(job interface{}) (interface{}, error) {
	if r, err := w.pool.SendWork(job); err != nil {
		return nil, err
	} else {
		return r.(*wrappedReturn).ReturnData, r.(*wrappedReturn).Err
	}
}

func (w *repoDeleteWorker) Close() {
	w.pool.Close()
}
