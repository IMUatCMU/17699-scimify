package service

import (
	"net/http"
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/worker"
	"github.com/go-scim/scimify/resource"
	"fmt"
	"github.com/spf13/viper"
)

func getAllResourceTypes(rw http.ResponseWriter, _ *http.Request) {
	var (
		statusCode int
		headers    map[string]string
		body       []byte
	)

	repository := persistence.GetResourceTypeRepository()
	serializer := worker.GetDefaultJsonSerializerWorker()

	if allResourceTypes, _ := repository.GetAll(); len(allResourceTypes) == 0 {
		e := resource.CreateError(resource.NotFound, "No resource type was found.")
		statusCode, headers, body = handleError(e)
		writeResponse(rw, statusCode, headers, body)
		return
	} else if bytes, err := serializer.Do(&worker.JsonSerializeInput{
		Target: resource.NewListResponse(allResourceTypes, 1, viper.GetInt("scim.itemsPerPage")),
	}); err != nil {
		e := resource.CreateError(resource.ServerError, fmt.Sprintf("Error occured during serializing resource types: %s", err.Error()))
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
