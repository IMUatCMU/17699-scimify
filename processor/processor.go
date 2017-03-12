package processor

import (
	"context"
	"github.com/go-scim/scimify/resource"
)

type Processor interface {
	Process(r *resource.Resource, ctx context.Context) error
}

type SerialProcessor struct {
	processors []Processor
}

func (sp *SerialProcessor) Process(r *resource.Resource, ctx context.Context) error {
	for _, p := range sp.processors {
		err := p.Process(r, ctx)
		if nil != err {
			return err
		}
	}
	return nil
}
