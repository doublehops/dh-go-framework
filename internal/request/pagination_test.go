package request

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRequestParamsSorting(t *testing.T) {
	sortable := []string{"id", "createdAt"}

	tests := []struct {
		name          string
		sortBy        string
		order         string
		expectedSort  string
		expectedOrder string
	}{
		{name: "allowed field and order", sortBy: "createdAt", order: "desc", expectedSort: "createdAt", expectedOrder: "DESC"},
		{name: "allowed field without order", sortBy: "id", expectedSort: "id", expectedOrder: ""},
		{name: "field not in allowed list is dropped", sortBy: "password", order: "asc", expectedSort: "", expectedOrder: ""},
		{name: "injected sort field is dropped", sortBy: "id; DROP TABLE user", expectedSort: "", expectedOrder: ""},
		{name: "injected order is dropped", sortBy: "id", order: "ASC, (SELECT SLEEP(5))", expectedSort: "id", expectedOrder: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := url.Values{}
			q.Set("sortBy", tt.sortBy)
			q.Set("order", tt.order)

			r := httptest.NewRequest(http.MethodGet, "/v1/author?"+q.Encode(), http.NoBody)

			p := GetRequestParams(r, nil, sortable)
			assert.Equal(t, tt.expectedSort, p.Sort)
			assert.Equal(t, tt.expectedOrder, p.Order)
		})
	}
}
