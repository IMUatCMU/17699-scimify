package processor

import (
	"github.com/go-scim/scimify/resource"
	"net/http"
)

type ProcessorContext struct {
	// basic
	Identity  string
	Resource  *resource.Resource
	Reference *resource.Resource
	Schema    *resource.Schema
	Request   *http.Request

	// query args
	QueryFilter	string
	QuerySortBy	string
	QuerySortOrder	bool
	QueryPageStart	int
	QueryPageSize	int
	ParsedFilter 	interface{}

	// inclusion and exclusion
	Inclusion	[]string
	Exclusion	[]string

	// Serialization
	SerializationTargetFunc func() interface{}

	// Error
	Err 		error

	// Results
	SingleResult	resource.ScimObject
	MultiResults	[]resource.ScimObject
	ListResponse 	*resource.ListResponse

	// HTTP Response
	ResponseStatus	int
	ResponseHeaders	map[string]string
	ResponseBody	[]byte
}

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
			ctx.Err = err
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
