package ynab

import (
	"context"
	"fmt"
	"net/url"

	"github.com/google/uuid"
)

type transactionData struct {
	Data struct {
		Transaction     Transaction `json:"transaction"`
		ServerKnowledge int64       `json:"server_knowledge"`
	} `json:"data"`
}

type transactionsData struct {
	Data struct {
		Transactions    []Transaction `json:"transactions"`
		ServerKnowledge int64         `json:"server_knowledge"`
	} `json:"data"`
}

// Transaction represents a single YNAB transaction. Amounts are in milliunits (divide by 1000 for display).
//
// ID is a string rather than uuid.UUID because upcoming scheduled transaction instances use a
// compound format "{scheduled_uuid}_{date}" (e.g. "abc123..._2025-06-01"). Regular posted
// transactions have standard UUID ids.
type Transaction struct {
	ID                      string           `json:"id"`
	Date                    Date             `json:"date"`
	Amount                  int64            `json:"amount"`
	Memo                    *string          `json:"memo"`
	Cleared                 ClearedStatus    `json:"cleared"`
	Approved                bool             `json:"approved"`
	FlagColor               *FlagColor       `json:"flag_color"`
	FlagName                *string          `json:"flag_name"`
	AccountID               uuid.UUID        `json:"account_id"`
	PayeeID                 *uuid.UUID       `json:"payee_id"`
	AccountName             string           `json:"account_name"`
	PayeeName               *string          `json:"payee_name"`
	CategoryID              *uuid.UUID       `json:"category_id"`
	CategoryName            *string          `json:"category_name"`
	MatchedTransactionID    *uuid.UUID       `json:"matched_transaction_id"`
	TransferAccountID       *uuid.UUID       `json:"transfer_account_id"`
	TransferTransactionID   *uuid.UUID       `json:"transfer_transaction_id"`
	ImportID                *string          `json:"import_id"`
	ImportPayeeName         *string          `json:"import_payee_name"`
	ImportPayeeNameOriginal *string          `json:"import_payee_name_original"`
	Subtransactions         []Subtransaction `json:"subtransactions"`
	Deleted                 bool             `json:"deleted"`
}

// Subtransaction is a line item within a split transaction.
type Subtransaction struct {
	ID                    uuid.UUID  `json:"id"`
	TransactionID         uuid.UUID  `json:"transaction_id"`
	Amount                int64      `json:"amount"`
	Memo                  *string    `json:"memo"`
	PayeeID               *uuid.UUID `json:"payee_id"`
	PayeeName             *string    `json:"payee_name"`
	CategoryID            *uuid.UUID `json:"category_id"`
	CategoryName          *string    `json:"category_name"`
	TransferAccountID     *uuid.UUID `json:"transfer_account_id"`
	TransferTransactionID *uuid.UUID `json:"transfer_transaction_id"`
	Deleted               bool       `json:"deleted"`
}

// ClearedStatus represents the cleared state of a transaction.
type ClearedStatus string

const (
	ClearedStatusCleared    ClearedStatus = "cleared"
	ClearedStatusUncleared  ClearedStatus = "uncleared"
	ClearedStatusReconciled ClearedStatus = "reconciled"
)

// FlagColor represents the color of a transaction flag.
type FlagColor string

const (
	FlagColorNone   FlagColor = ""
	FlagColorRed    FlagColor = "red"
	FlagColorOrange FlagColor = "orange"
	FlagColorYellow FlagColor = "yellow"
	FlagColorGreen  FlagColor = "green"
	FlagColorBlue   FlagColor = "blue"
	FlagColorPurple FlagColor = "purple"
)

type scheduledTransactionData struct {
	Data struct {
		ScheduledTransaction ScheduledTransaction `json:"scheduled_transaction"`
		ServerKnowledge      int64                `json:"server_knowledge"`
	} `json:"data"`
}

