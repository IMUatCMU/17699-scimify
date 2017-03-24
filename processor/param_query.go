package processor

import (
	"encoding/json"
	"fmt"
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/resource"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type SearchRequest struct {
	Schemas            []string `json:"schemas"`
	Attributes         []string `json:"attributes"`
	ExcludedAttributes []string `json:"excludedAttributes"`
	Filter             string   `json:"filter"`
	SortBy             string   `json:"sortBy"`
	SortOrder          string   `json:"sortOrder"`
	StartIndex         int      `json:"startIndex"`
	Count              int      `json:"count"`
}

func (sr *SearchRequest) validate() error {
	if len(sr.Schemas) != 1 || sr.Schemas[0] != resource.SearchUrn {
		return resource.CreateError(resource.InvalidSyntax, fmt.Sprintf("search request must have urn '%s'", resource.SearchUrn))
	}

	if len(sr.Filter) == 0 {
		sr.Filter = "id pr"
	}

	if sr.StartIndex < 1 {
		sr.StartIndex = 1
	}

	if sr.Count < 0 {
		sr.Count = 0
	}

	switch sr.SortOrder {
	case "", "ascending", "descending":
	default:
		return resource.CreateError(resource.InvalidValue, "sortOrder param should have value [ascending] or [descending].")
	}

	return nil
}

func (sr *SearchRequest) copyToContext(ctx *ProcessorContext) {
	ctx.Inclusion = sr.Attributes
	ctx.Exclusion = sr.ExcludedAttributes
	ctx.QueryFilter = sr.Filter
	ctx.QuerySortBy = sr.SortBy
	switch sr.SortOrder {
	case "", "ascending":
		ctx.QuerySortOrder = true
	case "descending":
		ctx.QuerySortOrder = false
	}
	ctx.QueryPageStart = sr.StartIndex
	ctx.QueryPageSize = sr.Count
}

var (
	oneUserQueryParser,
	oneGroupQueryParser,
	oneRootQueryParser sync.Once

	userQueryParser,
	groupQueryParser,
	rootQueryParser Processor
)

func ParseParamForUserQueryEndpointProcessor() Processor {
	oneUserQueryParser.Do(func() {
		userQueryParser = &parseParamForQueryEndpointProcessor{
			internalSchemaRepo: persistence.GetInternalSchemaRepository(),
			schemaId:           viper.GetString("scim.internalSchemaId.user"),
		}
	})
	return userQueryParser
}

func ParseParamForGroupQueryEndpointProcessor() Processor {
	oneGroupQueryParser.Do(func() {
		groupQueryParser = &parseParamForQueryEndpointProcessor{
			internalSchemaRepo: persistence.GetInternalSchemaRepository(),
			schemaId:           viper.GetString("scim.internalSchemaId.group"),
		}
	})
	return groupQueryParser
}

func ParseParamForRootQueryEndpointProcessor() Processor {
	oneRootQueryParser.Do(func() {
		rootQueryParser = &parseParamForQueryEndpointProcessor{
			internalSchemaRepo: persistence.GetInternalSchemaRepository(),
			schemaId:           viper.GetString("scim.internalSchemaId.root"),
		}
	})
	return rootQueryParser
}

type parseParamForQueryEndpointProcessor struct {
	internalSchemaRepo persistence.Repository
	schemaId           string
}

type parseFunc func(*http.Request) (*SearchRequest, error)

func (qep *parseParamForQueryEndpointProcessor) parseParamsFromHttpGet(httpRequest *http.Request) (*SearchRequest, error) {
	sr := &SearchRequest{
		Schemas:    []string{resource.SearchUrn},
		StartIndex: 1,
		Count:      viper.GetInt("scim.itemsPerPage"),
	}

	sr.Attributes = strings.Split(httpRequest.URL.Query().Get("attributes"), ",")
	sr.ExcludedAttributes = strings.Split(httpRequest.URL.Query().Get("excludedAttributes"), ",")
	sr.Filter = httpRequest.URL.Query().Get("filter")
	sr.SortBy = httpRequest.URL.Query().Get("sortBy")
	sr.SortOrder = httpRequest.URL.Query().Get("sortOrder")
	if v := httpRequest.URL.Query().Get("startIndex"); len(v) > 0 {
		if i, err := strconv.Atoi(v); err != nil {
			return nil, resource.CreateError(resource.InvalidValue, "startIndex param must be a 1-based integer.")
		} else {
			if i < 1 {
				sr.StartIndex = 1
			} else {
				sr.StartIndex = i
			}
		}
	} else {
		sr.StartIndex = 1
	}
	if v := httpRequest.URL.Query().Get("count"); len(v) > 0 {
		if i, err := strconv.Atoi(v); err != nil {
			return nil, resource.CreateError(resource.InvalidValue, "count param must be a non-negative integer.")
		} else {
			if i < 0 {
				sr.Count = 0
			} else {
				sr.Count = i
			}
		}
	} else {
		sr.Count = viper.GetInt("scim.itemsPerPage")
	}

	return sr, nil
}

func (qep *parseParamForQueryEndpointProcessor) parseParamsFromHttpPost(httpRequest *http.Request) (*SearchRequest, error) {
	sr := &SearchRequest{
		StartIndex: 1,
		Count:      viper.GetInt("scim.itemsPerPage"),
	}
	bodyBytes, err := ioutil.ReadAll(httpRequest.Body)
	if err != nil {
		return nil, resource.CreateError(resource.ServerError, fmt.Sprintf("failed to read request body: %s", err.Error()))
	}

	err = json.Unmarshal(bodyBytes, sr)
	if err != nil {
		return nil, resource.CreateError(resource.InvalidSyntax, fmt.Sprintf("failed to deserialize request body: %s", err.Error()))
	}

	return sr, nil
}

func (qep *parseParamForQueryEndpointProcessor) Process(ctx *ProcessorContext) error {
	httpRequest := qep.getHttpRequest(ctx)

	if sch, err := qep.getSchema(); err != nil {
		return err
	} else {
		ctx.Schema = sch
	}

	var f parseFunc
	switch httpRequest.Method {
	case http.MethodGet:
		f = qep.parseParamsFromHttpGet
	case http.MethodPost:
		f = qep.parseParamsFromHttpPost
	default:
		return resource.CreateError(resource.NotImplemented, fmt.Sprintf("resource query by http method %s is not supported.", httpRequest.Method))
	}

	sr, err := f(httpRequest)
	if err != nil {
		return err
	}

	err = sr.validate()
	if err != nil {
		return err
	}

	sr.copyToContext(ctx)

	return nil
}

func (qep *parseParamForQueryEndpointProcessor) getSchema() (*resource.Schema, error) {
	obj, err := qep.internalSchemaRepo.Get(qep.schemaId)
	if err != nil {
		return nil, resource.CreateError(resource.ServerError, fmt.Sprintf("failed to get schema for resource query: %s", err.Error()))
	} else {
		return obj.(*resource.Schema), nil
	}
}

func (qep *parseParamForQueryEndpointProcessor) getHttpRequest(ctx *ProcessorContext) *http.Request {
	if ctx.Request == nil {
		panic(&MissingContextValueError{"http request"})
	}
	return ctx.Request
}
