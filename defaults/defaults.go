package defaults

import (
	"github.com/go-scim/scimify/resource"
	"context"
)

type ValueDefaulter interface {
	Default(*resource.Resource, context.Context) (bool, error)
}
