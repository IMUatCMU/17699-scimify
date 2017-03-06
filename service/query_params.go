package service

import (
	"github.com/go-scim/scimify/resource"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
	"strings"
)

func newQueryParameters() *queryParameters {
	return &queryParameters{
		filter:             "",
		sortBy:             "",
		ascending:          true,
		pageStart:          1,
		pageSize:           viper.GetInt("scim.itemsPerPage"),
		attributes:         []string{},
		excludedAttributes: []string{},
	}
}

type queryParameters struct {
	filter             string
	sortBy             string
	ascending          bool
	pageStart          int
	pageSize           int
	attributes         []string
	excludedAttributes []string
}

// TODO parse query params from POST requests
func (p *queryParameters) parse(req *http.Request) error {
	p.filter = req.URL.Query().Get("filter")
	if len(p.filter) == 0 {
		return resource.CreateError(resource.InvalidValue, "[filter] is required.")
	}

	p.sortBy = req.URL.Query().Get("sortBy")
	switch req.URL.Query().Get("sortOrder") {
	case "", "ascending":
		p.ascending = true
	case "descending":
		p.ascending = false
	default:
		return resource.CreateError(resource.InvalidValue, "[sortOrder] should have value [ascending] or [descending].")
	}

	if v := req.URL.Query().Get("startIndex"); len(v) > 0 {
		if i, err := strconv.Atoi(v); err != nil {
			return resource.CreateError(resource.InvalidValue, "[startIndex] must be a 1-based integer.")
		} else {
			if i < 1 {
				p.pageStart = 1
			} else {
				p.pageStart = i
			}
		}
	}
	if v := req.URL.Query().Get("count"); len(v) > 0 {
		if i, err := strconv.Atoi(v); err != nil {
			return resource.CreateError(resource.InvalidValue, "[count] must be a non-negative integer.")
		} else {
			if i < 0 {
				p.pageSize = 0
			} else {
				p.pageSize = i
			}
		}
	}

	p.attributes = strings.Split(req.URL.Query().Get("attributes"), ",")
	p.excludedAttributes = strings.Split(req.URL.Query().Get("excludedAttributes"), ",")

	return nil
}
