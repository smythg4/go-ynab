package ynab

import (
	"io"
	"net/http"
	"strings"
)

type mockTransport struct {
	body       string
	statusCode int
	lastReq    *http.Request
	lastBody   []byte
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	m.lastReq = req.Clone(req.Context())
	if req.Body != nil {
		m.lastBody, _ = io.ReadAll(req.Body)
	}
	return &http.Response{
		StatusCode: m.statusCode,
		Body:       io.NopCloser(strings.NewReader(m.body)),
		Header:     make(http.Header),
	}, nil
}

func newTestClient(body string, statusCode int) (*Client, *mockTransport) {
	transport := &mockTransport{body: body, statusCode: statusCode}
	c := NewClient("test-token")
	c.WithTransport(transport)
	return c, transport
}
