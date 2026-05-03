package ynab

import (
	"context"
	"fmt"
	"net/url"

	"github.com/google/uuid"
)

type transactionData struct {
	Data struct {
		Transaction Transaction `json:"transaction"`
	} `json:"data"`
}

type transactionsData struct {
	Data struct {
		Transactions    []Transaction `json:"transactions"`
		ServerKnowledge int64         `json:"server_knowledge"`
	} `json:"data"`
}

type Transaction struct {
	ID                   string           `json:"id"`
	Date                 Date             `json:"date"`
	Amount               int64            `json:"amount"`
	Memo                 *string          `json:"memo"`
	Cleared              ClearedStatus    `json:"cleared"`
	Approved             bool             `json:"approved"`
	FlagColor            *FlagColor       `json:"flag_color"`
	FlagName             *string          `json:"flag_name"`
	AccountID            uuid.UUID        `json:"account_id"`
	PayeeID              *uuid.UUID       `json:"payee_id"`
	AccountName          string           `json:"account_name"`
	PayeeName            *string          `json:"payee_name"`
	CategoryID           *uuid.UUID       `json:"category_id"`
	CategoryName         *string          `json:"category_name"`
	MatchedTransactionID *string          `json:"matched_transaction_id"`
	Subtransactions      []Subtransaction `json:"subtransactions"`
}

type Subtransaction struct {
	ID                    string     `json:"id"`
	TransactionID         string     `json:"transaction_id"`
	Amount                int64      `json:"amount"`
	Memo                  *string    `json:"memo"`
	PayeeID               *uuid.UUID `json:"payee_id"`
	PayeeName             *string    `json:"payee_name"`
	CategoryID            *uuid.UUID `json:"category_id"`
	CategoryName          *string    `json:"category_name"`
	TransferAccountID     *uuid.UUID `json:"transfer_account_id"`
	TransferTransactionID *string    `json:"transfer_transaction_id"`
}

type ClearedStatus string

const (
	ClearedStatusCleared    ClearedStatus = "cleared"
	ClearedStatusUncleared  ClearedStatus = "uncleared"
	ClearedStatusReconciled ClearedStatus = "reconciled"
)

type FlagColor string

const (
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
	} `json:"data"`
}

type scheduledTransactionsData struct {
	Data struct {
		ScheduledTransactions []ScheduledTransaction `json:"scheduled_transactions"`
		ServerKnowledge       int64                  `json:"server_knowledge"`
	} `json:"data"`
}

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
}

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
	FrequencyEvery3Months    Frequency = "everyThreeMonths"
	FrequencyEvery4Months    Frequency = "everyFourMonths"
	FrequencyTwiceAYear      Frequency = "twiceAYear"
	FrequencyYearly          Frequency = "yearly"
	FrequencyEveryOtherYear  Frequency = "everyOtherYear"
)

// GET Methods using transactions
func (c *Client) GetTransactions(ctx context.Context, planId uuid.UUID, params *TransactionListParams) ([]Transaction, int64, error) {
	q := url.Values{}
	if params != nil {
		if params.SinceDate != nil {
			q.Set("since_date", params.SinceDate.String())
		}
		if params.Type != nil {
			q.Set("type", *params.Type)
		}
		if params.LastKnowledgeOfServer != nil {
			q.Set("last_knowledge_of_server", fmt.Sprintf("%d", *params.LastKnowledgeOfServer))
		}
	}
	var result transactionsData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/transactions", planId), q, &result); err != nil {
		return nil, -1, err
	}
	return result.Data.Transactions, result.Data.ServerKnowledge, nil
}

func (c *Client) GetTransaction(ctx context.Context, planId uuid.UUID, txId string) (*Transaction, error) {
	var result transactionData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/transactions/%s", planId, txId), nil, &result); err != nil {
		return nil, err
	}
	return &result.Data.Transaction, nil
}

func (c *Client) GetTransactionsByAccount(ctx context.Context, planId uuid.UUID, accountId uuid.UUID) ([]Transaction, int64, error) {
	var result transactionsData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/accounts/%s/transactions", planId, accountId), nil, &result); err != nil {
		return nil, -1, err
	}
	return result.Data.Transactions, result.Data.ServerKnowledge, nil
}

func (c *Client) GetTransactionsByCategory(ctx context.Context, planId uuid.UUID, categoryId uuid.UUID) ([]Transaction, int64, error) {
	var result transactionsData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/categories/%s/transactions", planId, categoryId), nil, &result); err != nil {
		return nil, -1, err
	}
	return result.Data.Transactions, result.Data.ServerKnowledge, nil
}

func (c *Client) GetTransactionsByPayee(ctx context.Context, planId uuid.UUID, payeeId uuid.UUID) ([]Transaction, int64, error) {
	var result transactionsData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/payees/%s/transactions", planId, payeeId), nil, &result); err != nil {
		return nil, -1, err
	}
	return result.Data.Transactions, result.Data.ServerKnowledge, nil
}

func (c *Client) GetTransactionsByMonth(ctx context.Context, planId uuid.UUID, month Date) ([]Transaction, int64, error) {
	var result transactionsData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/months/%s/transactions", planId, month), nil, &result); err != nil {
		return nil, -1, err
	}
	return result.Data.Transactions, result.Data.ServerKnowledge, nil
}

