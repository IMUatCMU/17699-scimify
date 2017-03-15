package processor

import (
	"github.com/satori/go.uuid"
	"strings"
)

type generateIdProcessor struct{}

func (gip *generateIdProcessor) Process(ctx *ProcessorContext) error {
	r := getResource(ctx, true)
	r.Attributes["id"] = strings.ToLower(uuid.NewV4().String())
	return nil
}
