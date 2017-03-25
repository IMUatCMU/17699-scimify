package processor

import (
	"encoding/json"
	"fmt"
	"github.com/go-scim/scimify/modify"
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/resource"
	"github.com/spf13/viper"
	"sync"
)

var (
	oneUserPatchParser,
	oneGroupPatchParser sync.Once

	userPatchParser,
	groupPatchParser Processor
)

func ParseParamForUserPatchEndpointProcessor() Processor {
	oneUserPatchParser.Do(func() {
		userPatchParser = &parseParamForPatchEndpointProcessor{
			internalSchemaRepo: persistence.GetInternalSchemaRepository(),
			schemaId:           viper.GetString("scim.internalSchemaId.user"),
			resourceIdUrlParam: viper.GetString("scim.api.userIdUrlParam"),
		}
	})
	return userPatchParser
}

func ParseParamForGroupPatchEndpointProcessor() Processor {
	oneGroupPatchParser.Do(func() {
		groupPatchParser = &parseParamForPatchEndpointProcessor{
			internalSchemaRepo: persistence.GetInternalSchemaRepository(),
			schemaId:           viper.GetString("scim.internalSchemaId.group"),
			resourceIdUrlParam: viper.GetString("scim.api.groupIdUrlParam"),
		}
	})
	return groupPatchParser
}

type parseParamForPatchEndpointProcessor struct {
	internalSchemaRepo persistence.Repository
	schemaId           string
	resourceIdUrlParam string
}

func (rep *parseParamForPatchEndpointProcessor) Process(ctx *ProcessorContext) error {
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

	if mod, err := rep.parseModification(req); err != nil {
		return err
	} else {
		ctx.Mod = mod
	}

	return nil
}

func (rep *parseParamForPatchEndpointProcessor) parseModification(req RequestSource) (*modify.Modification, error) {
	bodyBytes, err := req.Body()
	if err != nil {
		return nil, resource.CreateError(resource.ServerError, fmt.Sprintf("failed to read request body: %s", err.Error()))
	}

	mod := &modify.Modification{}
	err = json.Unmarshal(bodyBytes, mod)
	if err != nil {
		return nil, resource.CreateError(resource.InvalidSyntax, fmt.Sprintf("failed to serialize request body: %s", err.Error()))
	}

	return mod, nil
}

func (rep *parseParamForPatchEndpointProcessor) getResourceId(req RequestSource) (string, error) {
	if id := req.UrlParam(rep.resourceIdUrlParam); len(id) == 0 {
		return "", resource.CreateError(resource.InvalidSyntax, "failed to obtain resource id from url")
	} else {
		return id, nil
	}
}

func (rep *parseParamForPatchEndpointProcessor) getSchema() (*resource.Schema, error) {
	obj, err := rep.internalSchemaRepo.Get(rep.schemaId)
	if err != nil {
		return nil, resource.CreateError(resource.ServerError, fmt.Sprintf("failed to get schema for resource replace: %s", err.Error()))
	} else {
		return obj.(*resource.Schema), nil
	}
}

func (rep *parseParamForPatchEndpointProcessor) getRequestSource(ctx *ProcessorContext) RequestSource {
	if ctx.Request == nil {
		panic(&MissingContextValueError{"request source"})
	}
	return ctx.Request
}
