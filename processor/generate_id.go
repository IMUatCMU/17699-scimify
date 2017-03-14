package processor

import (
	"context"
	"github.com/go-scim/scimify/resource"
	"github.com/satori/go.uuid"
	"strings"
)

type generateIdProcessor struct{}

func (gip *generateIdProcessor) Process(r *resource.Resource, ctx context.Context) error {
	r.Attributes["id"] = strings.ToLower(uuid.NewV4().String())
	return nil
}
