package service

import (
	p "github.com/go-scim/scimify/processor"
	"net/http"
	"sync"
)

type resourceTypeService struct {
	oneResourceTypeGetAll       sync.Once
	getAllResourceTypeProcessor p.Processor
}

func (srv *resourceTypeService) getGetAllResourceTypeProcessor() p.Processor {
	srv.oneResourceTypeGetAll.Do(func() {
		srv.getAllResourceTypeProcessor = &p.ErrorHandlingProcessor{
			Op: []p.Processor{
				p.GetWorkerBean(p.DbResourceTypeGetAll),
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
	return srv.getAllResourceTypeProcessor
}

func (srv *resourceTypeService) getAllResourceTypes(req *http.Request) (response, error) {
	processor := srv.getGetAllResourceTypeProcessor()
	ctx := &p.ProcessorContext{Request: &p.HttpRequestSource{Req: req}}
	processor.Process(ctx)
	return response{
		statusCode: ctx.ResponseStatus,
		headers:    ctx.ResponseHeaders,
		body:       ctx.ResponseBody,
	}, nil
}
