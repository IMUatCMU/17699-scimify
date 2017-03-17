package processor

import (
	"fmt"
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/resource"
	"github.com/go-zoo/bone"
	"net/http"
)

type parseParamForDeleteEndpointProcessor struct {
	internalSchemaRepo persistence.Repository
	schemaId           string
	userIdUrlParam     string
}

func (dep *parseParamForDeleteEndpointProcessor) Process(ctx *ProcessorContext) error {
	httpRequest := dep.getHttpRequest(ctx)

	if sch, err := dep.getSchema(); err != nil {
		return err
	} else {
		ctx.Schema = sch
	}

	if id, err := dep.getUserId(httpRequest); len(id) == 0 {
		return err
	} else {
		ctx.Identity = id
	}

	// TODO parse version if any

	return nil
}

func (dep *parseParamForDeleteEndpointProcessor) getSchema() (*resource.Schema, error) {
	obj, err := dep.internalSchemaRepo.Get(dep.schemaId)
	if err != nil {
		return nil, resource.CreateError(resource.ServerError, fmt.Sprintf("failed to get schema for resource delete: %s", err.Error()))
	} else {
		return obj.(*resource.Schema), nil
	}
}

func (dep *parseParamForDeleteEndpointProcessor) getUserId(req *http.Request) (string, error) {
	if id := bone.GetValue(req, dep.userIdUrlParam); len(id) == 0 {
		return "", resource.CreateError(resource.InvalidSyntax, "failed to obtain resource id from url")
	} else {
		return id, nil
	}
}

func (dep *parseParamForDeleteEndpointProcessor) getHttpRequest(ctx *ProcessorContext) *http.Request {
	if ctx.Request == nil {
		panic(&MissingContextValueError{"http request"})
	}
	return ctx.Request
}
