package processor

import (
	"github.com/go-scim/scimify/filter"
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/resource"
	"sync"
)

var (
	oneFilter   sync.Once
	parseFilter Processor
)

func ParseFilterProcessor() Processor {
	oneFilter.Do(func() {
		parseFilter = &parseFilterProcessor{}
	})
	return parseFilter
}

type parseFilterProcessor struct{}

func (pfp *parseFilterProcessor) Process(ctx *ProcessorContext) error {
	schema := pfp.getSchema(ctx)

	tokens, err := filter.Tokenize(ctx.QueryFilter)
	if err != nil {
		return err
	}

	node, err := filter.Parse(tokens)
	if err != nil {
		return err
	}

	bson, err := persistence.TranspileToMongoQuery(node, schema)
	if err != nil {
		return err
	}

	ctx.ParsedFilter = bson
	return nil
}

func (pfp *parseFilterProcessor) getSchema(ctx *ProcessorContext) *resource.Schema {
	if ctx.Schema == nil {
		panic(&MissingContextValueError{"schema"})
	}
	return ctx.Schema
}
