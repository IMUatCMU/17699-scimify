package defaults

import (
	"context"
	"github.com/go-scim/scimify/resource"
)

type ValueDefaulter interface {
	Default(*resource.Resource, context.Context) (bool, error)
}
