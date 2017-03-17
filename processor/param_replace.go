package processor

import (
	"fmt"
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/resource"
	"github.com/go-zoo/bone"
	"io/ioutil"
	"net/http"
)

type parseParamForReplaceEndpointProcessor struct {
	internalSchemaRepo persistence.Repository
	schemaId           string
	userIdUrlParam     string
}

func (rep *parseParamForReplaceEndpointProcessor) Process(ctx *ProcessorContext) error {
	httpRequest := rep.getHttpRequest(ctx)

	if sch, err := rep.getSchema(); err != nil {
		return err
	} else {
		ctx.Schema = sch
	}

	if id, err := rep.getUserId(httpRequest); len(id) == 0 {
		return err
	} else {
		ctx.Identity = id
	}

	if r, err := rep.parseResource(httpRequest); err != nil {
		return err
	} else {
		ctx.Resource = r
	}

	return nil
}

func (rep *parseParamForReplaceEndpointProcessor) parseResource(req *http.Request) (*resource.Resource, error) {
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

func (rep *parseParamForReplaceEndpointProcessor) getUserId(req *http.Request) (string, error) {
	if id := bone.GetValue(req, rep.userIdUrlParam); len(id) == 0 {
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

func (rep *parseParamForReplaceEndpointProcessor) getHttpRequest(ctx *ProcessorContext) *http.Request {
	if ctx.Request == nil {
		panic(&MissingContextValueError{"http request"})
	}
	return ctx.Request
}
