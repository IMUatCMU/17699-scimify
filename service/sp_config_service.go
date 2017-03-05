package service

import (
	"fmt"
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/resource"
	"github.com/go-scim/scimify/worker"
	"net/http"
)

func getServiceProviderConfig(rw http.ResponseWriter, _ *http.Request) {
	var (
		statusCode int
		headers    map[string]string
		body       []byte
	)

	repository := persistence.GetServiceProviderConfigRepository()
	serializer := worker.GetDefaultJsonSerializerWorker()

	if spConfig, err := repository.Get("", nil); err != nil {
		e := resource.CreateError(resource.NotFound, "No service provider configuration was found.")
		statusCode, headers, body = handleError(e)
		writeResponse(rw, statusCode, headers, body)
		return
	} else if bytes, err := serializer.Do(&worker.JsonSerializeInput{Target: spConfig}); err != nil {
		e := resource.CreateError(resource.ServerError, fmt.Sprintf("Error occured during serializing service provider configuration: %s", err.Error()))
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