type scheduledTransactionsData struct {
	Data struct {
		ScheduledTransactions []ScheduledTransaction `json:"scheduled_transactions"`
		ServerKnowledge       int64                  `json:"server_knowledge"`
	} `json:"data"`
}

// ScheduledTransaction represents a recurring scheduled transaction.
type ScheduledTransaction struct {
	ID                uuid.UUID                 `json:"id"`
	DateFirst         Date                      `json:"date_first"`
	DateNext          Date                      `json:"date_next"`
	Frequency         Frequency                 `json:"frequency"`
	Amount            int64                     `json:"amount"`
	Memo              *string                   `json:"memo"`
	FlagColor         *FlagColor                `json:"flag_color"`
	FlagName          *string                   `json:"flag_name"`
	AccountID         uuid.UUID                 `json:"account_id"`
	PayeeID           *uuid.UUID                `json:"payee_id"`
	CategoryID        *uuid.UUID                `json:"category_id"`
	AccountName       string                    `json:"account_name"`
	PayeeName         *string                   `json:"payee_name"`
	CategoryName      *string                   `json:"category_name"`
	Subtransactions   []ScheduledSubtransaction `json:"subtransactions"`
	TransferAccountID *uuid.UUID                `json:"transfer_account_id"`
	Deleted           bool                      `json:"deleted"`
}

// ScheduledSubtransaction is a line item within a split scheduled transaction.
type ScheduledSubtransaction struct {
	ID                     uuid.UUID  `json:"id"`
	ScheduledTransactionID uuid.UUID  `json:"scheduled_transaction_id"`
	Amount                 int64      `json:"amount"`
	Memo                   *string    `json:"memo"`
	PayeeID                *uuid.UUID `json:"payee_id"`
	PayeeName              *string    `json:"payee_name"`
	CategoryID             *uuid.UUID `json:"category_id"`
	CategoryName           *string    `json:"category_name"`
	TransferAccountID      *uuid.UUID `json:"transfer_account_id"`
	Deleted                bool       `json:"deleted"`
}

// Frequency represents the recurrence interval for a scheduled transaction.
type Frequency string

const (
	FrequencyNever           Frequency = "never"
	FrequencyDaily           Frequency = "daily"
	FrequencyWeekly          Frequency = "weekly"
	FrequencyEveryOtherWeek  Frequency = "everyOtherWeek"
	FrequencyTwiceAMonth     Frequency = "twiceAMonth"
	FrequencyEvery4Weeks     Frequency = "every4Weeks"
	FrequencyMonthly         Frequency = "monthly"
	FrequencyEveryOtherMonth Frequency = "everyOtherMonth"
	FrequencyEvery3Months    Frequency = "every3Months"
	FrequencyEvery4Months    Frequency = "every4Months"
	FrequencyTwiceAYear      Frequency = "twiceAYear"
	FrequencyYearly          Frequency = "yearly"
	FrequencyEveryOtherYear  Frequency = "everyOtherYear"
)

// GET Methods using transactions

func buildTransactionParams(params *TransactionListParams) url.Values {
	q := url.Values{}
	if params != nil {
		if params.SinceDate != nil {
			q.Set("since_date", params.SinceDate.String())
		}
		if params.Type != nil {
			q.Set("type", string(*params.Type))
		}
		if params.LastKnowledgeOfServer != nil {
			q.Set("last_knowledge_of_server", fmt.Sprintf("%d", *params.LastKnowledgeOfServer))
		}
	}
	return q
}

func (c *Client) GetTransactions(ctx context.Context, planID uuid.UUID, params *TransactionListParams) ([]Transaction, int64, error) {
	var result transactionsData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/transactions", planID), buildTransactionParams(params), &result); err != nil {
		return nil, -1, err
	}
	return result.Data.Transactions, result.Data.ServerKnowledge, nil
}

func (c *Client) GetTransaction(ctx context.Context, planID uuid.UUID, txID string) (*Transaction, int64, error) {
	var result transactionData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/transactions/%s", planID, txID), nil, &result); err != nil {
		return nil, -1, err
	}
	return &result.Data.Transaction, result.Data.ServerKnowledge, nil
}

