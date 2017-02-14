package worker

import "sync"

// Common interface for a worker that can do some work
type Worker interface {

	// initialize the thread pool (internally called)
	initialize(numProcs int)

	// Perform work
	Do(interface{}) (interface{}, error)

	// Destroy the thread pool
	Close()
}

type wrappedReturn struct {
	ReturnData interface{}
	Err        error
}

var once sync.Once

var filterWorkerInstance *filterWorker

func GetFilterWorker() Worker {
	once.Do(func() {
		filterWorkerInstance = &filterWorker{}
		filterWorkerInstance.initialize(2)
	})
	return filterWorkerInstance
}
