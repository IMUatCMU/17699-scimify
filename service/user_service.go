package service

import (
	p "github.com/go-scim/scimify/processor"
	"net/http"
	"sync"
)

type userService struct {
	oneGetUser    sync.Once
	oneCreateUser sync.Once
	oneDeleteUser sync.Once
	oneQueryUser  sync.Once
	oneUpdateUser sync.Once

	getUserProcessor    p.Processor
	createUserProcessor p.Processor
	deleteUserProcessor p.Processor
	queryUserProcessor  p.Processor
	updateUserProcessor p.Processor
}

func (srv *userService) getQueryUserProcessor() p.Processor {
	srv.oneQueryUser.Do(func() {
		srv.queryUserProcessor = &p.ErrorHandlingProcessor{
			Op: []p.Processor{
				p.GetWorkerBean(p.ParamUserQuery),
				p.GetWorkerBean(p.ParseFilter),
				p.GetWorkerBean(p.DbUserQuery),
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
	return srv.queryUserProcessor
}

func (srv *userService) getGetUserProcessor() p.Processor {
	srv.oneGetUser.Do(func() {
		srv.getUserProcessor = &p.ErrorHandlingProcessor{
			Op: []p.Processor{
				p.GetWorkerBean(p.ParamUserGet),
				p.GetWorkerBean(p.DbUserGetToSingleResult),
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
	return srv.getUserProcessor
}

func (srv *userService) getCreateUserProcessor() p.Processor {
	srv.oneCreateUser.Do(func() {
		srv.createUserProcessor = &p.ErrorHandlingProcessor{
			Op: []p.Processor{
				p.GetWorkerBean(p.ParamUserCreate),
				p.GetWorkerBean(p.ValidateType),
				p.GetWorkerBean(p.ValidateRequired),
				p.GetWorkerBean(p.GenerateId),
				p.GetWorkerBean(p.GenerateUserMeta),
				p.GetWorkerBean(p.DbUserCreate),
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
	return srv.createUserProcessor
}

func (srv *userService) getUpdateUserProcessor() p.Processor {
	srv.oneUpdateUser.Do(func() {
		srv.updateUserProcessor = &p.ErrorHandlingProcessor{
			Op: []p.Processor{
				p.GetWorkerBean(p.ParamUserReplace),
				p.GetWorkerBean(p.DbUserGetToReference),
				p.GetWorkerBean(p.ValidateType),
				p.GetWorkerBean(p.ValidateRequired),
				p.GetWorkerBean(p.ValidateMutability),
				p.GetWorkerBean(p.UpdateMeta),
				p.GetWorkerBean(p.DbUserReplace),
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
	return srv.updateUserProcessor
}

func (srv *userService) getDeleteUserProcessor() p.Processor {
	srv.oneDeleteUser.Do(func() {
		srv.deleteUserProcessor = &p.ErrorHandlingProcessor{
			Op: []p.Processor{
				p.GetWorkerBean(p.ParamUserDelete),
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
	return srv.deleteUserProcessor
}

func (srv *userService) getUserById(req *http.Request) (response, error) {
	processor := srv.getGetUserProcessor()
	ctx := &p.ProcessorContext{Request: req}
	processor.Process(ctx)
	return response{
		statusCode: ctx.ResponseStatus,
		headers:    ctx.ResponseHeaders,
		body:       ctx.ResponseBody,
	}, nil
}

func (srv *userService) createUser(req *http.Request) (response, error) {
	processor := srv.getCreateUserProcessor()
	ctx := &p.ProcessorContext{Request: req}
	processor.Process(ctx)
	return response{
		statusCode: ctx.ResponseStatus,
		headers:    ctx.ResponseHeaders,
		body:       ctx.ResponseBody,
	}, nil
}

func (srv *userService) updateUserById(req *http.Request) (response, error) {
	processor := srv.getUpdateUserProcessor()
	ctx := &p.ProcessorContext{Request: req}
	processor.Process(ctx)
	return response{
		statusCode: ctx.ResponseStatus,
		headers:    ctx.ResponseHeaders,
		body:       ctx.ResponseBody,
	}, nil
}

func (srv *userService) patchUserById(req *http.Request) (response, error) {
	return nil_response, nil
}

func (srv *userService) deleteUserById(req *http.Request) (response, error) {
	processor := srv.getDeleteUserProcessor()
	ctx := &p.ProcessorContext{Request: req}
	processor.Process(ctx)
	return response{
		statusCode: ctx.ResponseStatus,
		headers:    ctx.ResponseHeaders,
		body:       ctx.ResponseBody,
	}, nil
}

func (srv *userService) queryUser(req *http.Request) (response, error) {
	processor := srv.getQueryUserProcessor()
	ctx := &p.ProcessorContext{Request: req}
	processor.Process(ctx)
	return response{
		statusCode: ctx.ResponseStatus,
		headers:    ctx.ResponseHeaders,
		body:       ctx.ResponseBody,
	}, nil
}
