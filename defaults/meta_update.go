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
		id, ok := r.Attributes["id"].(string)
		if !ok {
			return false, errors.New("id must be set before updating meta")
		}

		newMeta := make(map[string]interface{})
		for k, v := range meta {
			newMeta[k] = v
		}
		newMeta["version"] = GenerateNewVersion(id)
		newMeta["lastModified"] = CurrentTime()

		r.Lock()
		r.Attributes["meta"] = newMeta
		r.Unlock()
		return true, nil
	}
}
