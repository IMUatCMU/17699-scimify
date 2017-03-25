package service

import (
	p "github.com/go-scim/scimify/processor"
	"net/http"
	"sync"
)

type spConfigService struct {
	oneGet               sync.Once
	getSpConfigProcessor p.Processor
}

func (srv *spConfigService) getGetSpConfigProcessor() p.Processor {
	srv.oneGet.Do(func() {
		srv.getSpConfigProcessor = &p.ErrorHandlingProcessor{
			Op: []p.Processor{
				p.GetWorkerBean(p.DbSPConfigGet),
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
	return srv.getSpConfigProcessor
}

func (srv *spConfigService) getServiceProviderConfig(req *http.Request) (response, error) {
	processor := srv.getGetSpConfigProcessor()
	ctx := &p.ProcessorContext{Request: &p.HttpRequestSource{Req:req}}
	processor.Process(ctx)
	return response{
		statusCode: ctx.ResponseStatus,
		headers:    ctx.ResponseHeaders,
		body:       ctx.ResponseBody,
	}, nil
}