func (c *Client) GetTransactionsByAccount(ctx context.Context, planID uuid.UUID, accountID uuid.UUID, params *TransactionListParams) ([]Transaction, int64, error) {
	q := buildTransactionParams(params)
	var result transactionsData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/accounts/%s/transactions", planID, accountID), q, &result); err != nil {
		return nil, -1, err
	}
	return result.Data.Transactions, result.Data.ServerKnowledge, nil
}

func (c *Client) GetTransactionsByCategory(ctx context.Context, planID uuid.UUID, categoryID uuid.UUID, params *TransactionListParams) ([]Transaction, int64, error) {
	q := buildTransactionParams(params)
	var result transactionsData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/categories/%s/transactions", planID, categoryID), q, &result); err != nil {
		return nil, -1, err
	}
	return result.Data.Transactions, result.Data.ServerKnowledge, nil
}

func (c *Client) GetTransactionsByPayee(ctx context.Context, planID uuid.UUID, payeeID uuid.UUID, params *TransactionListParams) ([]Transaction, int64, error) {
	q := buildTransactionParams(params)
	var result transactionsData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/payees/%s/transactions", planID, payeeID), q, &result); err != nil {
		return nil, -1, err
	}
	return result.Data.Transactions, result.Data.ServerKnowledge, nil
}

func (c *Client) GetTransactionsByMonth(ctx context.Context, planID uuid.UUID, month Date, params *TransactionListParams) ([]Transaction, int64, error) {
	q := buildTransactionParams(params)
	var result transactionsData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/months/%s/transactions", planID, month), q, &result); err != nil {
		return nil, -1, err
	}
	return result.Data.Transactions, result.Data.ServerKnowledge, nil
}

func (c *Client) GetScheduledTransactions(ctx context.Context, planID uuid.UUID, params *ListParams) ([]ScheduledTransaction, int64, error) {
	var result scheduledTransactionsData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/scheduled_transactions", planID), buildListParams(params), &result); err != nil {
		return nil, -1, err
	}
	return result.Data.ScheduledTransactions, result.Data.ServerKnowledge, nil
}

func (c *Client) GetScheduledTransaction(ctx context.Context, planID, txID uuid.UUID) (*ScheduledTransaction, error) {
	var result scheduledTransactionData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/scheduled_transactions/%s", planID, txID), nil, &result); err != nil {
		return nil, err
	}
	return &result.Data.ScheduledTransaction, nil
}

// POST Methods and infrastructure using transactions

// SaveTransaction is the request body for creating a transaction. Amount is in milliunits.
type SaveTransaction struct {
	AccountID       uuid.UUID            `json:"account_id"`
	Date            Date                 `json:"date"`
	Amount          int64                `json:"amount"`
	PayeeID         *uuid.UUID           `json:"payee_id,omitempty"`
	PayeeName       *string              `json:"payee_name,omitempty"`
	CategoryID      *uuid.UUID           `json:"category_id,omitempty"`
	Memo            *string              `json:"memo,omitempty"`
	Cleared         ClearedStatus        `json:"cleared,omitempty"`
	Approved        *bool                `json:"approved,omitempty"`
	FlagColor       *FlagColor           `json:"flag_color,omitempty"`
	ImportID        *string              `json:"import_id,omitempty"`
	Subtransactions []SaveSubtransaction `json:"subtransactions,omitempty"`
}

type saveTransactionWrapper struct {
	Transaction SaveTransaction `json:"transaction"`
}

type saveTransactionsWrapper struct {
	Transactions []SaveTransaction `json:"transactions"`
}

type createTransactionResponseData struct {
	Data CreateTransactionResponse `json:"data"`
}

