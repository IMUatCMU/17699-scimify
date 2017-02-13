package worker

import (
	"github.com/go-scim/scimify/filter"
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/resource"
)

type FilterWorkerInput struct {
	filterText string
	schema     *resource.Schema
}

func FilterWorker(input interface{}) interface{} {
	result := &WrappedReturn{}

	tokens, err := filter.Tokenize(input.(*FilterWorkerInput).filterText)
	if err != nil {
		result.Err = err
		return result
	}

	root, err := filter.Parse(tokens)
	if err != nil {
		result.Err = err
		return result
	}

	bson, err := persistence.TranspileToMongoQuery(root, input.(*FilterWorkerInput).schema)
	if err != nil {
		result.Err = err
		return result
	}

	result.ReturnData = bson
	return result
}
