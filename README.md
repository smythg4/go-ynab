# ynab-go
[![Go Report Card](https://goreportcard.com/badge/github.com/smythg4/go-ynab)](https://goreportcard.com/report/github.com/smythg4/go-ynab)
[![GoDoc](https://pkg.go.dev/badge/github.com/smythg4/go-ynab/ynab.svg)](https://pkg.go.dev/github.com/smythg4/go-ynab/ynab)
  
A Go client for the [YNAB API](https://api.ynab.com). Supports full access  to all published YNAB API endpoints. Requires a YNAB account and a [Personal Access Token](https://app.ynab.com/settings/developer).

## Installation

```
go get github.com/smythg4/go-ynab
```

## Usage

### Authentication

All API access requires a Personal Access Token. Pass it to `NewClient`:

```go
client := ynab.NewClient(os.Getenv("YNAB_TOKEN"))
```

### Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/smythg4/go-ynab/ynab"
)

func main() {
	client := ynab.NewClient(os.Getenv("YNAB_TOKEN"))

	plans, err := client.GetPlans(context.Background(), true) // include account information in the return
	if err != nil {
		log.Fatal(err)
	}

	for _, plan := range plans {
		fmt.Println(plan.Name)
		for _, acct := range plan.Accounts {
			fmt.Printf("   %s\n", acct.Name)
		}
	}
}
```

## Context

All methods accept a `context.Context` as their first argument, enabling callers to cancel in-flight requests, enforce deadlines, and propagate request-scoped values through the call chain.

## Rate Limiting

The [YNAB API](https://api.ynab.com/#rate-limiting) allows 200 requests per hour. `WithRateLimit` enables a token bucket limiter that automatically spaces requests to stay within that limit:

```go
client := ynab.NewClient(os.Getenv("YNAB_TOKEN")).WithRateLimit(200, 10)
```

The first argument is the request budget per hour; the second is the burst size — the number of requests that can be made immediately before throttling begins. To keep total consumption within YNAB's limit, the sustained rate is reduced by the burst size: `WithRateLimit(200, 10)` allows 10 immediate requests, then throttles to 190 per hour. Calls block until a token is available rather than returning an error, so no retry logic is needed on the caller's side.

Rate limiting is opt-in. Omit `WithRateLimit` for scripts or one-off tools where request volume is not a concern.

## Timeout

The default request timeout is 10 seconds. Use `WithTimeout` to override it:

```go
client := ynab.NewClient(os.Getenv("YNAB_TOKEN")).WithTimeout(30)
```

Both methods return the client, so they can be chained:

```go
client := ynab.NewClient(os.Getenv("YNAB_TOKEN")).
    WithRateLimit(200, 10).
    WithTimeout(30)
```

## Error Handling

Errors from the API are returned as typed errors that can be inspected with `errors.As`:

```go
plan, err := client.GetPlan(ctx, id)
if err != nil {
    var notFound ynab.ErrNotFound
    if errors.As(err, &notFound) {
        // handle missing plan
    }
    return err
}
```

Available error types: `ErrBadRequest`, `ErrUnauthorized`, `ErrForbidden`, `ErrNotFound`, `ErrConflict`, `ErrRateLimit`, `ErrServerError`, `ErrServiceUnavailable`.

## Examples

- [List plans](examples/list-plans/main.go)
- [Get plan month](examples/get-plan-month/main.go)
- [Get category balance](examples/get-category-balance/main.go)
- [List transactions](examples/list-transactions/main.go)
- [Create transaction](examples/create-transaction/main.go)
- [Create multiple transactions](examples/create-transactions/main.go)
- [Update transaction](examples/update-transaction/main.go)
- [Update multiple transactions](examples/update-transactions/main.go)
- [Update category budget](examples/update-category-budget/main.go)
- [Delete transaction](examples/delete-transaction/main.go)
- [Split transaction](examples/split-transactions/main.go)
- [Delta request](examples/delta-request/main.go)

## API Coverage

### Plans
| Method | Endpoint |
|--------|----------|
| `GetPlans` | `GET /plans` |
| `GetPlan` † | `GET /plans/{plan_id}` |
| `GetLastUsedPlan` | `GET /plans/last-used` |
| `GetPlanSettings` | `GET /plans/{plan_id}/settings` |

### Accounts
| Method | Endpoint |
|--------|----------|
| `GetAccounts` † | `GET /plans/{plan_id}/accounts` |
| `GetAccount` | `GET /plans/{plan_id}/accounts/{account_id}` |
| `CreateAccount` | `POST /plans/{plan_id}/accounts` |

### Categories
| Method | Endpoint |
|--------|----------|
| `GetCategories` † | `GET /plans/{plan_id}/categories` |
| `GetCategory` | `GET /plans/{plan_id}/categories/{category_id}` |
| `GetCategoryForMonth` | `GET /plans/{plan_id}/months/{month}/categories/{category_id}` |
| `CreateCategory` † | `POST /plans/{plan_id}/categories` |
| `CreateCategoryGroup` † | `POST /plans/{plan_id}/category_groups` |
| `UpdateCategory` † | `PATCH /plans/{plan_id}/categories/{category_id}` |
| `UpdateCategoryForMonth` † | `PATCH /plans/{plan_id}/months/{month}/categories/{category_id}` |
| `UpdateCategoryGroup` † | `PATCH /plans/{plan_id}/category_groups/{category_group_id}` |

### Months
| Method | Endpoint |
|--------|----------|
| `GetMonths` † | `GET /plans/{plan_id}/months` |
| `GetMonth` | `GET /plans/{plan_id}/months/{month}` |

### Payees
| Method | Endpoint |
|--------|----------|
| `GetPayees` † | `GET /plans/{plan_id}/payees` |
| `GetPayee` | `GET /plans/{plan_id}/payees/{payee_id}` |
| `GetPayeeLocations` | `GET /plans/{plan_id}/payee_locations` |
| `GetPayeeLocation` | `GET /plans/{plan_id}/payee_locations/{payee_location_id}` |
| `GetPayeeLocationsByPayee` | `GET /plans/{plan_id}/payees/{payee_id}/payee_locations` |
| `CreatePayee` † | `POST /plans/{plan_id}/payees` |
| `UpdatePayee` † | `PATCH /plans/{plan_id}/payees/{payee_id}` |

### Transactions
| Method | Endpoint |
|--------|----------|
| `GetTransactions` † | `GET /plans/{plan_id}/transactions` |
| `GetTransaction` † | `GET /plans/{plan_id}/transactions/{transaction_id}` |
| `GetTransactionsByAccount` † | `GET /plans/{plan_id}/accounts/{account_id}/transactions` |
| `GetTransactionsByCategory` † | `GET /plans/{plan_id}/categories/{category_id}/transactions` |
| `GetTransactionsByPayee` † | `GET /plans/{plan_id}/payees/{payee_id}/transactions` |
| `GetTransactionsByMonth` † | `GET /plans/{plan_id}/months/{month}/transactions` |
| `CreateTransaction` | `POST /plans/{plan_id}/transactions` |
| `CreateTransactions` | `POST /plans/{plan_id}/transactions` |
| `UpdateTransaction` | `PUT /plans/{plan_id}/transactions/{transaction_id}` |
| `UpdateTransactions` | `PATCH /plans/{plan_id}/transactions` |
| `DeleteTransaction` † | `DELETE /plans/{plan_id}/transactions/{transaction_id}` |
| `ImportTransactions` | `POST /plans/{plan_id}/transactions/import` |

### Scheduled Transactions
| Method | Endpoint |
|--------|----------|
| `GetScheduledTransactions` † | `GET /plans/{plan_id}/scheduled_transactions` |
| `GetScheduledTransaction` | `GET /plans/{plan_id}/scheduled_transactions/{scheduled_transaction_id}` |
| `CreateScheduledTransaction` | `POST /plans/{plan_id}/scheduled_transactions` |
| `UpdateScheduledTransaction` | `PUT /plans/{plan_id}/scheduled_transactions/{scheduled_transaction_id}` |
| `DeleteScheduledTransaction` | `DELETE /plans/{plan_id}/scheduled_transactions/{scheduled_transaction_id}` |

### Money Movements
| Method | Endpoint |
|--------|----------|
| `GetMoneyMovements` † | `GET /plans/{plan_id}/money_movements` |
| `GetMoneyMovementsByMonth` † | `GET /plans/{plan_id}/months/{month}/money_movements` |
| `GetMoneyMovementGroups` † | `GET /plans/{plan_id}/money_movement_groups` |
| `GetMoneyMovementGroupsByMonth` † | `GET /plans/{plan_id}/months/{month}/money_movement_groups` |

### User
| Method | Endpoint |
|--------|----------|
| `GetUser` | `GET /user` |

† Returns server knowledge as a second return value for use with delta requests.

## Test Coverage

Unit tests cover all endpoints (GET, POST, PATCH, PUT, DELETE), client configuration, error type dispatch, and auth header injection. Write operation tests verify the HTTP method and request body serialization.

```
go test ./ynab/...
```

Integration tests exercise the live API against a real plan and require `YNAB_TOKEN` and `YNAB_TEST_PLAN_ID` environment variables. They are opt-in via a build tag:

```
YNAB_TOKEN=... YNAB_TEST_PLAN_ID=... go test -tags integration -v ./integration/
```

## License

[MIT](LICENSE)

---

Not affiliated with YNAB. [YNAB API Terms of Service](https://api.ynab.com/#terms).
