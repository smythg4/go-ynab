# ynab-go

A Go client for the [YNAB API](https://api.ynab.com). Requires a YNAB account and a [Personal Access Token](https://app.ynab.com/settings/developer).

## Installation

```
go get github.com/smythg4/go-ynab
```

## Usage

### Authentication

All API access requires a Personal Access Token. Pass it to `NewClient`:

```go
client := ynab.NewClient(os.Getenv("YNAB_API_KEY"))
```

### Quick Start

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/smythg4/go-ynab"
)

func main() {
    client := ynab.NewClient(os.Getenv("YNAB_API_KEY"))

    plans, err := client.GetPlans(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    for _, plan := range plans {
        fmt.Println(plan.Name)
    }
}
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

Available error types: `ErrUnauthorized`, `ErrForbidden`, `ErrNotFound`, `ErrRateLimit`, `ErrServerError`, `ErrServiceUnavailable`.

## Examples

- [ ] List plans
- [ ] Get plan month
- [ ] Get category balance
- [ ] List transactions
- [ ] Create transaction
- [ ] Create multiple transactions
- [ ] Update transaction
- [ ] Update multiple transactions
- [ ] Update category budget
- [ ] Delete transaction
- [ ] Split transaction
- [ ] Delta request

## API Coverage

### Plans
| Method | Endpoint |
|--------|----------|
| `GetPlans` | `GET /plans` |
| `GetPlan` | `GET /plans/{plan_id}` |
| `GetPlanSettings` | `GET /plans/{plan_id}/settings` |

### Accounts
| Method | Endpoint |
|--------|----------|
| `GetAccounts` | `GET /plans/{plan_id}/accounts` |
| `GetAccount` | `GET /plans/{plan_id}/accounts/{account_id}` |
| `CreateAccount` | `POST /plans/{plan_id}/accounts` |

### Categories
| Method | Endpoint |
|--------|----------|
| `GetCategories` | `GET /plans/{plan_id}/categories` |
| `GetCategory` | `GET /plans/{plan_id}/categories/{category_id}` |
| `GetCategoryForMonth` | `GET /plans/{plan_id}/months/{month}/categories/{category_id}` |
| `CreateCategory` | `POST /plans/{plan_id}/categories` |
| `CreateCategoryGroup` | `POST /plans/{plan_id}/category_groups` |
| `UpdateCategory` | `PATCH /plans/{plan_id}/categories/{category_id}` |
| `UpdateCategoryForMonth` | `PATCH /plans/{plan_id}/months/{month}/categories/{category_id}` |
| `UpdateCategoryGroup` | `PATCH /plans/{plan_id}/category_groups/{category_group_id}` |

### Months
| Method | Endpoint |
|--------|----------|
| `GetMonths` | `GET /plans/{plan_id}/months` |
| `GetMonth` | `GET /plans/{plan_id}/months/{month}` |

### Payees
| Method | Endpoint |
|--------|----------|
| `GetPayees` | `GET /plans/{plan_id}/payees` |
| `GetPayee` | `GET /plans/{plan_id}/payees/{payee_id}` |
| `GetPayeeLocations` | `GET /plans/{plan_id}/payee_locations` |
| `GetPayeeLocation` | `GET /plans/{plan_id}/payee_locations/{payee_location_id}` |
| `GetPayeeLocationsByPayee` | `GET /plans/{plan_id}/payees/{payee_id}/payee_locations` |
| `CreatePayee` | `POST /plans/{plan_id}/payees` |
| `UpdatePayee` | `PATCH /plans/{plan_id}/payees/{payee_id}` |

### Transactions
| Method | Endpoint |
|--------|----------|
| `GetTransactions` | `GET /plans/{plan_id}/transactions` |
| `GetTransaction` | `GET /plans/{plan_id}/transactions/{transaction_id}` |
| `GetTransactionsByAccount` | `GET /plans/{plan_id}/accounts/{account_id}/transactions` |
| `GetTransactionsByCategory` | `GET /plans/{plan_id}/categories/{category_id}/transactions` |
| `GetTransactionsByPayee` | `GET /plans/{plan_id}/payees/{payee_id}/transactions` |
| `GetTransactionsByMonth` | `GET /plans/{plan_id}/months/{month}/transactions` |
| `CreateTransaction` | `POST /plans/{plan_id}/transactions` |
| `CreateTransactions` | `POST /plans/{plan_id}/transactions` |
| `CreateScheduledTransaction` | `POST /plans/{plan_id}/scheduled_transactions` |
| `UpdateTransaction` | `PUT /plans/{plan_id}/transactions/{transaction_id}` |
| `UpdateTransactions` | `PATCH /plans/{plan_id}/transactions` |
| `DeleteTransaction` | `DELETE /plans/{plan_id}/transactions/{transaction_id}` |
| `ImportTransactions` | `POST /plans/{plan_id}/transactions/import` |

### Scheduled Transactions
| Method | Endpoint |
|--------|----------|
| `GetScheduledTransactions` | `GET /plans/{plan_id}/scheduled_transactions` |
| `GetScheduledTransaction` | `GET /plans/{plan_id}/scheduled_transactions/{scheduled_transaction_id}` |
| `CreateScheduledTransaction` | `POST /plans/{plan_id}/scheduled_transactions` |
| `UpdateScheduledTransaction` | `PUT /plans/{plan_id}/scheduled_transactions/{scheduled_transaction_id}` |
| `DeleteScheduledTransaction` | `DELETE /plans/{plan_id}/scheduled_transactions/{scheduled_transaction_id}` |

### Money Movements
| Method | Endpoint |
|--------|----------|
| `GetMoneyMovements` | `GET /plans/{plan_id}/money_movements` |
| `GetMoneyMovementsByMonth` | `GET /plans/{plan_id}/months/{month}/money_movements` |
| `GetMoneyMovementGroups` | `GET /plans/{plan_id}/money_movement_groups` |
| `GetMoneyMovementGroupsByMonth` | `GET /plans/{plan_id}/months/{month}/money_movement_groups` |

### User
| Method | Endpoint |
|--------|----------|
| `GetUser` | `GET /user` |

---

Not affiliated with YNAB. [YNAB API Terms of Service](https://api.ynab.com/#terms).
