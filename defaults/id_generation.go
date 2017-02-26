package defaults

import (
	"github.com/go-scim/scimify/resource"
	"context"
	"github.com/satori/go.uuid"
	"strings"
)

type idGenerationValueDefaulter struct {}

func (d *idGenerationValueDefaulter) Default(r *resource.Resource, ctx context.Context) (bool, error) {
	r.Lock()
	r.Attributes["id"] = strings.ToLower(uuid.NewV4().String())
	r.Unlock()
	return true, nil
}
