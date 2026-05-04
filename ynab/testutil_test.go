package ynab

import (
	"io"
	"net/http"
	"strings"
)

type mockTransport struct {
	body       string
	statusCode int
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: m.statusCode,
		Body:       io.NopCloser(strings.NewReader(m.body)),
		Header:     make(http.Header),
	}, nil
}

func newTestClient(body string, statusCode int) *Client {
	c := NewClient("test-token")
	c.WithTransport(&mockTransport{body: body, statusCode: statusCode})
	return c
}
