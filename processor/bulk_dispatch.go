package processor

import (
	"encoding/json"
	"strings"
	"net/http"
	"github.com/spf13/viper"
	"github.com/go-scim/scimify/resource"
	"sync"
)

var (
	oneBulkDispatch 	sync.Once
	bulkDispatch 	Processor
)

func BulkDispatchProcessor() Processor {
	oneBulkDispatch.Do(func() {
		bulkDispatch = &bulkDispatchProcessor{
			validationErrProcessor:NewSerialProcessor(
				GetWorkerBean(TranslateError),
				GetWorkerBean(SetJsonToError),
				GetWorkerBean(JsonSimple),
				GetWorkerBean(SetStatusToError),
			),
		}
	})
	return bulkDispatch
}

type bulkDispatchProcessor struct {
	validationErrProcessor 		Processor
}

func (p *bulkDispatchProcessor) Process(ctx *ProcessorContext) error {
	bulk := p.getBulk(ctx)

	errCount := 0
	allResponses := make([]*BulkResponseOperation, 0, len(bulk.Operations))
	for _, op := range bulk.Operations {
		if errCount > bulk.FailOnErrors {
			break
		}

		if err := op.validate(); err != nil {
			errCount++
			ctx, _ := p.createValidationErrorContext(err)
			allResponses = append(allResponses, p.createBulkResponse(op, ctx))
			continue
		}

		ctx, _ := p.createContext(op)
		p.dispatch(op, ctx)
		if ctx.ResponseStatus > 299 {
			errCount++
		}
		allResponses = append(allResponses, p.createBulkResponse(op, ctx))
	}

	bulkResponse := &BulkResponse{
		Schemas: []string{resource.BulkResponseUrn},
		Operations:allResponses,
	}
	ctx.SerializationTargetFunc = func() interface{} {
		return bulkResponse
	}

	return nil
}

func (p *bulkDispatchProcessor) createBulkResponse(op BulkRequestOperation, ctx *ProcessorContext) *BulkResponseOperation {
	br := &BulkResponseOperation{}
	br.Method = strings.ToLower(ctx.Request.Method())
	br.BulkId = op.BulkId
	if tag, ok := ctx.ResponseHeaders["ETag"]; ok {
		br.Version = tag
	}
	if loc, ok := ctx.ResponseHeaders["Location"]; ok {
		br.Location = loc
	}
	br.Status = ctx.ResponseStatus
	if ctx.ResponseStatus > 299 {
		br.Response = json.RawMessage(ctx.ResponseBody)
	} else {
		br.Response = json.RawMessage{}
	}
	return br
}

func (p *bulkDispatchProcessor) dispatch(op BulkRequestOperation, ctx *ProcessorContext) error {
	userUri := viper.GetString("scim.resourceTypeUri.user")
	groupUri := viper.GetString("scim.resourceTypeUri.group")

	var processor Processor
	switch op.Method {
	case http.MethodPost:
		switch {
		case strings.HasPrefix(op.Path, userUri):
			processor = GetWorkerBean(SrvUserCreate)
		case strings.HasPrefix(op.Path, groupUri):
			processor = GetWorkerBean(SrvGroupCreate)
		}
	case http.MethodPut:
		switch {
		case strings.HasPrefix(op.Path, userUri):
			processor = GetWorkerBean(SrvUserReplace)
		case strings.HasPrefix(op.Path, groupUri):
			processor = GetWorkerBean(SrvGroupReplace)
		}
	case http.MethodPatch:
		switch {
		case strings.HasPrefix(op.Path, userUri):
			processor = GetWorkerBean(SrvUserPatch)
		case strings.HasPrefix(op.Path, groupUri):
			processor = GetWorkerBean(SrvGroupPatch)
		}
	case http.MethodDelete:
		switch {
		case strings.HasPrefix(op.Path, userUri):
			processor = GetWorkerBean(SrvUserDelete)
		case strings.HasPrefix(op.Path, groupUri):
			processor = GetWorkerBean(SrvGroupDelete)
		}
	default:
		return resource.CreateError(resource.InvalidSyntax, "unsupported bulk operation")
	}

	err := processor.Process(ctx)
	return err
}

func (p *bulkDispatchProcessor) createContext(op BulkRequestOperation) (*ProcessorContext, error) {
	ctx := &ProcessorContext{}

	rs := &BulkRequestSource{}
	rs.populate(op)
	ctx.Request = rs

	return ctx, nil
}

func (p *bulkDispatchProcessor) createValidationErrorContext(err error) (*ProcessorContext, error) {
	ctx := &ProcessorContext{Err:err}
	e := p.validationErrProcessor.Process(ctx)
	return ctx, e
}

func (p *bulkDispatchProcessor) getBulk(ctx *ProcessorContext) *BulkRequest {
	if ctx.Bulk == nil {
		panic(&MissingContextValueError{"bulk"})
	}
	return ctx.Bulk
}

type BulkResponse struct {
	Schemas 	[]string			`json:"schemas"`
	Operations 	[]*BulkResponseOperation	`json:"Operations"`
}

type BulkResponseOperation struct {
	BulkOperation
	Location	string			`json:"location"`
	Response 	json.RawMessage		`json:"response,omitempty"`
	Status 		int 			`json:"status"`
}

type BulkRequestSource struct {
	target 		string
	method 		string
	urlParams 	map[string]string
	params 		map[string]string
	body 		[]byte
}
func (rs *BulkRequestSource) Target() string {
	return rs.target
}
func (rs *BulkRequestSource) Method() string {
	return rs.method
}
func (rs *BulkRequestSource) UrlParam(name string) string {
	return rs.urlParams[name]
}
func (rs *BulkRequestSource) Param(name string) string {
	return rs.params[name]
}
func (rs *BulkRequestSource) Body() ([]byte, error) {
	return rs.body, nil
}
func (rs *BulkRequestSource) populate(op BulkRequestOperation) {
	userUri := viper.GetString("scim.resourceTypeUri.user")
	groupUri := viper.GetString("scim.resourceTypeUri.group")

	rs.target = op.Path
	rs.method = strings.ToUpper(op.Method)

	rs.urlParams = make(map[string]string, 0)
	switch rs.method {
	case http.MethodPut, http.MethodPatch, http.MethodDelete:
		if strings.HasPrefix(rs.target, userUri + "/") {
			rs.urlParams[viper.GetString("scim.api.userIdUrlParam")] = strings.TrimPrefix(rs.target, userUri + "/")
		} else if strings.HasPrefix(rs.target, groupUri + "/") {
			rs.urlParams[viper.GetString("scim.api.groupIdUrlParam")] = strings.TrimPrefix(rs.target, groupUri + "/")
		}
	}

	rs.params = make(map[string]string, 0)

	switch rs.method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		rs.body = []byte(op.Data)
	}
}