package ynab

import (
	"context"

	"github.com/google/uuid"
)

// API is the full YNAB API surface. *Client implements this interface and is the
// recommended type for dependency injection and testing.
type API interface {
	UserService
	PlanService
	AccountService
	CategoryService
	PayeeService
	PayeeLocationService
	MonthService
	MoneyMovementService
	TransactionService
	ScheduledTransactionService
}

// UserService covers the /user endpoint.
type UserService interface {
	// GetUser returns the authenticated user.
	GetUser(ctx context.Context) (*User, error)
}

// PlanService covers the /budgets endpoints.
type PlanService interface {
	// GetPlans returns all plans for the authenticated user. Pass includeAccounts true
	// to include each plan's account list in the response.
	GetPlans(ctx context.Context, includeAccounts bool) ([]Plan, error)

	// GetPlan returns the full export for the given plan, including all sub-resources.
	// The second return value is server knowledge for delta requests.
	// For large plans this response can be substantial — consider using specific
	// resource endpoints for targeted queries.
	GetPlan(ctx context.Context, id uuid.UUID, params *ListParams) (*PlanDetails, int64, error)

	// GetLastUsedPlan returns the full export for the most recently used plan.
	// Use the returned plan's ID for subsequent sub-resource calls (GetAccounts,
	// GetTransactions, etc.) — there is no "last-used" shortcut for sub-resource endpoints.
	GetLastUsedPlan(ctx context.Context) (*PlanDetails, error)

	// GetPlanSettings returns the date and currency format settings for a plan.
	GetPlanSettings(ctx context.Context, id uuid.UUID) (*PlanSettings, error)
}

// AccountService covers the /budgets/{id}/accounts endpoints.
type AccountService interface {
	// GetAccounts returns all accounts for a plan. The second return value is server
	// knowledge for delta requests.
	GetAccounts(ctx context.Context, planID uuid.UUID, params *ListParams) ([]Account, int64, error)

	// GetAccount returns a single account by ID.
	GetAccount(ctx context.Context, planID uuid.UUID, accountID uuid.UUID) (*Account, error)

	// CreateAccount creates a new account for a plan.
	CreateAccount(ctx context.Context, planID uuid.UUID, a SaveAccount) (*Account, error)
}

// CategoryService covers the /budgets/{id}/categories endpoints.
type CategoryService interface {
	// GetCategories returns all category groups and their categories for a plan.
	// The second return value is server knowledge for delta requests.
	GetCategories(ctx context.Context, planID uuid.UUID, params *ListParams) ([]CategoryGroup, int64, error)

	// GetCategory returns a single category by ID.
	GetCategory(ctx context.Context, planID uuid.UUID, categoryID uuid.UUID) (*Category, error)

	// GetCategoryForMonth returns a category's data for a specific budget month.
	GetCategoryForMonth(ctx context.Context, planID uuid.UUID, month Date, categoryID uuid.UUID) (*Category, error)

	// CreateCategory creates a new category within a category group.
	// The second return value is server knowledge for delta requests.
	CreateCategory(ctx context.Context, planID uuid.UUID, sc SaveCategory) (*Category, int64, error)

	// CreateCategoryGroup creates a new category group.
	// The second return value is server knowledge for delta requests.
	CreateCategoryGroup(ctx context.Context, planID uuid.UUID, scg SaveCategoryGroup) (*CategoryGroup, int64, error)

	// UpdateCategory updates an existing category.
	// The second return value is server knowledge for delta requests.
	UpdateCategory(ctx context.Context, planID, categoryID uuid.UUID, sc SaveCategory) (*Category, int64, error)

	// UpdateCategoryForMonth updates a category's budgeted amount for a specific month.
	// The second return value is server knowledge for delta requests.
	UpdateCategoryForMonth(ctx context.Context, planID uuid.UUID, month Date, categoryID uuid.UUID, smc SaveMonthCategory) (*Category, int64, error)

	// UpdateCategoryGroup updates an existing category group.
	// The second return value is server knowledge for delta requests.
	UpdateCategoryGroup(ctx context.Context, planID, categoryGroupID uuid.UUID, scg SaveCategoryGroup) (*CategoryGroup, int64, error)
}

