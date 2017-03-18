package service

import (
	p "github.com/go-scim/scimify/processor"
	"net/http"
	"sync"
)

type schemaService struct {
	oneGetSchema,
	oneGetAllSchema sync.Once

	getSchemaProcessor,
	getAllSchemaProcessor p.Processor
}

func (srv *schemaService) getGetSchemaProcessor() p.Processor {
	srv.oneGetSchema.Do(func() {
		srv.getSchemaProcessor = &p.ErrorHandlingProcessor{
			Op: []p.Processor{
				p.GetWorkerBean(p.ParamSchemaGet),
				p.GetWorkerBean(p.DbSchemaGet),
				p.GetWorkerBean(p.SetJsonToSingle),
				p.GetWorkerBean(p.JsonSimple),
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
	return srv.getSchemaProcessor
}

func (srv *schemaService) getGetAllSchemaProcessor() p.Processor {
	srv.oneGetAllSchema.Do(func() {
		srv.getAllSchemaProcessor = &p.ErrorHandlingProcessor{
			Op: []p.Processor{
				p.GetWorkerBean(p.DbSchemaGetAll),
				p.GetWorkerBean(p.SetJsonToMultiple),
				p.GetWorkerBean(p.JsonSimple),
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
	return srv.getAllSchemaProcessor
}

func (srv *schemaService) getAllSchemas(req *http.Request) (response, error) {
	processor := srv.getGetAllSchemaProcessor()
	ctx := &p.ProcessorContext{Request: req}
	processor.Process(ctx)
	return response{
		statusCode: ctx.ResponseStatus,
		headers:    ctx.ResponseHeaders,
		body:       ctx.ResponseBody,
	}, nil
}

func (srv *schemaService) getSchemaById(req *http.Request) (response, error) {
	processor := srv.getGetSchemaProcessor()
	ctx := &p.ProcessorContext{Request: req}
	processor.Process(ctx)
	return response{
		statusCode: ctx.ResponseStatus,
		headers:    ctx.ResponseHeaders,
		body:       ctx.ResponseBody,
	}, nil
}
