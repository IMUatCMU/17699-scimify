package processor

import (
	"github.com/go-scim/scimify/resource"
	"net/http"
)

type ProcessorContext struct {
	Identity	string
	Resource 	*resource.Resource
	Reference 	*resource.Resource
	Schema 		*resource.Schema
	Request 	*http.Request
	MiscArgs	map[AName]interface{}
	Results 	map[RName]interface{}
}

type AName string

const (
	ArgResource = AName("resource")
	ArgReference = AName("reference")
	ArgSchema = AName("schema")
	ArgFilter = AName("filter")
	ArgSortBy = AName("sortBy")
	ArgSortOrder = AName("sortOrder")
	ArgPageStart = AName("pageStart")
	ArgPageSize = AName("pageSize")
)

type RName string

const (
	RSingleResource = RName("singleResource")
	RAllResources = RName("allResources")
)

type Processor interface {
	Process(ctx *ProcessorContext) error
}

func NewSerialProcessor(processors ...Processor) Processor {
	return &SerialProcessor{processors: processors}
}

type SerialProcessor struct {
	processors []Processor
}

func (sp *SerialProcessor) Process(ctx *ProcessorContext) error {
	for _, p := range sp.processors {
		err := p.Process(ctx)
		if nil != err {
			return err
		}
	}
	return nil
}
