package ynab

import (
	"io"
	"net/http"
	"strings"
)

const (
	testID1 = "123e4567-e89b-12d3-a456-426614174000"
	testID2 = "223e4567-e89b-12d3-a456-426614174000"
	testID3 = "323e4567-e89b-12d3-a456-426614174000"
	testID4 = "423e4567-e89b-12d3-a456-426614174000"
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
