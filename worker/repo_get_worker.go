package worker

import (
	"context"
	"github.com/go-scim/scimify/persistence"
	"github.com/jeffail/tunny"
)

type RepoGetWorkerInput struct {
	Id  string
	Ctx context.Context
}

type repoGetWorker struct {
	pool *tunny.WorkPool
	Repo persistence.Repository
}

func (w *repoGetWorker) initialize(numProcs int) {
	if pool, err := tunny.CreatePool(numProcs, func(input interface{}) interface{} {
		args := input.(*RepoGetWorkerInput)
		r := &wrappedReturn{}
		r.ReturnData, r.Err = w.Repo.Get(args.Id, args.Ctx)
		return r
	}).Open(); err != nil {
		panic("Failed to initialize repository get worker pool")
	} else {
		w.pool = pool
	}
}

func (w *repoGetWorker) Do(job interface{}) (interface{}, error) {
	if r, err := w.pool.SendWork(job); err != nil {
		return nil, err
	} else {
		return r.(*wrappedReturn).ReturnData, r.(*wrappedReturn).Err
	}
}

func (w *repoGetWorker) Close() {
	w.pool.Close()
}
