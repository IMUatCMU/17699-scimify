package defaults

import (
	"context"
	"github.com/go-scim/scimify/resource"
)

type delegateValueDefaulter struct {
	Defaulters []ValueDefaulter
}

func (d *delegateValueDefaulter) Default(r *resource.Resource, ctx context.Context) (bool, error) {
	for _, defaulter := range d.Defaulters {
		_, err := defaulter.Default(r, ctx)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}
