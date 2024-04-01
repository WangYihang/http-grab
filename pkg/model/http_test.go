package model_test

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/WangYihang/http-grab/pkg/model"
)

func TestHTTPRequestBody(t *testing.T) {
	testcases := []struct {
		method string
		scheme string
		path   string
		host   string
		body   string
	}{
		{
			method: "POST",
			scheme: "http",
			path:   "/",
			host:   "localhost:8080",
			body:   "Hello, World!",
		},
	}
	for _, tc := range testcases {
		bodyReader := bytes.NewReader([]byte(tc.body))
		req, err := http.NewRequest(
			tc.method,
			fmt.Sprintf(
				"%s://%s%s",
				tc.scheme,
				tc.host,
				tc.path,
			),
			bodyReader,
		)
		if err != nil {
			t.Errorf("failed to create request: %v", err)
		}
		httpRequest, err := model.NewHTTPRequest(req)
		if err != nil {
			t.Errorf("failed to create http request: %v", err)
		}
		fmt.Println(httpRequest.Body, tc.body)
		if httpRequest.Body != tc.body {
			t.Errorf("expected body %s, got %s", tc.body, httpRequest.Body)
		}
	}
}
