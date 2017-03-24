package processor

import (
	"fmt"
	"github.com/go-scim/scimify/persistence"
	"github.com/go-scim/scimify/resource"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

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

func (qep *parseParamForQueryEndpointProcessor) parseParam(req *http.Request, name string) string {
	switch req.Method {
	case http.MethodGet:
		return req.URL.Query().Get(name)
	case http.MethodPost:
		return req.FormValue(name)
	default:
		return ""
	}
}

func (qep *parseParamForQueryEndpointProcessor) Process(ctx *ProcessorContext) error {
	httpRequest := qep.getHttpRequest(ctx)

	if sch, err := qep.getSchema(); err != nil {
		return err
	} else {
		ctx.Schema = sch
	}

	// filter
	ctx.QueryFilter = qep.parseParam(httpRequest, "filter")
	if len(ctx.QueryFilter) == 0 {
		return resource.CreateError(resource.InvalidValue, "filter param is required.")
	}

	// sortBy
	ctx.QuerySortBy = qep.parseParam(httpRequest, "sortBy")

	// sortOrder
	switch qep.parseParam(httpRequest, "sortOrder") {
	case "", "ascending":
		ctx.QuerySortOrder = true
	case "descending":
		ctx.QuerySortOrder = false
	default:
		return resource.CreateError(resource.InvalidValue, "sortOrder param should have value [ascending] or [descending].")
	}

	// startIndex
	if v := qep.parseParam(httpRequest, "startIndex"); len(v) > 0 {
		if i, err := strconv.Atoi(v); err != nil {
			return resource.CreateError(resource.InvalidValue, "startIndex param must be a 1-based integer.")
		} else {
			if i < 1 {
				ctx.QueryPageStart = 1
			} else {
				ctx.QueryPageStart = i
			}
		}
	} else {
		ctx.QueryPageStart = 1
	}

	// count
	if v := qep.parseParam(httpRequest, "count"); len(v) > 0 {
		if i, err := strconv.Atoi(v); err != nil {
			return resource.CreateError(resource.InvalidValue, "count param must be a non-negative integer.")
		} else {
			if i < 0 {
				ctx.QueryPageSize = 0
			} else {
				ctx.QueryPageSize = i
			}
		}
	} else {
		ctx.QueryPageSize = viper.GetInt("scim.itemsPerPage")
	}

	// attributes
	ctx.Inclusion = strings.Split(qep.parseParam(httpRequest, "attributes"), ",")

	// excludedAttributes
	ctx.Exclusion = strings.Split(qep.parseParam(httpRequest, "excludedAttributes"), ",")

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
