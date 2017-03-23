package processor

import (
	"github.com/go-scim/scimify/modify"
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
	QueryFilter    string
	QuerySortBy    string
	QuerySortOrder bool
	QueryPageStart int
	QueryPageSize  int
	ParsedFilter   interface{}

	// inclusion and exclusion
	Inclusion []string
	Exclusion []string

	// modification
	Mod *modify.Modification

	// Serialization
	SerializationTargetFunc func() interface{}

	// Error
	Err error

	// Results
	SingleResult resource.ScimObject
	MultiResults []resource.ScimObject
	ListResponse *resource.ListResponse

	// HTTP Response
	ResponseStatus  int
	ResponseHeaders map[string]string
	ResponseBody    []byte
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
	return &ErrorHandlingProcessor{Op: opProc, ErrOp: errProc}
}

type ErrorHandlingProcessor struct {
	Op    []Processor
	ErrOp []Processor
}

func (ehp *ErrorHandlingProcessor) Process(ctx *ProcessorContext) error {
	for _, op := range ehp.Op {
		err := op.Process(ctx)
		if nil != err {
			ctx.Err = err
			for _, ep := range ehp.ErrOp {
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
