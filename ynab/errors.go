package ynab

import (
	"fmt"
	"net/http"
)

type errorData struct {
	Error APIError `json:"error"`
}

// APIError contains the error details returned by the YNAB API.
type APIError struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Detail string `json:"detail"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("(%s)%s: %s", e.ID, e.Name, e.Detail)
}

// ErrBadRequest is returned when the request is malformed or fails validation (400).
type ErrBadRequest struct {
	APIError
}

func (e ErrBadRequest) Error() string {
	return fmt.Sprintf("bad request (%s): %s", e.ID, e.Detail)
}

// ErrUnauthorized is returned when the API token is missing or invalid (401).
type ErrUnauthorized struct {
	APIError
}

func (e ErrUnauthorized) Error() string {
	return fmt.Sprintf("unauthorized (%s): %s", e.ID, e.Detail)
}

// ErrRateLimit is returned when the API rate limit has been exceeded (429).
type ErrRateLimit struct {
	APIError
}

func (e ErrRateLimit) Error() string {
	return fmt.Sprintf("rate limited (%s): %s", e.ID, e.Detail)
}

// ErrNotFound is returned when the requested resource does not exist (404).
type ErrNotFound struct {
	APIError
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("not found (%s): %s", e.ID, e.Detail)
}

// ErrForbidden is returned when the request is not permitted (403).
type ErrForbidden struct {
	APIError
}

func (e ErrForbidden) Error() string {
	return fmt.Sprintf("forbidden (%s): %s", e.ID, e.Detail)
}

// ErrServerError is returned when the YNAB API encounters an internal error (500).
type ErrServerError struct {
	APIError
}

func (e ErrServerError) Error() string {
	return fmt.Sprintf("server error (%s): %s", e.ID, e.Detail)
}

// ErrConflict is returned when a resource cannot be saved due to a conflict with an existing resource (409).
type ErrConflict struct {
	APIError
}

func (e ErrConflict) Error() string {
	return fmt.Sprintf("conflict (%s): %s", e.ID, e.Detail)
}

// ErrServiceUnavailable is returned when the YNAB API is temporarily unavailable (503).
type ErrServiceUnavailable struct {
	APIError
}

func (e ErrServiceUnavailable) Error() string {
	return fmt.Sprintf("service unavailable (%s): %s", e.ID, e.Detail)
}

func newAPIError(status int, apiErr APIError) error {
	switch status {
	case http.StatusBadRequest:
		return ErrBadRequest{apiErr}
	case http.StatusUnauthorized:
		return ErrUnauthorized{apiErr}
	case http.StatusForbidden:
		return ErrForbidden{apiErr}
	case http.StatusNotFound:
		return ErrNotFound{apiErr}
	case http.StatusConflict:
		return ErrConflict{apiErr}
	case http.StatusTooManyRequests:
		return ErrRateLimit{apiErr}
	case http.StatusInternalServerError:
		return ErrServerError{apiErr}
	case http.StatusServiceUnavailable:
		return ErrServiceUnavailable{apiErr}
	default:
		return apiErr
	}
}
