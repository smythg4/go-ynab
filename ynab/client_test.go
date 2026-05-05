package ynab

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type captureTransport struct {
	req *http.Request
}

func (t *captureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.req = req
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{}`)),
		Header:     make(http.Header),
	}, nil
}

func TestNewClient(t *testing.T) {
	c := NewClient("test-token")

	if c.baseURL != "https://api.ynab.com/v1" {
		t.Errorf("got baseURL %v, want https://api.ynab.com/v1", c.baseURL)
	}

	if c.httpClient.Timeout != 10*time.Second {
		t.Errorf("got timeout %v, want 10s", c.httpClient.Timeout)
	}
}

func TestWithTimeout(t *testing.T) {
	c := NewClient("test-token").WithTimeout(30)

	if c.httpClient.Timeout != 30*time.Second {
		t.Errorf("got timeout %v, want 30s", c.httpClient.Timeout)
	}
}

func TestWithRateLimit(t *testing.T) {
	c := NewClient("test-token").WithRateLimit(200, 10)

	if c.limiter == nil {
		t.Error("expected limiter to be set, got nil")
	}
}

func TestAuthTransport(t *testing.T) {
	capture := &captureTransport{}
	transport := &authTransport{apiKey: "my-secret-token", base: capture}

	req, _ := http.NewRequest(http.MethodGet, "https://example.com", nil)
	_, _ = transport.RoundTrip(req)

	authHeader := capture.req.Header.Get("Authorization")
	want := "Bearer my-secret-token"
	if authHeader != want {
		t.Errorf("got Authorization %q, want %q", authHeader, want)
	}
}
