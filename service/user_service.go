package service

import (
	p "github.com/go-scim/scimify/processor"
	"net/http"
)

type userService struct {}

func (srv *userService) getUserById(req *http.Request) (response, error) {
	processor := p.GetWorkerBean(p.SrvUserGet)
	ctx := &p.ProcessorContext{Request: &p.HttpRequestSource{Req:req}}
	processor.Process(ctx)
	return response{
		statusCode: ctx.ResponseStatus,
		headers:    ctx.ResponseHeaders,
		body:       ctx.ResponseBody,
	}, nil
}

func (srv *userService) createUser(req *http.Request) (response, error) {
	processor := p.GetWorkerBean(p.SrvUserCreate)
	ctx := &p.ProcessorContext{Request: &p.HttpRequestSource{Req:req}}
	processor.Process(ctx)
	return response{
		statusCode: ctx.ResponseStatus,
		headers:    ctx.ResponseHeaders,
		body:       ctx.ResponseBody,
	}, nil
}

func (srv *userService) updateUserById(req *http.Request) (response, error) {
	processor := p.GetWorkerBean(p.SrvUserReplace)
	ctx := &p.ProcessorContext{Request: &p.HttpRequestSource{Req:req}}
	processor.Process(ctx)
	return response{
		statusCode: ctx.ResponseStatus,
		headers:    ctx.ResponseHeaders,
		body:       ctx.ResponseBody,
	}, nil
}

func (srv *userService) patchUserById(req *http.Request) (response, error) {
	processor := p.GetWorkerBean(p.SrvUserPatch)
	ctx := &p.ProcessorContext{Request: &p.HttpRequestSource{Req:req}}
	processor.Process(ctx)
	return response{
		statusCode: ctx.ResponseStatus,
		headers:    ctx.ResponseHeaders,
		body:       ctx.ResponseBody,
	}, nil
}

func (srv *userService) deleteUserById(req *http.Request) (response, error) {
	processor := p.GetWorkerBean(p.SrvUserDelete)
	ctx := &p.ProcessorContext{Request: &p.HttpRequestSource{Req:req}}
	processor.Process(ctx)
	return response{
		statusCode: ctx.ResponseStatus,
		headers:    ctx.ResponseHeaders,
		body:       ctx.ResponseBody,
	}, nil
}

func (srv *userService) queryUser(req *http.Request) (response, error) {
	processor := p.GetWorkerBean(p.SrvUserQuery)
	ctx := &p.ProcessorContext{Request: &p.HttpRequestSource{Req:req}}
	processor.Process(ctx)
	return response{
		statusCode: ctx.ResponseStatus,
		headers:    ctx.ResponseHeaders,
		body:       ctx.ResponseBody,
	}, nil
}
