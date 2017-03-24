package service

import (
	p "github.com/go-scim/scimify/processor"
	"net/http"
	"sync"
)

type groupService struct {
	oneGetGroup    sync.Once
	oneCreateGroup sync.Once
	oneDeleteGroup sync.Once
	oneQueryGroup  sync.Once
	oneUpdateGroup sync.Once
	onePatchGroup  sync.Once

	getGroupProcessor    p.Processor
	createGroupProcessor p.Processor
	deleteGroupProcessor p.Processor
	queryGroupProcessor  p.Processor
	updateGroupProcessor p.Processor
	patchGroupProcessor  p.Processor
}

func (srv *groupService) getGetGroupProcessor() p.Processor {
	srv.oneGetGroup.Do(func() {
		srv.getGroupProcessor = &p.ErrorHandlingProcessor{
			Op: []p.Processor{
				p.GetWorkerBean(p.ParamGroupGet),
				p.GetWorkerBean(p.DbGroupGetToSingleResult),
				p.GetWorkerBean(p.SetJsonToSingle),
				p.GetWorkerBean(p.SetAllHeader),
				p.GetWorkerBean(p.JsonAssisted),
				p.GetWorkerBean(p.SetStatusToOk),
			},
			ErrOp: []p.Processor{
				p.GetWorkerBean(p.TranslateError),
				p.GetWorkerBean(p.SetJsonToError),
				p.GetWorkerBean(p.JsonSimple),
				p.GetWorkerBean(p.SetStatusToError),
			},
		}
	})
	return srv.getGroupProcessor
}

func (srv *groupService) getCreateGroupProcessor() p.Processor {
	srv.oneCreateGroup.Do(func() {
		srv.createGroupProcessor = &p.ErrorHandlingProcessor{
			Op: []p.Processor{
				p.GetWorkerBean(p.ParamGroupCreate),
				p.GetWorkerBean(p.ValidateType),
				p.GetWorkerBean(p.ValidateRequired),
				p.GetWorkerBean(p.GenerateId),
				p.GetWorkerBean(p.GenerateGroupMeta),
				p.GetWorkerBean(p.DbGroupCreate),
				p.GetWorkerBean(p.SetJsonToResource),
				p.GetWorkerBean(p.SetAllHeader),
				p.GetWorkerBean(p.JsonAssisted),
				p.GetWorkerBean(p.SetStatusToCreated),
			},
			ErrOp: []p.Processor{
				p.GetWorkerBean(p.TranslateError),
				p.GetWorkerBean(p.SetJsonToError),
				p.GetWorkerBean(p.JsonSimple),
				p.GetWorkerBean(p.SetStatusToError),
			},
		}
	})
	return srv.createGroupProcessor
}

func (srv *groupService) getDeleteGroupProcessor() p.Processor {
	srv.oneDeleteGroup.Do(func() {
		srv.deleteGroupProcessor = &p.ErrorHandlingProcessor{
			Op: []p.Processor{
				p.GetWorkerBean(p.ParamGroupDelete),
				p.GetWorkerBean(p.DbUserDelete),
				p.GetWorkerBean(p.SetStatusToNoContent),
			},
			ErrOp: []p.Processor{
				p.GetWorkerBean(p.TranslateError),
				p.GetWorkerBean(p.SetJsonToError),
				p.GetWorkerBean(p.JsonSimple),
				p.GetWorkerBean(p.SetStatusToError),
			},
		}
	})
	return srv.deleteGroupProcessor
}

func (srv *groupService) getQueryGroupProcessor() p.Processor {
	srv.oneQueryGroup.Do(func() {
		srv.queryGroupProcessor = &p.ErrorHandlingProcessor{
			Op: []p.Processor{
				p.GetWorkerBean(p.ParamGroupQuery),
				p.GetWorkerBean(p.ParseFilter),
				p.GetWorkerBean(p.DbGroupQuery),
				p.GetWorkerBean(p.SetJsonToMultiple),
				p.GetWorkerBean(p.JsonHybridList),
				p.GetWorkerBean(p.SetStatusToOk),
			},
			ErrOp: []p.Processor{
				p.GetWorkerBean(p.TranslateError),
				p.GetWorkerBean(p.SetJsonToError),
				p.GetWorkerBean(p.JsonSimple),
				p.GetWorkerBean(p.SetStatusToError),
			},
		}
	})
	return srv.queryGroupProcessor
}

