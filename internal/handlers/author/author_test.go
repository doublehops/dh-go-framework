package author

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/doublehops/dh-go-framework/internal/config"
	"github.com/doublehops/dh-go-framework/internal/httprequest"
	"github.com/doublehops/dh-go-framework/internal/model/author"
	"github.com/doublehops/dh-go-framework/internal/request"
	"github.com/doublehops/dh-go-framework/internal/testtools"
)

var cfg *config.Config

func TestMain(m *testing.M) {
	var err error

	cfg, err = config.New("./config_test.json")
	if err != nil {
		log.Printf("error starting main. %s", err.Error())
		os.Exit(1)
	}
	code := m.Run()
	os.Exit(code)
}

//nolint:funlen
func TestAuthorCRUD(t *testing.T) {
	var ok bool
	var d *author.Author
	req, _ := httprequest.GetRequester(cfg.Host.TestURL)
	ctx := context.TODO()

	payload := author.Author{
		Name: "author1",
	}

	// Test CREATE new record.
	statusCode, res, err := req.MakeRequest(ctx, http.MethodPost, "v1/author", nil, payload)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusCreated))

	record := request.SingleItemResp{
		Data: &author.Author{},
	}
	err = json.Unmarshal(res, &record)
	assert.NoError(t, err, "unable to unmarshal record")
	if d, ok = record.Data.(*author.Author); !ok {
		t.Error("unable to convert response")
	}

	assert.NoError(t, err, "error unmarshalling record")
	assert.Equal(t, payload.Name, d.Name)
	assert.Greater(t, d.ID, int32(0))
	expectedTime, duration := testtools.GetTolerance(5)
	assert.WithinDuration(t, expectedTime, *d.CreatedAt, duration)
	assert.WithinDuration(t, expectedTime, *d.UpdatedAt, duration)

	// Test GET new record.
	path := fmt.Sprintf("v1/author/%d", d.ID)
	statusCode, res, err = req.MakeRequest(ctx, http.MethodGet, path, nil, nil)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusOK))

	err = json.Unmarshal(res, &record)
	assert.NoError(t, err, "unable to unmarshal record")
	if d, ok = record.Data.(*author.Author); !ok {
		t.Error("unable to convert response")
	}

	assert.NoError(t, err, "error unmarshalling record")
	assert.Equal(t, payload.Name, d.Name)
	assert.Greater(t, d.ID, int32(0))
	expectedTime, duration = testtools.GetTolerance(5)
	assert.WithinDuration(t, expectedTime, *d.UpdatedAt, duration)

	// Test UPDATE new record.
	payload = author.Author{
		Name: "authorABC",
	}

	path = fmt.Sprintf("v1/author/%d", d.ID)
	statusCode, res, err = req.MakeRequest(ctx, http.MethodPut, path, nil, payload)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusOK))

	err = json.Unmarshal(res, &record)
	assert.NoError(t, err, "unable to unmarshal record")
	if d, ok = record.Data.(*author.Author); !ok {
		t.Error("unable to convert response")
	}

	assert.NoError(t, err, "error unmarshalling record")
	assert.Equal(t, payload.Name, d.Name)
	assert.Greater(t, d.ID, int32(0))
	expectedTime, duration = testtools.GetTolerance(5)
	assert.WithinDuration(t, expectedTime, *d.UpdatedAt, duration)

	// Test DELETE new record.
	path = fmt.Sprintf("v1/author/%d", d.ID)
	statusCode, _, err = req.MakeRequest(ctx, http.MethodDelete, path, nil, nil)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusNoContent))

	// Test that record has been deleted.
	path = fmt.Sprintf("v1/author/%d", d.ID)
	statusCode, _, err = req.MakeRequest(ctx, http.MethodGet, path, nil, nil)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusNotFound))
}
