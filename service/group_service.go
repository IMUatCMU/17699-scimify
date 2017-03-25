package service

import (
	p "github.com/go-scim/scimify/processor"
	"net/http"
)

type groupService struct {}

func (srv *groupService) getGroupById(req *http.Request) (response, error) {
	processor := p.GetWorkerBean(p.SrvGroupGet)
	ctx := &p.ProcessorContext{Request: &p.HttpRequestSource{Req:req}}
	processor.Process(ctx)
	return response{
		statusCode: ctx.ResponseStatus,
		headers:    ctx.ResponseHeaders,
		body:       ctx.ResponseBody,
	}, nil
}

func (srv *groupService) createGroup(req *http.Request) (response, error) {
	processor := p.GetWorkerBean(p.SrvGroupCreate)
	ctx := &p.ProcessorContext{Request: &p.HttpRequestSource{Req:req}}
	processor.Process(ctx)
	return response{
		statusCode: ctx.ResponseStatus,
		headers:    ctx.ResponseHeaders,
		body:       ctx.ResponseBody,
	}, nil
}

func (srv *groupService) updateGroupById(req *http.Request) (response, error) {
	processor := p.GetWorkerBean(p.SrvGroupReplace)
	ctx := &p.ProcessorContext{Request: &p.HttpRequestSource{Req:req}}
	processor.Process(ctx)
	return response{
		statusCode: ctx.ResponseStatus,
		headers:    ctx.ResponseHeaders,
		body:       ctx.ResponseBody,
	}, nil
}

func (srv *groupService) patchGroupById(req *http.Request) (response, error) {
	processor := p.GetWorkerBean(p.SrvGroupPatch)
	ctx := &p.ProcessorContext{Request: &p.HttpRequestSource{Req:req}}
	processor.Process(ctx)
	return response{
		statusCode: ctx.ResponseStatus,
		headers:    ctx.ResponseHeaders,
		body:       ctx.ResponseBody,
	}, nil
}

func (srv *groupService) deleteGroupById(req *http.Request) (response, error) {
	processor := p.GetWorkerBean(p.SrvGroupDelete)
	ctx := &p.ProcessorContext{Request: &p.HttpRequestSource{Req:req}}
	processor.Process(ctx)
	return response{
		statusCode: ctx.ResponseStatus,
		headers:    ctx.ResponseHeaders,
		body:       ctx.ResponseBody,
	}, nil
}

func (srv *groupService) queryGroup(req *http.Request) (response, error) {
	processor := p.GetWorkerBean(p.SrvGroupQuery)
	ctx := &p.ProcessorContext{Request: &p.HttpRequestSource{Req:req}}
	processor.Process(ctx)
	return response{
		statusCode: ctx.ResponseStatus,
		headers:    ctx.ResponseHeaders,
		body:       ctx.ResponseBody,
	}, nil
}
