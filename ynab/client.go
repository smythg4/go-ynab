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

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	limiter    *rate.Limiter
}

func NewClient(apiKey string) *Client {
	return &Client{
		baseURL: "https://api.ynab.com/v1",
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *Client) WithRateLimit(requestsPerHour, burstVolume int) *Client {
	interval := time.Hour / time.Duration(requestsPerHour)
	c.limiter = rate.NewLimiter(rate.Every(interval), burstVolume)
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
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var apiErr ErrorData
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
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		var apiErr ErrorData
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
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var apiErr ErrorData
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
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != 209 {
		var apiErr ErrorData
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
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var apiErr ErrorData
		if err := json.NewDecoder(res.Body).Decode(&apiErr); err != nil {
			return fmt.Errorf("request failed with status %d", res.StatusCode)
		}
		return newAPIError(res.StatusCode, apiErr.Error)
	}

	return json.NewDecoder(res.Body).Decode(out)
}
