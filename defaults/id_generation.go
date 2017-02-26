package defaults

import (
	"context"
	"github.com/go-scim/scimify/resource"
	"github.com/satori/go.uuid"
	"strings"
)

type idGenerationValueDefaulter struct{}

func (_ *idGenerationValueDefaulter) Default(r *resource.Resource, ctx context.Context) (bool, error) {
	r.Lock()
	r.Attributes["id"] = strings.ToLower(uuid.NewV4().String())
	r.Unlock()
	return true, nil
}
