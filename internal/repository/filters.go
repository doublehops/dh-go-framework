package repository

import (
	"fmt"
	"log"
	"strings"
	"unicode"

	req "github.com/doublehops/dh-go-framework/internal/request"
)

type Filter struct {
	Order Order
}

type Order string

const (
	ASC  Order = "ASC"
	DESC Order = "DESC"
)

func BuildQuery(query string, p *req.Request, getCount bool) (string, []any) {
	var pParams []any

	q, params := addFilters(query, p.Filters)

	if !getCount {
		q = addSorting(q, p.Sort, p.Order)

		q, pParams = addPagination(q, p, true)
		params = append(params, pParams...)
	}

	return q, params
}

func addSorting(q, sort, order string) string {
	o := ASC

	var sq string

	if sort == "" {
		return q
	}

	s := camelToSnake(sort)
	if order != "" {
		o = Order(order)
		if o != ASC && o != DESC {
			sq = fmt.Sprintf(" ORDER BY %s %s", s, o)

			return q + sq
		}
	}

	sq = fmt.Sprintf(" ORDER BY %s %s", s, o)

	return q + sq
}

func addPagination(query string, pagination *req.Request, includePagination bool) (string, []any) {
	if !includePagination {
		log.Fatal("we should not be here")
	}

	params := []any{pagination.Offset, pagination.PerPage}

	return query + " LIMIT ?, ?", params
}

func addFilters(query string, filters []req.FilterRule) (string, req.Params) {
	var params req.Params

	var whereClauses []string
	for _, f := range filters {
		field := camelToSnake(f.Field)

		switch f.Type {
		case req.FilterEquals:
			clause := fmt.Sprintf("%s = ?", field)
			whereClauses = append(whereClauses, clause)
			params = append(params, f.Value)
		case req.FilterLike:
			clause := field + " LIKE ?"
			whereClauses = append(whereClauses, clause)
			val := fmt.Sprintf("%%%v%%", f.Value) // Will be equivalent to '%val%'
			params = append(params, val)
		case req.FilterIsNull:
			clause := field + " IS NULL"
			whereClauses = append(whereClauses, clause)
		}
	}

	str := " WHERE " + strings.Join(whereClauses, " AND ")

	return query + str, params
}

func camelToSnake(s string) string {
	var snakeCase strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			snakeCase.WriteRune('_')
		}
		snakeCase.WriteRune(unicode.ToLower(r))
	}

	return snakeCase.String()
}