func (c *Client) GetScheduledTransactions(ctx context.Context, planId uuid.UUID) ([]ScheduledTransaction, int64, error) {
	var result scheduledTransactionsData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/scheduled_transactions", planId), nil, &result); err != nil {
		return nil, -1, err
	}
	return result.Data.ScheduledTransactions, result.Data.ServerKnowledge, nil
}

func (c *Client) GetScheduledTransaction(ctx context.Context, planId, txId uuid.UUID) (*ScheduledTransaction, error) {
	var result scheduledTransactionData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/scheduled_transactions/%s", planId, txId), nil, &result); err != nil {
		return nil, err
	}
	return &result.Data.ScheduledTransaction, nil
}

// POST Methods and infrastructure using transactions
type SaveTransaction struct {
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

type SaveTransactionWrapper struct {
	Transaction SaveTransaction `json:"transaction"`
}

type SaveTransactionsWrapper struct {
	Transactions []SaveTransaction `json:"transactions"`
}

type createTransactionResponseData struct {
	Data CreateTransactionResponse `json:"data"`
}

type CreateTransactionResponse struct {
	TransactionIDs     []string    `json:"transaction_ids"`
	Transaction        Transaction `json:"transaction"`
	DuplicateImportIDs []string    `json:"duplicate_import_ids"`
	ServerKnowledge    int64       `json:"server_knowledge"`
}

type createTransactionsResponseData struct {
	Data CreateTransactionsResponse `json:"data"`
}

type CreateTransactionsResponse struct {
	TransactionIDs     []string      `json:"transaction_ids"`
	Transactions       []Transaction `json:"transactions"`
	DuplicateImportIDs []string      `json:"duplicate_import_ids"`
	ServerKnowledge    int64         `json:"server_knowledge"`
}

func (c *Client) CreateTransaction(ctx context.Context, planId uuid.UUID, t SaveTransaction) (*CreateTransactionResponse, error) {
	var result createTransactionResponseData
	err := c.post(ctx, fmt.Sprintf("plans/%s/transactions", planId), SaveTransactionWrapper{t}, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) CreateTransactions(ctx context.Context, planId uuid.UUID, t []SaveTransaction) (*CreateTransactionsResponse, error) {
	var result createTransactionsResponseData
	err := c.post(ctx, fmt.Sprintf("plans/%s/transactions", planId), SaveTransactionsWrapper{t}, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

type importTransactionsResponseData struct {
	Data ImportTransactionsResponse `json:"data"`
}

type ImportTransactionsResponse struct {
	TransactionIDs  []string `json:"transaction_ids"`
	ServerKnowledge int64    `json:"server_knowledge"`
}

func (c *Client) ImportTransactions(ctx context.Context, planId uuid.UUID) (*ImportTransactionsResponse,
	error) {
	var result importTransactionsResponseData
	err := c.post(ctx, fmt.Sprintf("plans/%s/transactions/import", planId), struct{}{}, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

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

type SaveScheduledTransactionWrapper struct {
	Transaction SaveScheduledTransaction `json:"scheduled_transaction"`
}

func (c *Client) CreateScheduledTransaction(ctx context.Context, planId uuid.UUID, st SaveScheduledTransaction) (*ScheduledTransaction, error) {
	var result scheduledTransactionData
	err := c.post(ctx, fmt.Sprintf("plans/%s/scheduled_transactions", planId),
		SaveScheduledTransactionWrapper{st}, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data.ScheduledTransaction, nil
}

// DELETE Methods and infrastructure using transactions
func (c *Client) DeleteTransaction(ctx context.Context, planId uuid.UUID, transId string) (*Transaction, error) {
	var result transactionData
	err := c.delete(ctx, fmt.Sprintf("plans/%s/transactions/%s", planId, transId), &result)
	if err != nil {
		return nil, err
	}
	return &result.Data.Transaction, nil
}

func (c *Client) DeleteScheduledTransaction(ctx context.Context, planId uuid.UUID, transId uuid.UUID) (*ScheduledTransaction, error) {
	var result scheduledTransactionData
	err := c.delete(ctx, fmt.Sprintf("plans/%s/scheduled_transactions/%s", planId, transId), &result)
	if err != nil {
		return nil, err
	}
	return &result.Data.ScheduledTransaction, nil
}

// PATCH Methods and infrastructure using transactions
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

type UpdateTransactionsWrapper struct {
	Transactions []UpdateTransaction `json:"transactions"`
}
type UpdateTransactionWrapper struct {
	Transactions UpdateTransaction `json:"transactions"`
}

func (c *Client) UpdateTransactions(ctx context.Context, planId uuid.UUID, t []UpdateTransaction) (*CreateTransactionsResponse, error) {
	var result createTransactionsResponseData
	err := c.patch(ctx, fmt.Sprintf("plans/%s/transactions", planId), UpdateTransactionsWrapper{t},
		&result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// PUT Methods and infrastructure using transactions
func (c *Client) UpdateTransaction(ctx context.Context, planId uuid.UUID, txId string, t UpdateTransaction) (*CreateTransactionResponse, error) {
	var result createTransactionResponseData
	err := c.put(ctx, fmt.Sprintf("plans/%s/transactions/%s", planId, txId), UpdateTransactionWrapper{t}, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) UpdateScheduledTransaction(ctx context.Context, planId uuid.UUID, txId uuid.UUID, t SaveScheduledTransaction) (*ScheduledTransaction, error) {
	var result scheduledTransactionData
	err := c.put(ctx, fmt.Sprintf("plans/%s/scheduled_transactions/%s", planId, txId), SaveScheduledTransactionWrapper{t}, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data.ScheduledTransaction, nil
}
