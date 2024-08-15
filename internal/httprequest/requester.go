package httprequest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/doublehops/dh-go-framework/internal/config"
	"io"
	"net/http"

	"github.com/doublehops/dh-go-framework/internal/logga"
)

type Requester struct {
	Log  *logga.Logga
	Host string
}

const (
	testHost = "http://localhost:8088/"
)

func GetRequester() (Requester, error) {
	logg := &config.Logging{
		Writer:       "stdout",
		LogLevel:     "DEBUG",
		OutputFormat: "JSON",
	}
	l, err := logga.New(logg)
	if err != nil {
		return Requester{}, err
	}

	req := Requester{
		Log:  l,
		Host: testHost,
	}

	return req, nil
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

	req, err := http.NewRequest(method, r.Host+path, p)
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
