package processor

import (
	"github.com/go-scim/scimify/resource"
	"net/http"
	"sync"
)

var (
	oneSetStatusToError,
	oneStatusOK,
	oneStatusCreated,
	oneStatusNoContent sync.Once

	setStatusToError,
	statusOk,
	statusCreated,
	statusNoContent Processor
)

func SetStatusToErrorProcessor() Processor {
	oneSetStatusToError.Do(func() {
		setStatusToError = &setStatusProcessor{useError: true}
	})
	return setStatusToError
}

func SetStatusToOKProcessor() Processor {
	oneStatusOK.Do(func() {
		statusOk = &setStatusProcessor{status: http.StatusOK}
	})
	return statusOk
}

func SetStatusToCreatedProcessor() Processor {
	oneStatusCreated.Do(func() {
		statusCreated = &setStatusProcessor{status: http.StatusCreated}
	})
	return statusCreated
}

func SetStatusToNoContentProcessor() Processor {
	oneStatusNoContent.Do(func() {
		statusNoContent = &setStatusProcessor{status: http.StatusNoContent}
	})
	return statusNoContent
}

type setStatusProcessor struct {
	useError bool
	status   int
}

func (ssp *setStatusProcessor) Process(ctx *ProcessorContext) error {
	if ssp.useError {
		ctx.ResponseStatus = ctx.Err.(resource.Error).StatusCode
	} else {
		ctx.ResponseStatus = ssp.status
	}
	return nil
}
