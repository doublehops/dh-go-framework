// Package request provides HTTP request parsing, filtering, pagination, and response types.
package request

// FilterType is the type of filter to apply to a query field.
type FilterType string

// FilterEquals, FilterLike, FilterIsNull, and FilterRange are the supported filter types.
var (
	FilterEquals FilterType = "FilterEquals"
	FilterLike   FilterType = "FilterLike"
	FilterIsNull FilterType = "FilterIsNull"
	FilterRange  FilterType = "FilterRange"
)

// FilterRule defines a filter to apply to a collection query.
type FilterRule struct {
	Field string
	Type  FilterType
	Value any
}

// Params is a slice of values used as SQL query parameters.
type Params []any
