package processor

import (
	"fmt"
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/resource"
	"github.com/spf13/viper"
	"sync"
)

var (
	oneParseUserGet,
	oneParseGroupGet,
	oneParseSchemaGet sync.Once

	parseUserGet,
	parseGroupGet,
	parseSchemaGet Processor
)

func ParseParamForSchemaGetEndpointProcessor() Processor {
	oneParseSchemaGet.Do(func() {
		parseSchemaGet = &parseParamForGetEndpointProcessor{
			resourceIdUrlParam: viper.GetString("scim.api.schemaIdUrlParam"),
		}
	})
	return parseSchemaGet
}

func ParseParamForUserGetEndpointProcessor() Processor {
	oneParseUserGet.Do(func() {
		parseUserGet = &parseParamForGetEndpointProcessor{
			internalSchemaRepo: persistence.GetInternalSchemaRepository(),
			schemaId:           viper.GetString("scim.internalSchemaId.user"),
			resourceIdUrlParam: viper.GetString("scim.api.userIdUrlParam"),
		}
	})
	return parseUserGet
}

func ParseParamForGroupGetEndpointProcessor() Processor {
	oneParseGroupGet.Do(func() {
		parseGroupGet = &parseParamForGetEndpointProcessor{
			internalSchemaRepo: persistence.GetInternalSchemaRepository(),
			schemaId:           viper.GetString("scim.internalSchemaId.group"),
			resourceIdUrlParam: viper.GetString("scim.api.groupIdUrlParam"),
		}
	})
	return parseGroupGet
}

type parseParamForGetEndpointProcessor struct {
	internalSchemaRepo persistence.Repository
	schemaId           string
	resourceIdUrlParam string
}

func (gep *parseParamForGetEndpointProcessor) Process(ctx *ProcessorContext) error {
	req := gep.getRequestSource(ctx)

	if gep.internalSchemaRepo != nil {
		if sch, err := gep.getSchema(); err != nil {
			return err
		} else {
			ctx.Schema = sch
		}
	}

	if id, err := gep.getResourceId(req); len(id) == 0 {
		return err
	} else {
		ctx.Identity = id
	}

	// TODO parse version if any

	return nil
}

func (gep *parseParamForGetEndpointProcessor) getSchema() (*resource.Schema, error) {
	obj, err := gep.internalSchemaRepo.Get(gep.schemaId)
	if err != nil {
		return nil, resource.CreateError(resource.ServerError, fmt.Sprintf("failed to get schema for resource get: %s", err.Error()))
	} else {
		return obj.(*resource.Schema), nil
	}
}

func (gep *parseParamForGetEndpointProcessor) getResourceId(req RequestSource) (string, error) {
	if id := req.UrlParam(gep.resourceIdUrlParam); len(id) == 0 {
		return "", resource.CreateError(resource.InvalidSyntax, "failed to obtain resource id from url")
	} else {
		return id, nil
	}
}

func (gep *parseParamForGetEndpointProcessor) getRequestSource(ctx *ProcessorContext) RequestSource {
	if ctx.Request == nil {
		panic(&MissingContextValueError{"request source"})
	}
	return ctx.Request
}