// SaveSubtransaction is the request body for creating a sub-transaction. Amount is in milliunits.
type SaveSubtransaction struct {
	Amount     int64      `json:"amount"`
	PayeeID    *uuid.UUID `json:"payee_id,omitempty"`
	PayeeName  *string    `json:"payee_name,omitempty"`
	CategoryID *uuid.UUID `json:"category_id,omitempty"`
	Memo       *string    `json:"memo,omitempty"`
}

// CreateTransactionResponse is returned by CreateTransaction.
// DuplicateImportIDs contains any ImportIDs that matched existing transactions and were skipped.
type CreateTransactionResponse struct {
	TransactionIDs     []uuid.UUID `json:"transaction_ids"`
	Transaction        Transaction `json:"transaction"`
	DuplicateImportIDs []string    `json:"duplicate_import_ids"`
	ServerKnowledge    int64       `json:"server_knowledge"`
}

type createTransactionsResponseData struct {
	Data CreateTransactionsResponse `json:"data"`
}

// CreateTransactionsResponse is returned by CreateTransactions and UpdateTransactions.
type CreateTransactionsResponse struct {
	TransactionIDs     []uuid.UUID   `json:"transaction_ids"`
	Transactions       []Transaction `json:"transactions"`
	DuplicateImportIDs []string      `json:"duplicate_import_ids"`
	ServerKnowledge    int64         `json:"server_knowledge"`
}

