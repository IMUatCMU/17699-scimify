package service

import (
	p "github.com/go-scim/scimify/processor"
	"net/http"
	"sync"
)

type rootService struct {
	oneQueryRoot       sync.Once
	queryRootProcessor p.Processor
}

func (srv *rootService) getQueryRootProcessor() p.Processor {
	srv.oneQueryRoot.Do(func() {
		srv.queryRootProcessor = &p.ErrorHandlingProcessor{
			Op: []p.Processor{
				p.GetWorkerBean(p.ParamRootQuery),
				p.GetWorkerBean(p.ParseFilter),
				p.GetWorkerBean(p.DbRootQuery),
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
	return srv.queryRootProcessor
}

func (srv *rootService) query(req *http.Request) (response, error) {
	processor := srv.getQueryRootProcessor()
	ctx := &p.ProcessorContext{Request: &p.HttpRequestSource{Req:req}}
	processor.Process(ctx)
	return response{
		statusCode: ctx.ResponseStatus,
		headers:    ctx.ResponseHeaders,
		body:       ctx.ResponseBody,
	}, nil
}
