package service

import (
	"fmt"
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/resource"
	"github.com/go-scim/scimify/worker"
	"github.com/go-zoo/bone"
	"net/http"
)

type schemaService struct{}

func (srv *schemaService) getAllSchemas(_ *http.Request) (response, error) {
	repository := persistence.GetSchemaRepository()
	serializer := worker.GetDefaultJsonSerializerWorker()

	if allSchemas, _ := repository.GetAll(); len(allSchemas) == 0 {
		return nil_response, resource.CreateError(resource.NotFound, "No schema was found.")
	} else if bytes, err := serializer.Do(&worker.JsonSerializeInput{Target: allSchemas}); err != nil {
		return nil_response, resource.CreateError(resource.ServerError, fmt.Sprintf("Error occured during serializing schema: %s", err.Error()))
	} else {
		return response{
			statusCode: http.StatusOK,
			body:       bytes.([]byte),
		}, nil
	}
}

func (srv *schemaService) getSchemaById(req *http.Request) (response, error) {
	schemaId := bone.GetValue(req, "schemaId")
	repository := persistence.GetSchemaRepository()
	serializer := worker.GetDefaultJsonSerializerWorker()

	schema, _ := repository.Get(schemaId, nil)
	if nil == schema {
		return nil_response, resource.CreateError(resource.NotFound, fmt.Sprintf("Schema by id '%s' does not exist.", schemaId))
	} else if bytes, err := serializer.Do(&worker.JsonSerializeInput{Target: schema}); nil != err {
		return nil_response, resource.CreateError(resource.ServerError, fmt.Sprintf("Error occured during serializing schema: %s", err.Error()))
	} else {
		return response{
			statusCode: http.StatusOK,
			body:       bytes.([]byte),
		}, nil
	}
}