func (c *Client) CreateTransaction(ctx context.Context, planID uuid.UUID, t SaveTransaction) (*CreateTransactionResponse, error) {
	var result createTransactionResponseData
	err := c.post(ctx, fmt.Sprintf("plans/%s/transactions", planID), saveTransactionWrapper{t}, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) CreateTransactions(ctx context.Context, planID uuid.UUID, t []SaveTransaction) (*CreateTransactionsResponse, error) {
	var result createTransactionsResponseData
	err := c.post(ctx, fmt.Sprintf("plans/%s/transactions", planID), saveTransactionsWrapper{t}, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

type importTransactionsResponseData struct {
	Data ImportTransactionsResponse `json:"data"`
}

// ImportTransactionsResponse is returned by ImportTransactions.
type ImportTransactionsResponse struct {
	TransactionIDs  []uuid.UUID `json:"transaction_ids"`
	ServerKnowledge int64       `json:"server_knowledge"`
}

func (c *Client) ImportTransactions(ctx context.Context, planID uuid.UUID) (*ImportTransactionsResponse,
	error) {
	var result importTransactionsResponseData
	err := c.post(ctx, fmt.Sprintf("plans/%s/transactions/import", planID), struct{}{}, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// SaveScheduledTransaction is the request body for creating or updating a scheduled transaction.
type SaveScheduledTransaction struct {
	AccountID  uuid.UUID  `json:"account_id"`
	Date       Date       `json:"date"`
	Amount     int64      `json:"amount"`
	CategoryID *uuid.UUID `json:"category_id,omitempty"`
	FlagColor  *FlagColor `json:"flag_color,omitempty"`
	Frequency  Frequency  `json:"frequency"`
	Memo       *string    `json:"memo,omitempty"`
	PayeeID    *uuid.UUID `json:"payee_id,omitempty"`
	PayeeName  *string    `json:"payee_name,omitempty"`
}

type saveScheduledTransactionWrapper struct {
	Transaction SaveScheduledTransaction `json:"scheduled_transaction"`
}

func (c *Client) CreateScheduledTransaction(ctx context.Context, planID uuid.UUID, st SaveScheduledTransaction) (*ScheduledTransaction, error) {
	var result scheduledTransactionData
	err := c.post(ctx, fmt.Sprintf("plans/%s/scheduled_transactions", planID),
		saveScheduledTransactionWrapper{st}, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data.ScheduledTransaction, nil
}

// DELETE Methods and infrastructure using transactions

func (c *Client) DeleteTransaction(ctx context.Context, planID uuid.UUID, transID string) (*Transaction, int64, error) {
	var result transactionData
	err := c.delete(ctx, fmt.Sprintf("plans/%s/transactions/%s", planID, transID), &result)
	if err != nil {
		return nil, -1, err
	}
	return &result.Data.Transaction, result.Data.ServerKnowledge, nil
}

func (c *Client) DeleteScheduledTransaction(ctx context.Context, planID uuid.UUID, transID uuid.UUID) (*ScheduledTransaction, int64, error) {
	var result scheduledTransactionData
	err := c.delete(ctx, fmt.Sprintf("plans/%s/scheduled_transactions/%s", planID, transID), &result)
	if err != nil {
		return nil, -1, err
	}
	return &result.Data.ScheduledTransaction, result.Data.ServerKnowledge, nil
}

// PATCH Methods and infrastructure using transactions

// UpdateTransaction is the request body for replacing a transaction (PUT). ID, AccountID, Date, and Amount are required.
type UpdateTransaction struct {
	ID         string        `json:"id"`
	AccountID  uuid.UUID     `json:"account_id"`
	Date       Date          `json:"date"`
	Amount     int64         `json:"amount"`
	PayeeID    *uuid.UUID    `json:"payee_id,omitempty"`
	PayeeName  *string       `json:"payee_name,omitempty"`
	CategoryID *uuid.UUID    `json:"category_id,omitempty"`
	Memo       *string       `json:"memo,omitempty"`
	Cleared    ClearedStatus `json:"cleared,omitempty"`
	Approved   *bool         `json:"approved,omitempty"`
	FlagColor  *FlagColor    `json:"flag_color,omitempty"`
	ImportID   *string       `json:"import_id,omitempty"`
}

// PatchTransaction is the request body for partially updating a transaction (PATCH).
// Identify the target by setting ID (a *uuid.UUID pointer) or ImportID; only non-nil fields are applied.
type PatchTransaction struct {
	ID         *string       `json:"id,omitempty"`
	ImportID   *string       `json:"import_id,omitempty"`
	AccountID  *uuid.UUID    `json:"account_id,omitempty"`
	Date       *Date         `json:"date,omitempty"`
	Amount     *int64        `json:"amount,omitempty"`
	PayeeID    *uuid.UUID    `json:"payee_id,omitempty"`
	PayeeName  *string       `json:"payee_name,omitempty"`
	CategoryID *uuid.UUID    `json:"category_id,omitempty"`
	Memo       *string       `json:"memo,omitempty"`
	Cleared    ClearedStatus `json:"cleared,omitempty"`
	Approved   *bool         `json:"approved,omitempty"`
	FlagColor  *FlagColor    `json:"flag_color,omitempty"`
}

type patchTransactionsWrapper struct {
	Transactions []PatchTransaction `json:"transactions"`
}
type updateTransactionWrapper struct {
	Transactions UpdateTransaction `json:"transaction"`
}

func (c *Client) UpdateTransactions(ctx context.Context, planID uuid.UUID, t []PatchTransaction) (*CreateTransactionsResponse, error) {
	var result createTransactionsResponseData
	err := c.patch(ctx, fmt.Sprintf("plans/%s/transactions", planID), patchTransactionsWrapper{t},
		&result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// PUT Methods and infrastructure using transactions

func (c *Client) UpdateTransaction(ctx context.Context, planID uuid.UUID, txID string, t UpdateTransaction) (*CreateTransactionResponse, error) {
	var result createTransactionResponseData
	err := c.put(ctx, fmt.Sprintf("plans/%s/transactions/%s", planID, txID), updateTransactionWrapper{t}, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) UpdateScheduledTransaction(ctx context.Context, planID uuid.UUID, txID uuid.UUID, t SaveScheduledTransaction) (*ScheduledTransaction, error) {
	var result scheduledTransactionData
	err := c.put(ctx, fmt.Sprintf("plans/%s/scheduled_transactions/%s", planID, txID), saveScheduledTransactionWrapper{t}, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data.ScheduledTransaction, nil
}
