package service

import (
	"fmt"
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/resource"
	"github.com/go-scim/scimify/worker"
	"github.com/go-zoo/bone"
	"net/http"
)

func getAllSchemas(rw http.ResponseWriter, _ *http.Request) {
	var (
		statusCode int
		headers    map[string]string
		body       []byte
	)

	repository := persistence.GetSchemaRepository()
	serializer := worker.GetDefaultJsonSerializerWorker()

	if allSchemas, _ := repository.GetAll(); len(allSchemas) == 0 {
		e := resource.CreateError(resource.NotFound, "No schema was found.")
		statusCode, headers, body = handleError(e)
		writeResponse(rw, statusCode, headers, body)
		return
	} else if bytes, err := serializer.Do(&worker.JsonSerializeInput{Target: allSchemas}); err != nil {
		e := resource.CreateError(resource.ServerError, fmt.Sprintf("Error occured during serializing schema: %s", err.Error()))
		statusCode, headers, body = handleError(e)
		writeResponse(rw, statusCode, headers, body)
		return
	} else {
		statusCode = http.StatusOK
		headers = map[string]string{"Content-Type": "application/json+scim"}
		body = bytes.([]byte)
		writeResponse(rw, statusCode, headers, body)
	}
}

func getSchemaById(rw http.ResponseWriter, req *http.Request) {
	var (
		statusCode int
		headers    map[string]string
		body       []byte
	)

	schemaId := bone.GetValue(req, "schemaId")
	repository := persistence.GetSchemaRepository()
	serializer := worker.GetDefaultJsonSerializerWorker()

	schema, _ := repository.Get(schemaId, nil)
	if nil == schema {
		e := resource.CreateError(resource.NotFound, fmt.Sprintf("Schema by id '%s' does not exist.", schemaId))
		statusCode, headers, body = handleError(e)
		writeResponse(rw, statusCode, headers, body)
		return
	} else if bytes, err := serializer.Do(&worker.JsonSerializeInput{Target: schema}); nil != err {
		e := resource.CreateError(resource.ServerError, fmt.Sprintf("Error occured during serializing schema: %s", err.Error()))
		statusCode, headers, body = handleError(e)
		writeResponse(rw, statusCode, headers, body)
		return
	} else {
		statusCode = http.StatusOK
		headers = map[string]string{"Content-Type": "application/json+scim"}
		body = bytes.([]byte)
		writeResponse(rw, statusCode, headers, body)
	}
}
