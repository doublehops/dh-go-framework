package author

import (
	"context"
	"encoding/json"
	"fmt"
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

	statusCode, res, err := req.MakeRequest(ctx, http.MethodPost, "v1/author", nil, payload)
	assert.NoError(t, err, "unexpected error in request/response")
	assert.Contains(t, statusCode, fmt.Sprintf("%d", http.StatusCreated))

	record := request.SingleItemResp{
		Data: &author.Author{},
	}
	err = json.Unmarshal(res, &record)
	assert.NoError(t, err, "unable to unmaarshal record")
	d := record.Data.(*author.Author)

	assert.NoError(t, err, "error unmarshalling record")
	assert.Equal(t, payload.Name, d.Name)
}
