package author

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/doublehops/dh-go-framework/internal/testtools"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/doublehops/dh-go-framework/internal/httprequest"
	"github.com/doublehops/dh-go-framework/internal/model/author"
	"github.com/doublehops/dh-go-framework/internal/request"
)

func TestAuthorCRUD(t *testing.T) {
	req, _ := httprequest.GetRequester()
	ctx := context.TODO()

	payload := author.Author{
		Name: "author1",
	}

	// Test create new record.
	statusCode, res, err := req.MakeRequest(ctx, http.MethodPost, "v1/author", nil, payload)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusCreated))

	record := request.SingleItemResp{
		Data: &author.Author{},
	}
	err = json.Unmarshal(res, &record)
	assert.NoError(t, err, "unable to unmarshal record")
	d := record.Data.(*author.Author)

	assert.NoError(t, err, "error unmarshalling record")
	assert.Equal(t, payload.Name, d.Name)
	assert.Greater(t, d.ID, int32(0))
	expectedTime, duration := testtools.GetTolerance(5)
	assert.WithinDuration(t, expectedTime, *d.CreatedAt, duration)
	assert.WithinDuration(t, expectedTime, *d.UpdatedAt, duration)

	// Test get new record.
	path := fmt.Sprintf("v1/author/%d", d.ID)
	statusCode, res, err = req.MakeRequest(ctx, http.MethodGet, path, nil, nil)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusOK))

	err = json.Unmarshal(res, &record)
	assert.NoError(t, err, "unable to unmarshal record")
	d = record.Data.(*author.Author)

	assert.NoError(t, err, "error unmarshalling record")
	assert.Equal(t, payload.Name, d.Name)
	assert.Greater(t, d.ID, int32(0))
	expectedTime, duration = testtools.GetTolerance(5)
	assert.WithinDuration(t, expectedTime, *d.UpdatedAt, duration)

	// Test update new record.
	payload = author.Author{
		Name: "authorABC",
	}

	path = fmt.Sprintf("v1/author/%d", d.ID)
	statusCode, res, err = req.MakeRequest(ctx, http.MethodPut, path, nil, payload)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusOK))

	err = json.Unmarshal(res, &record)
	assert.NoError(t, err, "unable to unmarshal record")
	d = record.Data.(*author.Author)

	assert.NoError(t, err, "error unmarshalling record")
	assert.Equal(t, payload.Name, d.Name)
	assert.Greater(t, d.ID, int32(0))
	expectedTime, duration = testtools.GetTolerance(5)
	assert.WithinDuration(t, expectedTime, *d.UpdatedAt, duration)

	// Test delete new record.
	path = fmt.Sprintf("v1/author/%d", d.ID)
	statusCode, res, err = req.MakeRequest(ctx, http.MethodDelete, path, nil, nil)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusNoContent))

	// Test that record has been deleted.
	path = fmt.Sprintf("v1/author/%d", d.ID)
	statusCode, res, err = req.MakeRequest(ctx, http.MethodGet, path, nil, nil)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusNotFound))

}
