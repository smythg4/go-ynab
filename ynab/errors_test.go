package ynab

import (
	"errors"
	"net/http"
	"testing"
)

func TestNewAPIError(t *testing.T) {
	apiErr := APIError{ID: "test", Name: "test_error", Detail: "something went wrong"}

	tests := []struct {
		name       string
		status     int
		target     any
		errChecker func(error) bool
	}{
		{
			name:   "401 returns ErrUnauthorized",
			status: http.StatusUnauthorized,
			errChecker: func(err error) bool {
				var target ErrUnauthorized
				return errors.As(err, &target)
			},
		},
		{
			name:   "403 returns ErrForbidden",
			status: http.StatusForbidden,
			errChecker: func(err error) bool {
				var target ErrForbidden
				return errors.As(err, &target)
			},
		},
		{
			name:   "404 returns ErrNotFound",
			status: http.StatusNotFound,
			errChecker: func(err error) bool {
				var target ErrNotFound
				return errors.As(err, &target)
			},
		},
		{
			name:   "429 returns ErrRateLimit",
			status: http.StatusTooManyRequests,
			errChecker: func(err error) bool {
				var target ErrRateLimit
				return errors.As(err, &target)
			},
		},
		{
			name:   "500 returns ErrServerError",
			status: http.StatusInternalServerError,
			errChecker: func(err error) bool {
				var target ErrServerError
				return errors.As(err, &target)
			},
		},
		{
			name:   "503 returns ErrServiceUnavailable",
			status: http.StatusServiceUnavailable,
			errChecker: func(err error) bool {
				var target ErrServiceUnavailable
				return errors.As(err, &target)
			},
		},
		{
			name:   "unknown status returns APIError",
			status: 418,
			errChecker: func(err error) bool {
				var target APIError
				return errors.As(err, &target)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := newAPIError(tt.status, apiErr)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.errChecker(err) {
				t.Errorf("wrong error type for status %d: got %T", tt.status, err)
			}
		})
	}
}
