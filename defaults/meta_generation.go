package defaults

import (
	"context"
	"fmt"
	"github.com/go-scim/scimify/resource"
	"github.com/spf13/viper"
	"strings"
)

type metaGenerationValueDefaulter struct{}

func (_ *metaGenerationValueDefaulter) Default(r *resource.Resource, ctx context.Context) (bool, error) {
	resourceType, ok := ctx.Value(resource.CK_ResourceType).(string)
	if !ok {
		panic("missing required context parmaeter: CK_ResourceType")
	}

	resourceTypeURI, ok := ctx.Value(resource.CK_ResourceTypeURI).(string)
	if !ok {
		panic("missing required context parmaeter: CK_ResourceTypeURI")
	}

	id, ok := r.Attributes["id"].(string)
	if !ok || len(id) == 0 {
		panic("id has not been set")
	}

	now := CurrentTime()
	meta := map[string]interface{}{
		"resourceType": resourceType,
		"created":      now,
		"lastModified": now,
		"version":      GenerateNewVersion(id),
		"location": fmt.Sprintf("%s/%s/%s",
			strings.Trim(viper.GetString("server.rootPath"), "/"),
			strings.Trim(resourceTypeURI, "/"),
			id,
		),
	}

	r.Lock()
	r.Attributes["meta"] = meta
	r.Unlock()
	return true, nil
}
