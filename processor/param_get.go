package processor

import (
	"fmt"
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/resource"
	"github.com/go-zoo/bone"
	"net/http"
)

type parseParamForGetEndpointProcessor struct {
	internalSchemaRepo persistence.Repository
	schemaId           string
	userIdUrlParam     string
}

func (gep *parseParamForGetEndpointProcessor) Process(ctx *ProcessorContext) error {
	httpRequest := gep.getHttpRequest(ctx)

	if sch, err := gep.getSchema(); err != nil {
		return err
	} else {
		ctx.Schema = sch
	}

	if id, err := gep.getUserId(httpRequest); len(id) == 0 {
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

func (gep *parseParamForGetEndpointProcessor) getUserId(req *http.Request) (string, error) {
	if id := bone.GetValue(req, gep.userIdUrlParam); len(id) == 0 {
		return "", resource.CreateError(resource.InvalidSyntax, "failed to obtain resource id from url")
	} else {
		return id, nil
	}
}

func (gep *parseParamForGetEndpointProcessor) getHttpRequest(ctx *ProcessorContext) *http.Request {
	if ctx.Request == nil {
		panic(&MissingContextValueError{"http request"})
	}
	return ctx.Request
}
