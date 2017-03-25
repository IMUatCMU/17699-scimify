package processor

import (
	"fmt"
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/resource"
	"github.com/spf13/viper"
	"sync"
)

var (
	oneParseUserCreate,
	oneParseGroupCreate sync.Once

	parseUserCreate,
	parseGroupCreate Processor
)

func ParseParamForUserCreateEndpointProcessor() Processor {
	oneParseUserCreate.Do(func() {
		parseUserCreate = &parseParamForCreateEndpointProcessor{
			internalSchemaRepo: persistence.GetInternalSchemaRepository(),
			schemaId:           viper.GetString("scim.internalSchemaId.user"),
		}
	})
	return parseUserCreate
}

func ParseParamForGroupCreateEndpointProcessor() Processor {
	oneParseGroupCreate.Do(func() {
		parseGroupCreate = &parseParamForCreateEndpointProcessor{
			internalSchemaRepo: persistence.GetInternalSchemaRepository(),
			schemaId:           viper.GetString("scim.internalSchemaId.group"),
		}
	})
	return parseGroupCreate
}

type parseParamForCreateEndpointProcessor struct {
	internalSchemaRepo persistence.Repository
	schemaId           string
}

func (cep *parseParamForCreateEndpointProcessor) Process(ctx *ProcessorContext) error {
	requestSource := cep.getRequestSource(ctx)

	if sch, err := cep.getSchema(); err != nil {
		return err
	} else {
		ctx.Schema = sch
	}

	if r, err := cep.parseResource(requestSource); err != nil {
		return err
	} else {
		ctx.Resource = r
	}

	return nil
}

func (cep *parseParamForCreateEndpointProcessor) getSchema() (*resource.Schema, error) {
	obj, err := cep.internalSchemaRepo.Get(cep.schemaId)
	if err != nil {
		return nil, resource.CreateError(resource.ServerError, fmt.Sprintf("failed to get schema for resource creation: %s", err.Error()))
	} else {
		return obj.(*resource.Schema), nil
	}
}

func (cep *parseParamForCreateEndpointProcessor) parseResource(req RequestSource) (*resource.Resource, error) {
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

func (cep *parseParamForCreateEndpointProcessor) getRequestSource(ctx *ProcessorContext) RequestSource {
	if ctx.Request == nil {
		panic(&MissingContextValueError{"request source"})
	}
	return ctx.Request
}
