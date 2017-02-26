package defaults

import (
	"context"
	"errors"
	"github.com/go-scim/scimify/resource"
)

type metaUpdateValueDefaulter struct{}

func (_ *metaUpdateValueDefaulter) Default(r *resource.Resource, ctx context.Context) (bool, error) {
	if meta, ok := r.Attributes["meta"].(map[string]interface{}); !ok {
		return false, errors.New("meta was not set")
	} else {
		newVer, err := BumpVersion(meta["version"].(string))
		if err != nil {
			return false, errors.New("failed to generate new version: " + err.Error())
		}
		r.Lock()
		meta["version"] = newVer
		meta["lastModified"] = CurrentTime()
		r.Attributes["meta"] = meta
		r.Unlock()
		return true, nil
	}
}
