package processor

import (
	"fmt"
	"github.com/go-scim/scimify/resource"
	"github.com/spf13/viper"
	"strings"
	"sync"
)

var (
	oneGenerateUserMeta,
	oneGenerateGroupMeta sync.Once

	generateUserMetaProcessor,
	generateGroupMetaProcessor Processor
)

func GenerateUserMetaProcessor() Processor {
	oneGenerateUserMeta.Do(func() {
		generateUserMetaProcessor = &generateMetaProcessor{
			ResourceType:    viper.GetString("scim.resourceType.user"),
			ResourceTypeUri: viper.GetString("scim.resourceTypeUri.user"),
		}
	})
	return generateUserMetaProcessor
}

func GenerateGroupMetaProcessor() Processor {
	oneGenerateGroupMeta.Do(func() {
		generateGroupMetaProcessor = &generateMetaProcessor{
			ResourceType:    viper.GetString("scim.resourceType.group"),
			ResourceTypeUri: viper.GetString("scim.resourceTypeUri.group"),
		}
	})
	return generateGroupMetaProcessor
}

type generateMetaProcessor struct {
	ResourceType    string
	ResourceTypeUri string
}

func (gmp *generateMetaProcessor) Process(ctx *ProcessorContext) error {
	r := gmp.getResource(ctx)

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

func (gmp *generateMetaProcessor) getResource(ctx *ProcessorContext) *resource.Resource {
	if ctx.Resource == nil {
		panic(&MissingContextValueError{"resource"})
	}
	return ctx.Resource
}
