package ynab

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// AccountType represents the type of a YNAB account.
type AccountType string

const (
	AccountTypeChecking       AccountType = "checking"
	AccountTypeSavings        AccountType = "savings"
	AccountTypeCash           AccountType = "cash"
	AccountTypeCreditCard     AccountType = "creditCard"
	AccountTypeOtherAsset     AccountType = "otherAsset"
	AccountTypeOtherLiability AccountType = "otherLiability"
)

type accountData struct {
	Data struct {
		Account Account `json:"account"`
	} `json:"data"`
}

type accountsData struct {
	Data struct {
		Accounts        []Account `json:"accounts"`
		ServerKnowledge int64     `json:"server_knowledge"`
	} `json:"data"`
}

// Account represents a YNAB account such as checking, savings, or credit card.
type Account struct {
	ID                  uuid.UUID   `json:"id"`
	Name                string      `json:"name"`
	Type                AccountType `json:"type"`
	OnBudget            bool        `json:"on_budget"`
	Closed              bool        `json:"closed"`
	Note                *string     `json:"note"`
	Balance             int64       `json:"balance"`
	ClearedBalance      int64       `json:"cleared_balance"`
	UnclearedBalance    int64       `json:"uncleared_balance"`
	TransferPayeeID     *uuid.UUID  `json:"transfer_payee_id"`
	DirectImportLinked  bool        `json:"direct_import_linked"`
	DirectImportInError bool        `json:"direct_import_in_error"`
	LastReconciledAt    *time.Time  `json:"last_reconciled_at"`
	Deleted             bool        `json:"deleted"`
}

// GET Methods using accounts

// GetAccounts returns all accounts for a plan. The second return value is server knowledge for delta requests.
func (c *Client) GetAccounts(ctx context.Context, planId uuid.UUID, params *ListParams) ([]Account, int64, error) {
	var result accountsData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/accounts", planId), buildListParams(params), &result); err != nil {
		return nil, -1, err
	}
	return result.Data.Accounts, result.Data.ServerKnowledge, nil
}

// GetAccount returns a single account by ID.
func (c *Client) GetAccount(ctx context.Context, planId uuid.UUID, accountId uuid.UUID) (*Account, error) {
	var result accountData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/accounts/%s", planId, accountId), nil, &result); err != nil {
		return nil, err
	}
	return &result.Data.Account, nil
}

// POST Methods and infrastructure using accounts

// SaveAccount is the request body for creating a new account.
type SaveAccount struct {
	Name    string      `json:"name"`
	Type    AccountType `json:"type"`
	Balance int64       `json:"balance"`
}

type saveAccountWrapper struct {
	Account SaveAccount `json:"account"`
}

// CreateAccount creates a new account for a plan.
func (c *Client) CreateAccount(ctx context.Context, planId uuid.UUID, a SaveAccount) (*Account, error) {
	var result accountData
	err := c.post(ctx, fmt.Sprintf("plans/%s/accounts", planId), saveAccountWrapper{a}, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data.Account, nil
}