func (srv *groupService) getUpdateGroupProcessor() p.Processor {
	srv.oneUpdateGroup.Do(func() {
		srv.updateGroupProcessor = &p.ErrorHandlingProcessor{
			Op: []p.Processor{
				p.GetWorkerBean(p.ParamGroupReplace),
				p.GetWorkerBean(p.DbGroupGetToReference),
				p.GetWorkerBean(p.ValidateType),
				p.GetWorkerBean(p.ValidateRequired),
				p.GetWorkerBean(p.ValidateMutability),
				p.GetWorkerBean(p.UpdateMeta),
				p.GetWorkerBean(p.DbGroupReplace),
				p.GetWorkerBean(p.SetJsonToResource),
				p.GetWorkerBean(p.SetAllHeader),
				p.GetWorkerBean(p.JsonAssisted),
				p.GetWorkerBean(p.SetStatusToOk),
			},
			ErrOp: []p.Processor{
				p.GetWorkerBean(p.TranslateError),
				p.GetWorkerBean(p.SetJsonToError),
				p.GetWorkerBean(p.JsonSimple),
				p.GetWorkerBean(p.SetStatusToError),
			},
		}
	})
	return srv.updateGroupProcessor
}

func (srv *groupService) getPatchGroupProcessor() p.Processor {
	srv.onePatchGroup.Do(func() {
		srv.patchGroupProcessor = &p.ErrorHandlingProcessor{
			Op: []p.Processor{
				p.GetWorkerBean(p.ParamGroupPatch),
				p.GetWorkerBean(p.DbGroupGetToResource),
				p.GetWorkerBean(p.DbGroupGetToReference),
				p.GetWorkerBean(p.Modification),
				p.GetWorkerBean(p.ValidateType),
				p.GetWorkerBean(p.ValidateRequired),
				p.GetWorkerBean(p.ValidateMutability),
				p.GetWorkerBean(p.UpdateMeta),
				p.GetWorkerBean(p.DbGroupReplace),
				p.GetWorkerBean(p.SetJsonToResource),
				p.GetWorkerBean(p.SetAllHeader),
				p.GetWorkerBean(p.JsonAssisted),
				p.GetWorkerBean(p.SetStatusToOk),
			},
			ErrOp: []p.Processor{
				p.GetWorkerBean(p.TranslateError),
				p.GetWorkerBean(p.SetJsonToError),
				p.GetWorkerBean(p.JsonSimple),
				p.GetWorkerBean(p.SetStatusToError),
			},
		}
	})
	return srv.patchGroupProcessor
}

func (srv *groupService) getGroupById(req *http.Request) (response, error) {
	processor := srv.getGetGroupProcessor()
	ctx := &p.ProcessorContext{Request: req}
	processor.Process(ctx)
	return response{
		statusCode: ctx.ResponseStatus,
		headers:    ctx.ResponseHeaders,
		body:       ctx.ResponseBody,
	}, nil
}

func (srv *groupService) createGroup(req *http.Request) (response, error) {
	processor := srv.getCreateGroupProcessor()
	ctx := &p.ProcessorContext{Request: req}
	processor.Process(ctx)
	return response{
		statusCode: ctx.ResponseStatus,
		headers:    ctx.ResponseHeaders,
		body:       ctx.ResponseBody,
	}, nil
}

func (srv *groupService) updateGroupById(req *http.Request) (response, error) {
	processor := srv.getUpdateGroupProcessor()
	ctx := &p.ProcessorContext{Request: req}
	processor.Process(ctx)
	return response{
		statusCode: ctx.ResponseStatus,
		headers:    ctx.ResponseHeaders,
		body:       ctx.ResponseBody,
	}, nil
}

func (srv *groupService) patchGroupById(req *http.Request) (response, error) {
	processor := srv.getPatchGroupProcessor()
	ctx := &p.ProcessorContext{Request: req}
	processor.Process(ctx)
	return response{
		statusCode: ctx.ResponseStatus,
		headers:    ctx.ResponseHeaders,
		body:       ctx.ResponseBody,
	}, nil
}

func (srv *groupService) deleteGroupById(req *http.Request) (response, error) {
	processor := srv.getDeleteGroupProcessor()
	ctx := &p.ProcessorContext{Request: req}
	processor.Process(ctx)
	return response{
		statusCode: ctx.ResponseStatus,
		headers:    ctx.ResponseHeaders,
		body:       ctx.ResponseBody,
	}, nil
}

func (srv *groupService) queryGroup(req *http.Request) (response, error) {
	processor := srv.getQueryGroupProcessor()
	ctx := &p.ProcessorContext{Request: req}
	processor.Process(ctx)
	return response{
		statusCode: ctx.ResponseStatus,
		headers:    ctx.ResponseHeaders,
		body:       ctx.ResponseBody,
	}, nil
}