// PayeeService covers the /budgets/{id}/payees endpoints.
type PayeeService interface {
	// GetPayees returns all payees for a plan. The second return value is server
	// knowledge for delta requests.
	GetPayees(ctx context.Context, planID uuid.UUID, params *ListParams) ([]Payee, int64, error)

	// GetPayee returns a single payee by ID.
	GetPayee(ctx context.Context, planID, payeeID uuid.UUID) (*Payee, error)

	// CreatePayee creates a new payee. The second return value is server knowledge
	// for delta requests.
	CreatePayee(ctx context.Context, planID uuid.UUID, pp PostPayee) (*Payee, int64, error)

	// UpdatePayee updates an existing payee. The second return value is server knowledge
	// for delta requests.
	UpdatePayee(ctx context.Context, planID, payeeID uuid.UUID, pp PostPayee) (*Payee, int64, error)
}

// PayeeLocationService covers the /budgets/{id}/payee_locations endpoints.
type PayeeLocationService interface {
	// GetPayeeLocations returns all payee locations for a plan.
	GetPayeeLocations(ctx context.Context, planID uuid.UUID) ([]PayeeLocation, error)

	// GetPayeeLocationsByPayee returns all locations associated with a specific payee.
	GetPayeeLocationsByPayee(ctx context.Context, planID, payeeID uuid.UUID) ([]PayeeLocation, error)

	// GetPayeeLocation returns a single payee location by ID.
	GetPayeeLocation(ctx context.Context, planID, locationID uuid.UUID) (*PayeeLocation, error)
}

// MonthService covers the /budgets/{id}/months endpoints.
type MonthService interface {
	// GetMonths returns all budget months for a plan. The second return value is
	// server knowledge for delta requests.
	GetMonths(ctx context.Context, planID uuid.UUID, params *ListParams) ([]Month, int64, error)

	// GetMonth returns a single budget month including its category details.
	GetMonth(ctx context.Context, planID uuid.UUID, month Date) (*Month, error)
}

// MoneyMovementService covers the /budgets/{id}/money_movements endpoints.
type MoneyMovementService interface {
	// GetMoneyMovements returns all money movements for a plan. The second return
	// value is server knowledge for delta requests.
	GetMoneyMovements(ctx context.Context, planID uuid.UUID, params *ListParams) ([]MoneyMovement, int64, error)

	// GetMoneyMovementsByMonth returns money movements for a specific budget month.
	// The second return value is server knowledge for delta requests.
	GetMoneyMovementsByMonth(ctx context.Context, planID uuid.UUID, month Date) ([]MoneyMovement, int64, error)

	// GetMoneyMovementGroups returns all money movement groups for a plan. The second
	// return value is server knowledge for delta requests.
	GetMoneyMovementGroups(ctx context.Context, planID uuid.UUID, params *ListParams) ([]MoneyMovementGroup, int64, error)

	// GetMoneyMovementGroupsByMonth returns money movement groups for a specific budget month.
	// The second return value is server knowledge for delta requests.
	GetMoneyMovementGroupsByMonth(ctx context.Context, planID uuid.UUID, month Date) ([]MoneyMovementGroup, int64, error)
}

