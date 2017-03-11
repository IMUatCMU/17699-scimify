package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/resource"
	"github.com/go-scim/scimify/validation"
	"github.com/go-scim/scimify/worker"
	"github.com/go-zoo/bone"
	"io/ioutil"
	"net/http"
	"runtime"
	"sync"
)

type userService struct {
	oneUserSchema   sync.Once
	userSchemaCache *resource.Schema
}

func (srv *userService) getUserById(req *http.Request) (response, error) {
	userId := bone.GetValue(req, "userId")

	schema, err := srv.getUserSchema()
	if err != nil {
		return nil_response, resource.CreateError(resource.ServerError, "No schema was configured for user resource.")
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, resource.CK_Schema, schema)

	repoWorker := worker.GetRepoUserGetWorker()
	r, err := repoWorker.Do(&worker.RepoGetWorkerInput{Id: userId, Ctx: ctx})
	if err != nil {
		return nil_response, resource.CreateError(resource.NotFound, fmt.Sprintf("User not found by id %s", userId))
	}

	serializer := worker.GetSchemaAssistedJsonSerializerWorker()
	if bodyBytes, err := serializer.Do(&worker.JsonSerializeInput{
		Target:  r,
		Context: ctx,
	}); err != nil {
		return nil_response, resource.CreateError(resource.ServerError, err.Error())
	} else {
		meta := r.(*resource.Resource).Attributes["meta"].(map[string]interface{})
		return response{
			statusCode: http.StatusOK,
			headers: map[string]string{
				"ETag":     meta["version"].(string),
				"Location": meta["location"].(string),
			},
			body: bodyBytes.([]byte),
		}, nil
	}
}

func (srv *userService) createUser(req *http.Request) (response, error) {
	var subject *resource.Resource

	// load resource from request body
	if bodyBytes, err := ioutil.ReadAll(req.Body); err != nil {
		return nil_response, resource.CreateError(resource.ServerError, "Failed to read request body.")
	} else if subject, err = resource.NewResourceFromBytes(bodyBytes); err != nil {
		return nil_response, resource.CreateError(resource.InvalidSyntax, "The request body message was invalid or did not conform to the request schema")
	}

	// get schema guideline
	schema, err := srv.getUserSchema()
	if err != nil {
		return nil_response, resource.CreateError(resource.ServerError, "No schema was configured for user resource.")
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
		return nil_response, resource.CreateError(resource.InvalidSyntax, err.Error())
	}

	// validate the resource via resource creation validator worker
	validator := worker.GetCreationValidatorWorker()
	if _, err = validator.Do(&worker.ValidationInput{
		Resource: subject,
		Context:  ctx,
		Option:   validation.ValidationOptions{ReadOnlyIsMandatory: false, UnassignedImmutableIsIgnored: true},
	}); err != nil {
		return nil_response, resource.CreateError(resource.InvalidSyntax, err.Error())
	}

	// generate default values via resource creation value defaulter
	creationValueDefaulter := worker.GetCreationValueDefaulterWorker()
	if _, err := creationValueDefaulter.Do(&worker.ValueDefaulterInput{
		Resource: subject,
		Context:  ctx,
	}); err != nil {
		return nil_response, resource.CreateError(resource.InvalidSyntax, err.Error())
	}

	// persistence
	persistence := worker.GetRepoUserCreateWorker()
	if _, err = persistence.Do(&worker.RepoCreateWorkerInput{
		Resource: subject,
		Context:  ctx,
	}); err != nil {
		return nil_response, resource.CreateError(resource.ServerError, err.Error())
	}

	// serialization and return
	serializer := worker.GetSchemaAssistedJsonSerializerWorker()
	if bodyBytes, err := serializer.Do(&worker.JsonSerializeInput{
		Target:  subject,
		Context: ctx,
	}); err != nil {
		return nil_response, resource.CreateError(resource.ServerError, err.Error())
	} else {
		meta := subject.Attributes["meta"].(map[string]interface{})
		return response{
			statusCode: http.StatusCreated,
			headers: map[string]string{
				"ETag":     meta["version"].(string),
				"Location": meta["location"].(string),
			},
			body: bodyBytes.([]byte),
		}, nil
	}
}

