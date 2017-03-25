package service

import (
	p "github.com/go-scim/scimify/processor"
	"net/http"
	"sync"
)

type bulkService struct {
	oneBulk       sync.Once
	bulkProcessor p.Processor
}

func (srv *bulkService) getBulkProcessor() p.Processor {
	srv.oneBulk.Do(func() {
		srv.bulkProcessor = &p.ErrorHandlingProcessor{
			Op: []p.Processor{
				p.GetWorkerBean(p.ParamBulk),
				p.GetServiceBean(p.BulkDispatch),
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
	return srv.bulkProcessor
}

func (srv *bulkService) doBulk(req *http.Request) (response, error) {
	processor := srv.getBulkProcessor()
	ctx := &p.ProcessorContext{Request: &p.HttpRequestSource{Req: req}}
	processor.Process(ctx)
	return response{
		statusCode: ctx.ResponseStatus,
		headers:    ctx.ResponseHeaders,
		body:       ctx.ResponseBody,
	}, nil
}
