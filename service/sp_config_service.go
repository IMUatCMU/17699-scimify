package service

import (
	"fmt"
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/resource"
	"github.com/go-scim/scimify/worker"
	"net/http"
)

type spConfigService struct{}

func (srv *spConfigService) getServiceProviderConfig(_ *http.Request) (response, error) {
	repository := persistence.GetServiceProviderConfigRepository()
	serializer := worker.GetDefaultJsonSerializerWorker()

	if spConfig, err := repository.Get("", nil); err != nil {
		return nil_response, resource.CreateError(resource.NotFound, "No service provider configuration was found.")
	} else if bytes, err := serializer.Do(&worker.JsonSerializeInput{Target: spConfig}); err != nil {
		return nil_response, resource.CreateError(resource.ServerError, fmt.Sprintf("Error occured during serializing service provider configuration: %s", err.Error()))
	} else {
		return response{
			statusCode: http.StatusOK,
			body:       bytes.([]byte),
		}, nil
	}
}
