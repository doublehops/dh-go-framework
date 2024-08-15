package httprequest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/doublehops/dh-go-framework/internal/logga"
)

type Requester struct {
	Log              *logga.Logga
	aggregatorConfig *aggregatorConfig
}

type aggregatorConfig struct {
	Name       string     `json:"name"`
	Label      string     `json:"label"`
	HostConfig HostConfig `json:"hostConfig"`
}

type HostConfig struct {
	ApiKey  string `json:"apiKey"`
	ApiHost string `json:"apiHost"`
}

func (r *Requester) MakeRequest(ctx context.Context, method, path string, params map[string]string, payload any) (string, []byte, error) {

	client := &http.Client{}
	var p io.Reader

	msg := fmt.Sprintf("MakeRequest: %s %s", method, path)
	r.Log.Info(ctx, msg, nil)

	if payload != nil {
		pLoad, err := json.Marshal(payload)
		if err != nil {
			r.Log.Error(ctx, "Error marshalling JSON", nil)
			r.Log.Error(ctx, err.Error(), nil)
		}

		p = bytes.NewReader(pLoad)
	}

	req, err := http.NewRequest(method, r.aggregatorConfig.HostConfig.ApiHost+path, p)
	if err != nil {
		r.Log.Error(ctx, "There was an error instantiating a request", nil)
		r.Log.Error(ctx, err.Error(), nil)

		return "", nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	if params != nil {
		q := req.URL.Query()
		for key, value := range params {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	resp, err := client.Do(req)
	if err != nil {
		errMsg := fmt.Errorf("there was an error making HTTP request. %w", err)
		r.Log.Error(ctx, errMsg.Error(), nil)

		return "", nil, errMsg
	}

	statusCode := resp.Status
	respBody, _ := io.ReadAll(resp.Body)

	return statusCode, respBody, nil
}
