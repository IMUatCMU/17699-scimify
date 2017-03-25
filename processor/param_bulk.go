package processor

import (
	"encoding/json"
	"fmt"
	"github.com/go-scim/scimify/resource"
	"github.com/spf13/viper"
	"math"
	"net/http"
	"strings"
	"sync"
)

var (
	oneBulkParamParser sync.Once
	bulkParamParser    Processor
)

func ParseParamForBulkEndpointProcessor() Processor {
	oneBulkParamParser.Do(func() {
		bulkParamParser = &parseParamForBulkEndpointProcessor{}
	})
	return bulkParamParser
}

type BulkRequest struct {
	Schemas      []string               `json:"schemas"`
	FailOnErrors int                    `json:"failOnErrors"`
	Operations   []BulkRequestOperation `json:"Operations"`
}

func (br BulkRequest) validate() error {
	if len(br.Schemas) != 1 || br.Schemas[0] != resource.BulkRequestUrn {
		return resource.CreateError(resource.InvalidSyntax, fmt.Sprintf("bulk request must have schema '%s'", resource.BulkRequestUrn))
	}
	return nil
}

type BulkOperation struct {
	Method  string `json:"method"`
	BulkId  string `json:"bulkId"`
	Version string `json:"version"`
}

type BulkRequestOperation struct {
	BulkOperation
	Path string          `json:"path"`
	Data json.RawMessage `json:"data"`
}

func (op BulkRequestOperation) validate() error {
	switch strings.ToUpper(op.Method) {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
	default:
		return resource.CreateError(
			resource.InvalidSyntax,
			fmt.Sprintf("method must be one of [%s, %s, %s, %s]",
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete),
		)
	}

	userUri := viper.GetString("scim.resourceTypeUri.user")
	groupUri := viper.GetString("scim.resourceTypeUri.group")
	if http.MethodPost == strings.ToUpper(op.Method) {
		switch op.Path {
		case userUri, groupUri:
		default:
			return resource.CreateError(
				resource.InvalidPath,
				fmt.Sprintf("path must be one of [%s, %s] when using 'post' as method", userUri, groupUri),
			)
		}
	} else {
		switch {
		case strings.HasPrefix(op.Path, userUri+"/"):
		case strings.HasPrefix(op.Path, groupUri+"/"):
		default:
			return resource.CreateError(
				resource.InvalidPath,
				fmt.Sprintf("path must be one of [%s/<id>, %s/<id>] form when using '%s' as method", userUri, groupUri, op.Method),
			)
		}
	}

	switch strings.ToUpper(op.Method) {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		if len(op.Data) == 0 {
			return resource.CreateError(
				resource.InvalidSyntax,
				fmt.Sprintf("data is required when using %s as method", op.Method),
			)
		}
	}

	return nil
}

type parseParamForBulkEndpointProcessor struct{}

func (p *parseParamForBulkEndpointProcessor) Process(ctx *ProcessorContext) error {
	req := p.getRequestSource(ctx)

	bulk, err := p.parseBulkRequest(req)
	if err != nil {
		return err
	}

	err = bulk.validate()
	if err != nil {
		return err
	}

	ctx.Bulk = bulk
	return nil
}

func (p *parseParamForBulkEndpointProcessor) parseBulkRequest(req RequestSource) (*BulkRequest, error) {
	bodyBytes, err := req.Body()
	if err != nil {
		return nil, resource.CreateError(resource.ServerError, fmt.Sprintf("failed to read request body: %s", err.Error()))
	}

	bulk := &BulkRequest{FailOnErrors: math.MaxInt64}
	err = json.Unmarshal(bodyBytes, bulk)
	if err != nil {
		return nil, resource.CreateError(resource.InvalidSyntax, fmt.Sprintf("failed to read serialize request body: %s", err.Error()))
	}

	return bulk, nil
}

func (p *parseParamForBulkEndpointProcessor) getRequestSource(ctx *ProcessorContext) RequestSource {
	if ctx.Request == nil {
		panic(&MissingContextValueError{"request source"})
	}
	return ctx.Request
}