// TransactionService covers the /budgets/{id}/transactions endpoints.
type TransactionService interface {
	// GetTransactions returns all transactions for a plan. The second return value is
	// server knowledge for use with delta requests via TransactionListParams.LastKnowledgeOfServer.
	GetTransactions(ctx context.Context, planID uuid.UUID, params *TransactionListParams) ([]Transaction, int64, error)

	// GetTransaction returns a single transaction by ID. The second return value is
	// server knowledge for delta requests.
	GetTransaction(ctx context.Context, planID uuid.UUID, txID string) (*Transaction, int64, error)

	// GetTransactionsByAccount returns all transactions for a specific account.
	// The second return value is server knowledge for delta requests.
	GetTransactionsByAccount(ctx context.Context, planID uuid.UUID, accountID uuid.UUID, params *TransactionListParams) ([]Transaction, int64, error)

	// GetTransactionsByCategory returns all transactions for a specific category.
	// The second return value is server knowledge for delta requests.
	GetTransactionsByCategory(ctx context.Context, planID uuid.UUID, categoryID uuid.UUID, params *TransactionListParams) ([]Transaction, int64, error)

	// GetTransactionsByPayee returns all transactions for a specific payee.
	// The second return value is server knowledge for delta requests.
	GetTransactionsByPayee(ctx context.Context, planID uuid.UUID, payeeID uuid.UUID, params *TransactionListParams) ([]Transaction, int64, error)

	// GetTransactionsByMonth returns all transactions for a specific budget month.
	// The second return value is server knowledge for delta requests.
	GetTransactionsByMonth(ctx context.Context, planID uuid.UUID, month Date, params *TransactionListParams) ([]Transaction, int64, error)

	// CreateTransaction creates a single transaction.
	CreateTransaction(ctx context.Context, planID uuid.UUID, t SaveTransaction) (*CreateTransactionResponse, error)

	// CreateTransactions creates multiple transactions in a single request.
	CreateTransactions(ctx context.Context, planID uuid.UUID, t []SaveTransaction) (*CreateTransactionsResponse, error)

	// ImportTransactions triggers an import of transactions from linked accounts.
	// Returns the IDs of imported transactions.
	ImportTransactions(ctx context.Context, planID uuid.UUID) (*ImportTransactionsResponse, error)

	// UpdateTransaction replaces a transaction (PUT). Use UpdateTransactions for partial batch updates.
	UpdateTransaction(ctx context.Context, planID uuid.UUID, txID string, t UpdateTransaction) (*CreateTransactionResponse, error)

	// UpdateTransactions applies partial updates to multiple transactions (PATCH).
	UpdateTransactions(ctx context.Context, planID uuid.UUID, t []PatchTransaction) (*CreateTransactionsResponse, error)

	// DeleteTransaction deletes a transaction and returns the deleted transaction.
	// The second return value is server knowledge for delta requests.
	DeleteTransaction(ctx context.Context, planID uuid.UUID, transID string) (*Transaction, int64, error)
}

// ScheduledTransactionService covers the /budgets/{id}/scheduled_transactions endpoints.
type ScheduledTransactionService interface {
	// GetScheduledTransactions returns all scheduled transactions for a plan.
	// The second return value is server knowledge for delta requests.
	GetScheduledTransactions(ctx context.Context, planID uuid.UUID, params *ListParams) ([]ScheduledTransaction, int64, error)

	// GetScheduledTransaction returns a single scheduled transaction by ID.
	GetScheduledTransaction(ctx context.Context, planID, txID uuid.UUID) (*ScheduledTransaction, error)

	// CreateScheduledTransaction creates a new scheduled transaction.
	CreateScheduledTransaction(ctx context.Context, planID uuid.UUID, st SaveScheduledTransaction) (*ScheduledTransaction, error)

	// UpdateScheduledTransaction replaces a scheduled transaction (PUT).
	UpdateScheduledTransaction(ctx context.Context, planID uuid.UUID, txID uuid.UUID, t SaveScheduledTransaction) (*ScheduledTransaction, error)

	// DeleteScheduledTransaction deletes a scheduled transaction and returns the deleted
	// record. The second return value is server knowledge for delta requests.
	DeleteScheduledTransaction(ctx context.Context, planID uuid.UUID, transID uuid.UUID) (*ScheduledTransaction, int64, error)
}