func (srv *userService) replaceUserById(req *http.Request) (response, error) {
	return nil_response, nil
}

func (srv *userService) updateUserById(req *http.Request) (response, error) {
	return nil_response, nil
}

func (srv *userService) deleteUserById(req *http.Request) (response, error) {
	userId := bone.GetValue(req, "userId")

	schema, err := srv.getUserSchema()
	if err != nil {
		return nil_response, resource.CreateError(resource.ServerError, "No schema was configured for user resource.")
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, resource.CK_Schema, schema)

	repoWorker := worker.GetRepoUserDeleteWorker()
	if _, err := repoWorker.Do(&worker.RepoDeleteWorkerInput{Id: userId, Ctx: ctx}); err != nil {
		return nil_response, resource.CreateError(resource.NotFound, fmt.Sprintf("User not found by id %s", userId))
	}

	return response{statusCode: http.StatusNoContent, headers: nil, body: nil}, nil
}

func (srv *userService) queryUser(req *http.Request) (response, error) {
	p := newQueryParameters()
	err := p.parse(req)
	if nil != err {
		return nil_response, err
	}

	// obtain schema
	schema, err := srv.getUserSchema()
	if nil != err {
		return nil_response, resource.CreateError(resource.ServerError, "No schema was configured for user resource.")
	}

	// setup context
	ctx := context.Background()
	ctx = context.WithValue(ctx, resource.CK_Schema, schema)

	// parse filter
	filterWorker := worker.GetFilterWorker()
	query, err := filterWorker.Do(&worker.FilterWorkerInput{
		FilterText: p.filter,
		Schema:     schema,
	})
	if nil != err {
		return nil_response, resource.CreateError(resource.InvalidFilter, err.Error())
	}

	// run query
	repo := worker.GetRepoUserQueryWorker()
	results, err := repo.Do(&worker.RepoQueryWorkerInput{
		Context:   ctx,
		Filter:    query,
		PageStart: p.pageStart,
		PageSize:  p.pageSize,
		SortBy:    p.sortBy,
		Ascending: p.ascending,
	})
	if nil != err {
		return nil_response, resource.CreateError(resource.InvalidFilter, err.Error())
	}

	// Serialize results
	schemaSerializer := worker.GetSchemaAssistedJsonSerializerWorker()
	simpleSerializer := worker.GetDefaultJsonSerializerWorker()
	resultBytes, err := schemaSerializer.Do(&worker.JsonSerializeInput{
		Context:        ctx,
		InclusionPaths: p.attributes,
		ExclusionPaths: p.excludedAttributes,
		Target:         results,
	})
	if nil != err {
		return nil_response, resource.CreateError(resource.ServerError, err.Error())
	}
	rawJsonResults := json.RawMessage(resultBytes.([]byte))
	listResponse := resource.NewListResponse(&rawJsonResults, p.pageStart, p.pageSize, len(results.([]resource.ScimObject)))
	bytes, err := simpleSerializer.Do(&worker.JsonSerializeInput{
		Target: listResponse,
	})
	if nil != err {
		return nil_response, resource.CreateError(resource.ServerError, err.Error())
	}

	return response{
		statusCode: http.StatusOK,
		body:       bytes.([]byte),
	}, nil
}

func (srv *userService) getUserSchema() (r *resource.Schema, e error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			e = r.(error)
		}
	}()
	srv.oneUserSchema.Do(func() {
		repo := persistence.GetSchemaRepository()
		if coreSchema, err := repo.Get("core", nil); err != nil {
			panic(err)
		} else if userSchema, err := repo.Get(resource.UserUrn, nil); err != nil {
			panic(err)
		} else {
			srv.userSchemaCache = &resource.Schema{
				Schemas:    []string{resource.SchemaUrn},
				Id:         resource.UserUrn,
				Name:       "User Schema",
				Attributes: make([]*resource.Attribute, 0),
			}
			srv.userSchemaCache.MergeWith(coreSchema.(*resource.Schema), userSchema.(*resource.Schema))
			srv.userSchemaCache.ConstructAttributeIndex()
		}
	})
	r = srv.userSchemaCache
	e = nil
	return
}
