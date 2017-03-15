package processor

import (
	"github.com/go-scim/scimify/resource"
	"net/http"
)

type ProcessorContext struct {
	Identity  string
	Resource  *resource.Resource
	Reference *resource.Resource
	Schema    *resource.Schema
	Request   *http.Request
	MiscArgs  map[AName]interface{}
	Results   map[RName]interface{}
}

type AName string

const (
	ArgResource     = AName("resource")
	ArgReference    = AName("reference")
	ArgSchema       = AName("schema")
	ArgFilter       = AName("filter")
	ArgSortBy       = AName("sortBy")
	ArgSortOrder    = AName("sortOrder")
	ArgPageStart    = AName("pageStart")
	ArgPageSize     = AName("pageSize")
	ArgIncludePaths = AName("includePaths")
	ArgExcludePaths = AName("excludePaths")
	ArgError        = AName("error")
)

type RName string

const (
	RSingleResource = RName("singleResource")
	RAllResources   = RName("allResources")
	RFinalError     = RName("finalError")
	RBodyBytes      = RName("bodyBytes")
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

func NewErrorHandlingProcessor(opProc []Processor, errProc []Processor) Processor {
	return &ErrorHandlingProcessor{opProcessors: opProc, errProcessors: errProc}
}

type ErrorHandlingProcessor struct {
	opProcessors  []Processor
	errProcessors []Processor
}

func (ehp *ErrorHandlingProcessor) Process(ctx *ProcessorContext) error {
	for _, op := range ehp.opProcessors {
		err := op.Process(ctx)
		if nil != err {
			ctx.MiscArgs[ArgError] = err
			for _, ep := range ehp.errProcessors {
				err := ep.Process(ctx)
				if nil != err {
					return err
				}
			}
			return nil
		}
	}
	return nil
}
