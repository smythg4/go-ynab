package ynab

import (
	"fmt"
	"net/http"
)

type ErrorData struct {
	Error APIError `json:"error"`
}

type APIError struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Detail string `json:"detail"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("(%s)%s: %s", e.ID, e.Name, e.Detail)
}

type ErrUnauthorized struct {
	APIError
}

func (e ErrUnauthorized) Error() string {
	return fmt.Sprintf("unauthorized (%s): %s", e.ID, e.Detail)
}

type ErrRateLimit struct {
	APIError
}

func (e ErrRateLimit) Error() string {
	return fmt.Sprintf("rate limited (%s): %s", e.ID, e.Detail)
}

type ErrNotFound struct {
	APIError
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("not found (%s): %s", e.ID, e.Detail)
}

type ErrForbidden struct {
	APIError
}

func (e ErrForbidden) Error() string {
	return fmt.Sprintf("forbidden (%s): %s", e.ID, e.Detail)
}

type ErrServerError struct {
	APIError
}

func (e ErrServerError) Error() string {
	return fmt.Sprintf("server error (%s): %s", e.ID, e.Detail)
}

type ErrServiceUnavailable struct {
	APIError
}

func (e ErrServiceUnavailable) Error() string {
	return fmt.Sprintf("service unavailable (%s): %s", e.ID, e.Detail)
}

func newAPIError(status int, apiErr APIError) error {
	switch status {
	case http.StatusUnauthorized:
		return ErrUnauthorized{apiErr}
	case http.StatusTooManyRequests:
		return ErrRateLimit{apiErr}
	case http.StatusNotFound:
		return ErrNotFound{apiErr}
	case http.StatusForbidden:
		return ErrForbidden{apiErr}
	case http.StatusInternalServerError:
		return ErrServerError{apiErr}
	case http.StatusServiceUnavailable:
		return ErrServiceUnavailable{apiErr}
	default:
		return apiErr
	}
}
