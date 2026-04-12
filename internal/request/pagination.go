package request

import (
	"context"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/doublehops/go-common/str"
)

var (
	defaultPage    = 1
	defaultPerPage = 10
)

// GetRequestParams - find and return pagination request vars.
func GetRequestParams(r *http.Request, filters []FilterRule, sortableFields []string) *Request {
	query := r.URL.Query()
	page, perPage, offset := getPaginationVars(query)

	f := getFilterParams(r.Context(), query, filters)

	sort := query.Get("sortBy")
	order := query.Get("order")

	// Do not include sort field if it's not in allowed list.
	if !str.SliceContains(sort, sortableFields) {
		order = ""
	}

	pg := &Request{
		Page:    page,
		PerPage: perPage,
		Offset:  offset,
		Filters: f,
		Sort:    sort,
		Order:   strings.ToUpper(order),
	}

	return pg
}

func getFilterParams(_ context.Context, query url.Values, filters []FilterRule) []FilterRule {
	var newFilters []FilterRule
	for _, f := range filters {
		val := query.Get(f.Field)

		if f.Type == FilterIsNull {
			newFilters = append(newFilters, f)

			continue
		}

		if val == "" {
			continue
		}

		f.Value = val
		newFilters = append(newFilters, f)
	}

	return newFilters
}

// getPaginationVars - Search request query for wanted var and return value, if not found, return default value.
func getPaginationVars(query map[string][]string) (int, int, int) {
	page := defaultPage
	perPage := defaultPerPage
	offset := 0

	for key, value := range query {
		queryValue := value[len(value)-1]
		switch key {
		case "perPage":
			perPage, _ = strconv.Atoi(queryValue)
		case "page":
			page, _ = strconv.Atoi(queryValue)
		}
	}

	if page != 1 {
		offset = (page - 1) * perPage
	}

	return page, perPage, offset
}

// SetRecordCount - Set pagination data to return in API response.
func (r *Request) SetRecordCount(count int32) {
	pageApprox := float64(count) / float64(r.PerPage)
	totalPages := math.Ceil(pageApprox)
	tp := int(totalPages)

	r.PageCount = tp
	r.TotalCount = count
}
