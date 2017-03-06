package service

import (
	"fmt"
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/resource"
	"github.com/go-scim/scimify/worker"
	"github.com/spf13/viper"
	"net/http"
)

type resourceTypeService struct{}

func (srv *resourceTypeService) getAllResourceTypes(_ *http.Request) (response, error) {
	repository := persistence.GetResourceTypeRepository()
	serializer := worker.GetDefaultJsonSerializerWorker()

	if allResourceTypes, _ := repository.GetAll(); len(allResourceTypes) == 0 {
		return nil_response, resource.CreateError(resource.NotFound, "No resource type was found.")
	} else if bytes, err := serializer.Do(&worker.JsonSerializeInput{
		Target: resource.NewListResponse(allResourceTypes, 1, viper.GetInt("scim.itemsPerPage"), len(allResourceTypes)),
	}); err != nil {
		return nil_response, resource.CreateError(resource.ServerError, fmt.Sprintf("Error occured during serializing resource types: %s", err.Error()))
	} else {
		return response{
			statusCode: http.StatusOK,
			headers:    nil,
			body:       bytes.([]byte),
		}, nil
	}
}
