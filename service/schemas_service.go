package service

import (
	"fmt"
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/resource"
	"github.com/go-scim/scimify/worker"
	"github.com/go-zoo/bone"
	"net/http"
)

func getAllSchemas(rw http.ResponseWriter, req *http.Request) {

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
	} else if bytes, err := serializer.Do(&worker.JsonSerializeInput{Resource: schema}); nil != err {
		e := resource.CreateError(resource.ServerError, fmt.Sprintf("Error occured during serializing schema: %s", err.Error()))
		statusCode, headers, body = handleError(e)
	} else {
		statusCode = http.StatusOK
		headers = map[string]string{"Content-Type": "application/json+scim"}
		body = bytes.([]byte)
	}

	for k, v := range headers {
		rw.Header().Set(k, v)
	}
	rw.WriteHeader(statusCode)
	rw.Write(body)
}
