package service

import (
	"github.com/go-scim/scimify/resource"
	"github.com/go-scim/scimify/worker"
	"net/http"
)

func handleError(err error) (int, map[string]string, []byte) {
	var scimErr resource.Error
	if e, ok := err.(resource.Error); !ok {
		scimErr = resource.CreateError(resource.ServerError, err.Error())
	} else {
		scimErr = e
	}

	serializer := worker.GetDefaultJsonSerializerWorker()
	bytes, err := serializer.Do(&worker.JsonSerializeInput{Resource: scimErr.AsResource()})
	if nil != err {
		return http.StatusInternalServerError,
			map[string]string{"Content-Type": "text/plain"},
			[]byte(scimErr.Detail)
	} else {
		return scimErr.StatusCode,
			map[string]string{"Content-Type": "application/json+scim"},
			bytes.([]byte)
	}
}
