package processor

import (
	"fmt"
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/resource"
	"github.com/spf13/viper"
	"sync"
)

var (
	oneUserReplaceParser,
	oneGroupReplaceParser sync.Once

	userReplaceParser,
	groupReplaceParser Processor
)

func ParseParamForUserReplaceEndpointProcessor() Processor {
	oneUserReplaceParser.Do(func() {
		userReplaceParser = &parseParamForReplaceEndpointProcessor{
			internalSchemaRepo: persistence.GetInternalSchemaRepository(),
			schemaId:           viper.GetString("scim.internalSchemaId.user"),
			resourceIdUrlParam: viper.GetString("scim.api.userIdUrlParam"),
		}
	})
	return userReplaceParser
}

func ParseParamForGroupReplaceEndpointProcessor() Processor {
	oneGroupReplaceParser.Do(func() {
		groupReplaceParser = &parseParamForReplaceEndpointProcessor{
			internalSchemaRepo: persistence.GetInternalSchemaRepository(),
			schemaId:           viper.GetString("scim.internalSchemaId.group"),
			resourceIdUrlParam: viper.GetString("scim.api.groupIdUrlParam"),
		}
	})
	return groupReplaceParser
}

type parseParamForReplaceEndpointProcessor struct {
	internalSchemaRepo persistence.Repository
	schemaId           string
	resourceIdUrlParam string
}

func (rep *parseParamForReplaceEndpointProcessor) Process(ctx *ProcessorContext) error {
	req := rep.getRequestSource(ctx)

	if sch, err := rep.getSchema(); err != nil {
		return err
	} else {
		ctx.Schema = sch
	}

	if id, err := rep.getResourceId(req); len(id) == 0 {
		return err
	} else {
		ctx.Identity = id
	}

	if r, err := rep.parseResource(req); err != nil {
		return err
	} else {
		ctx.Resource = r
	}

	return nil
}

func (rep *parseParamForReplaceEndpointProcessor) parseResource(req RequestSource) (*resource.Resource, error) {
	bodyBytes, err := req.Body()
	if err != nil {
		return nil, resource.CreateError(resource.ServerError, fmt.Sprintf("failed to read request body: %s", err.Error()))
	}

	r, err := resource.NewResourceFromBytes(bodyBytes)
	if err != nil {
		return nil, resource.CreateError(resource.InvalidSyntax, fmt.Sprintf("failed to read serialize request body: %s", err.Error()))
	}

	return r, nil
}

func (rep *parseParamForReplaceEndpointProcessor) getResourceId(req RequestSource) (string, error) {
	if id := req.UrlParam(rep.resourceIdUrlParam); len(id) == 0 {
		return "", resource.CreateError(resource.InvalidSyntax, "failed to obtain resource id from url")
	} else {
		return id, nil
	}
}

func (rep *parseParamForReplaceEndpointProcessor) getSchema() (*resource.Schema, error) {
	obj, err := rep.internalSchemaRepo.Get(rep.schemaId)
	if err != nil {
		return nil, resource.CreateError(resource.ServerError, fmt.Sprintf("failed to get schema for resource replace: %s", err.Error()))
	} else {
		return obj.(*resource.Schema), nil
	}
}

func (rep *parseParamForReplaceEndpointProcessor) getRequestSource(ctx *ProcessorContext) RequestSource {
	if ctx.Request == nil {
		panic(&MissingContextValueError{"request source"})
	}
	return ctx.Request
}
