package service

import (
	"context"
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/resource"
	"github.com/go-scim/scimify/validation"
	"github.com/go-scim/scimify/worker"
	"io/ioutil"
	"net/http"
)

func getUserById(rw http.ResponseWriter, req *http.Request) {

}

func createUser(rw http.ResponseWriter, req *http.Request) {
	var (
		statusCode int
		headers    map[string]string
		body       []byte
		subject    *resource.Resource
	)

	// load resource from request body
	if bodyBytes, err := ioutil.ReadAll(req.Body); err != nil {
		e := resource.CreateError(resource.ServerError, "Failed to read request body.")
		statusCode, headers, body = handleError(e)
		writeResponse(rw, statusCode, headers, body)
		return
	} else if subject, err = resource.NewResourceFromBytes(bodyBytes); err != nil {
		e := resource.CreateError(resource.InvalidSyntax, "The request body message was invalid or did not conform to the request schema")
		statusCode, headers, body = handleError(e)
		writeResponse(rw, statusCode, headers, body)
		return
	}

	// get schema guideline
	schema, err := getUserSchema()
	if err != nil {
		e := resource.CreateError(resource.ServerError, "No schema was configured for user resource.")
		statusCode, headers, body = handleError(e)
		writeResponse(rw, statusCode, headers, body)
		return
	}

	// create context
	ctx := context.Background()
	ctx = context.WithValue(ctx, resource.CK_Schema, schema)
	ctx = context.WithValue(ctx, resource.CK_ResourceType, "User")
	ctx = context.WithValue(ctx, resource.CK_ResourceTypeURI, "/Users")

	// correct case via shared value defaulter
	sharedValueDefaulter := worker.GetSharedValueDefaulterWorker()
	if _, err = sharedValueDefaulter.Do(&worker.ValueDefaulterInput{
		Resource: subject,
		Context:  ctx,
	}); err != nil {
		e := resource.CreateError(resource.InvalidSyntax, err.Error())
		statusCode, headers, body = handleError(e)
		writeResponse(rw, statusCode, headers, body)
		return
	}

	// validate the resource via resource creation validator worker
	validator := worker.GetCreationValidatorWorker()
	if _, err = validator.Do(&worker.ValidationInput{
		Resource: subject,
		Context:  ctx,
		Option:   validation.ValidationOptions{ReadOnlyIsMandatory: false, UnassignedImmutableIsIgnored: true},
	}); err != nil {
		e := resource.CreateError(resource.InvalidSyntax, err.Error())
		statusCode, headers, body = handleError(e)
		writeResponse(rw, statusCode, headers, body)
		return
	}

	// generate default values via resource creation value defaulter
	creationValueDefaulter := worker.GetCreationValueDefaulterWorker()
	if _, err := creationValueDefaulter.Do(&worker.ValueDefaulterInput{
		Resource: subject,
		Context:  ctx,
	}); err != nil {
		e := resource.CreateError(resource.InvalidSyntax, err.Error())
		statusCode, headers, body = handleError(e)
		writeResponse(rw, statusCode, headers, body)
		return
	}

	// persistence
	persistence := worker.GetRepoUserCreateWorker()
	if _, err = persistence.Do(&worker.RepoCreateWorkerInput{
		Resource: subject,
		Context:  ctx,
	}); err != nil {
		e := resource.CreateError(resource.ServerError, err.Error())
		statusCode, headers, body = handleError(e)
		writeResponse(rw, statusCode, headers, body)
		return
	}

	// serialization and return
	serializer := worker.GetSchemaAssistedJsonSerializerWorker()
	if bodyBytes, err := serializer.Do(&worker.JsonSerializeInput{
		Target:  subject,
		Context: ctx,
	}); err != nil {
		e := resource.CreateError(resource.ServerError, err.Error())
		statusCode, headers, body = handleError(e)
		writeResponse(rw, statusCode, headers, body)
		return
	} else {
		meta := subject.Attributes["meta"].(map[string]interface{})
		statusCode = http.StatusCreated
		headers = map[string]string{
			"Content-Type": "application/json+scim",
			"ETag":         meta["version"].(string),
			"Location":     meta["location"].(string),
		}
		body = bodyBytes.([]byte)
		writeResponse(rw, statusCode, headers, body)
	}

	writeResponse(rw, 200, nil, nil)
}

func replaceUserById(rw http.ResponseWriter, req *http.Request) {

}

func updateUserById(rw http.ResponseWriter, req *http.Request) {

}

func deleteUserById(rw http.ResponseWriter, req *http.Request) {

}

func queryUser(rw http.ResponseWriter, req *http.Request) {

}

var userSchemaCache *resource.Schema

func getUserSchema() (*resource.Schema, error) {
	if nil == userSchemaCache {
		repo := persistence.GetSchemaRepository()
		if coreSchema, err := repo.Get("core", nil); err != nil {
			return nil, err
		} else if userSchema, err := repo.Get(resource.UserUrn, nil); err != nil {
			return nil, err
		} else {
			userSchemaCache = &resource.Schema{
				Schemas:    []string{resource.SchemaUrn},
				Id:         resource.UserUrn,
				Name:       "User Schema",
				Attributes: make([]*resource.Attribute, 0),
			}
			userSchemaCache.MergeWith(coreSchema.(*resource.Schema), userSchema.(*resource.Schema))
			userSchemaCache.ConstructAttributeIndex()
		}
	}
	return userSchemaCache, nil
}
