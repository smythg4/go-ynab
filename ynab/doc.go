// Package ynab provides a Go client for the YNAB (You Need A Budget) API.
//
// # Authentication
//
// All API access requires a Personal Access Token, which can be generated
// at https://app.ynab.com/settings/developer. Pass it to NewClient:
//
//	client := ynab.NewClient(os.Getenv("YNAB_TOKEN"))
//
// # Rate Limiting
//
// The YNAB API allows 200 requests per hour. Use WithRateLimit to enforce
// this automatically. The sustained rate is reduced by the burst size so that
// burst consumption is accounted for within the hourly limit:
//
//	client := ynab.NewClient(apiKey).WithRateLimit(200, 10)
//	// allows 10 immediate requests, then throttles to 190/hr
//
// # Error Handling
//
// API errors are returned as typed errors inspectable with errors.As:
//
//	_, _, err := client.GetPlan(ctx, id, nil)
//	var notFound ynab.ErrNotFound
//	if errors.As(err, &notFound) {
//	    // handle 404
//	}
//
// Available error types: ErrBadRequest, ErrUnauthorized, ErrForbidden,
// ErrNotFound, ErrConflict, ErrRateLimit, ErrServerError, ErrServiceUnavailable.
package ynab
