package processor

import (
	"context"
	"fmt"
	"github.com/go-scim/scimify/resource"
	"github.com/spf13/viper"
	"strings"
)

type generateMetaProcessor struct {
	ResourceType    string
	ResourceTypeUri string
}

func (gmp *generateMetaProcessor) Process(r *resource.Resource, ctx context.Context) error {
	id, ok := r.Attributes["id"].(string)
	if !ok || len(id) == 0 {
		panic(&PrerequisiteFailedError{reporter: "meta generation", requirement: "id"})
	}

	now := getCurrentTime()
	meta := map[string]interface{}{
		"resourceType": gmp.ResourceType,
		"created":      now,
		"lastModified": now,
		"version":      generateNewVersion(id),
		"location": fmt.Sprintf("%s/%s/%s",
			strings.Trim(viper.GetString("server.rootPath"), "/"),
			strings.Trim(gmp.ResourceTypeUri, "/"),
			id,
		),
	}
	r.Attributes["meta"] = meta

	return nil
}
