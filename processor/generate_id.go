package processor

import (
	"github.com/go-scim/scimify/resource"
	"context"
	"strings"
	"github.com/satori/go.uuid"
)

type generateIdProcessor struct {}

func (gip *generateIdProcessor) Process(r *resource.Resource, ctx context.Context) error {
	r.Attributes["id"] = strings.ToLower(uuid.NewV4().String())
	return nil
}