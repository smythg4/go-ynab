package ynab

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/time/rate"
)

type authTransport struct {
	apiKey string
	base   http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.Header.Set("Authorization", "Bearer "+t.apiKey)
	return t.base.RoundTrip(req)
}

// Client is the YNAB API client. Use NewClient to create one.
type Client struct {
	baseURL    string
	httpClient *http.Client
	limiter    *rate.Limiter
	apiKey     string
}

// NewClient returns a new Client with the given Personal Access Token.
// The client has a 15-second request timeout by default.
func NewClient(apiKey string) *Client {
	return &Client{
		baseURL: "https://api.ynab.com/v1",
		httpClient: &http.Client{
			Timeout:   15 * time.Second,
			Transport: &authTransport{apiKey: apiKey, base: http.DefaultTransport},
		},
	}
}

// WithRateLimit configures a token bucket rate limiter on the client.
// The YNAB API allows 200 requests per hour. Example:
//
//	client := ynab.NewClient(apiKey).WithRateLimit(200, 10)
func (c *Client) WithRateLimit(requestsPerHour, burstVolume int) *Client {
	interval := time.Hour / time.Duration(requestsPerHour)
	c.limiter = rate.NewLimiter(rate.Every(interval), burstVolume)
	return c
}

// WithTimeout sets the HTTP client timeout in seconds. Defaults to 15 seconds.
func (c *Client) WithTimeout(timeoutSeconds int) *Client {
	c.httpClient.Timeout = time.Duration(timeoutSeconds) * time.Second
	return c
}

// WithTransport replaces the HTTP transport. Primarily useful for testing.
// Auth is preserved across transport replacements.
func (c *Client) WithTransport(t http.RoundTripper) *Client {
	c.httpClient.Transport = &authTransport{apiKey: c.apiKey, base: t}
	return c
}

// Generic method for issuing GET requests, used for endpoint logic
func (c *Client) get(ctx context.Context, endpoint string, params url.Values, out any) error {
	if c.limiter != nil {
		if err := c.limiter.Wait(ctx); err != nil {
			return err
		}
	}
	url := fmt.Sprintf("%s/%s", c.baseURL, endpoint)
	if len(params) > 0 {
		url += "?" + params.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var apiErr errorData
		if err := json.NewDecoder(res.Body).Decode(&apiErr); err != nil {
			return fmt.Errorf("request failed with status %d", res.StatusCode)
		}
		return newAPIError(res.StatusCode, apiErr.Error)
	}

	return json.NewDecoder(res.Body).Decode(out)
}

// Generic method for issuing POST requests, used for endpoint logic
func (c *Client) post(ctx context.Context, endpoint string, payload any, out any) error {
	if c.limiter != nil {
		if err := c.limiter.Wait(ctx); err != nil {
			return err
		}
	}
	url := fmt.Sprintf("%s/%s", c.baseURL, endpoint)
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		var apiErr errorData
		if err := json.NewDecoder(res.Body).Decode(&apiErr); err != nil {
			return fmt.Errorf("request failed with status %d", res.StatusCode)
		}
		return newAPIError(res.StatusCode, apiErr.Error)
	}

	return json.NewDecoder(res.Body).Decode(out)
}

// Generic method for issuing DELETE requests, used for endpoint logic
func (c *Client) delete(ctx context.Context, endpoint string, out any) error {
	if c.limiter != nil {
		if err := c.limiter.Wait(ctx); err != nil {
			return err
		}
	}
	url := fmt.Sprintf("%s/%s", c.baseURL, endpoint)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var apiErr errorData
		if err := json.NewDecoder(res.Body).Decode(&apiErr); err != nil {
			return fmt.Errorf("request failed with status %d", res.StatusCode)
		}
		return newAPIError(res.StatusCode, apiErr.Error)
	}

	return json.NewDecoder(res.Body).Decode(out)
}

// Generic method for issuing PATCH requests, used for endpoint logic
func (c *Client) patch(ctx context.Context, endpoint string, payload any, out any) error {
	if c.limiter != nil {
		if err := c.limiter.Wait(ctx); err != nil {
			return err
		}
	}
	url := fmt.Sprintf("%s/%s", c.baseURL, endpoint)
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// YNAB uses the unsupported http status code 209 for some PATCH responses
	if res.StatusCode != http.StatusOK && res.StatusCode != 209 {
		var apiErr errorData
		if err := json.NewDecoder(res.Body).Decode(&apiErr); err != nil {
			return fmt.Errorf("request failed with status %d", res.StatusCode)
		}
		return newAPIError(res.StatusCode, apiErr.Error)
	}

	return json.NewDecoder(res.Body).Decode(out)
}

// Generic method for issuing PUT requests, used for endpoint logic
func (c *Client) put(ctx context.Context, endpoint string, payload any, out any) error {
	if c.limiter != nil {
		if err := c.limiter.Wait(ctx); err != nil {
			return err
		}
	}
	url := fmt.Sprintf("%s/%s", c.baseURL, endpoint)
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var apiErr errorData
		if err := json.NewDecoder(res.Body).Decode(&apiErr); err != nil {
			return fmt.Errorf("request failed with status %d", res.StatusCode)
		}
		return newAPIError(res.StatusCode, apiErr.Error)
	}

	return json.NewDecoder(res.Body).Decode(out)
}
