package processor

import (
	"fmt"
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/resource"
	"io/ioutil"
	"net/http"
)

type parseParamForCreateEndpointProcessor struct {
	internalSchemaRepo persistence.Repository
	schemaId           string
}

func (cep *parseParamForCreateEndpointProcessor) Process(ctx *ProcessorContext) error {
	httpRequest := cep.getHttpRequest(ctx)

	if sch, err := cep.getSchema(); err != nil {
		return err
	} else {
		ctx.Schema = sch
	}

	if r, err := cep.parseResource(httpRequest); err != nil {
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

func (cep *parseParamForCreateEndpointProcessor) parseResource(req *http.Request) (*resource.Resource, error) {
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, resource.CreateError(resource.ServerError, fmt.Sprintf("failed to read request body: %s", err.Error()))
	}

	r, err := resource.NewResourceFromBytes(bodyBytes)
	if err != nil {
		return nil, resource.CreateError(resource.InvalidSyntax, fmt.Sprintf("failed to read serialize request body: %s", err.Error()))
	}

	return r, nil
}

func (cep *parseParamForCreateEndpointProcessor) getHttpRequest(ctx *ProcessorContext) *http.Request {
	if ctx.Request == nil {
		panic(&MissingContextValueError{"http request"})
	}
	return ctx.Request
}
