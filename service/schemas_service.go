package service

import (
	"encoding/json"
	"fmt"
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/resource"
	"github.com/go-scim/scimify/worker"
	"github.com/go-zoo/bone"
	"net/http"
)

func getAllSchemas(rw http.ResponseWriter, req *http.Request) {
	var (
		statusCode int
		headers    map[string]string
		body       []byte
	)

	repository := persistence.GetSchemaRepository()
	serializer := worker.GetDefaultJsonSerializerWorker()

	allSchemas, _ := repository.GetAll()
	if len(allSchemas) == 0 {
		e := resource.CreateError(resource.NotFound, "No schema was found.")
		statusCode, headers, body = handleError(e)
		writeResponse(rw, statusCode, headers, body)
		return
	} else {
		rawJsons := make([]json.RawMessage, 0, len(allSchemas))
		for _, schema := range allSchemas {
			if bytes, err := serializer.Do(&worker.JsonSerializeInput{Resource: schema}); err != nil {
				e := resource.CreateError(resource.ServerError, fmt.Sprintf("Error occured during serializing schema: %s", err.Error()))
				statusCode, headers, body = handleError(e)
				writeResponse(rw, statusCode, headers, body)
				return
			} else {
				rawJsons = append(rawJsons, json.RawMessage(bytes.([]byte)))
			}
		}
		if json, err := json.Marshal(&rawJsons); err != nil {
			e := resource.CreateError(resource.ServerError, fmt.Sprintf("Error occured during serializing schema: %s", err.Error()))
			statusCode, headers, body = handleError(e)
		} else {
			statusCode = http.StatusOK
			headers = map[string]string{"Content-Type": "application/json+scim"}
			body = json
		}
	}

	writeResponse(rw, statusCode, headers, body)
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
	} else if bytes, err := serializer.Do(&worker.JsonSerializeInput{Resource: schema}); nil != err {
		e := resource.CreateError(resource.ServerError, fmt.Sprintf("Error occured during serializing schema: %s", err.Error()))
		statusCode, headers, body = handleError(e)
		writeResponse(rw, statusCode, headers, body)
		return
	} else {
		statusCode = http.StatusOK
		headers = map[string]string{"Content-Type": "application/json+scim"}
		body = bytes.([]byte)
	}

	writeResponse(rw, statusCode, headers, body)
}
