package ynab

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
)

type planSummaryData struct {
	Data struct {
		Plans       []Plan `json:"plans"`
		DefaultPlan *Plan  `json:"default_plan"`
	} `json:"data"`
}

// Plan represents a YNAB budget plan.
type Plan struct {
	ID             uuid.UUID      `json:"id"`
	Name           string         `json:"name"`
	LastModifiedOn time.Time      `json:"last_modified_on"`
	FirstMonth     Date           `json:"first_month"`
	LastMonth      Date           `json:"last_month"`
	DateFormat     DateFormat     `json:"date_format"`
	CurrencyFormat CurrencyFormat `json:"currency_format"`
	Accounts       []Account      `json:"accounts"`
}

type planSettingsData struct {
	Data struct {
		Settings PlanSettings `json:"settings"`
	} `json:"data"`
}

// PlanSettings contains the date and currency format preferences for a plan.
type PlanSettings struct {
	DateFormat     DateFormat     `json:"date_format"`
	CurrencyFormat CurrencyFormat `json:"currency_format"`
}

type planDetailsData struct {
	Data struct {
		Plan            PlanDetails `json:"plan"`
		ServerKnowledge int64       `json:"server_knowledge"`
	} `json:"data"`
}

// PlanDetails is the full plan export returned by GetPlan, including all
// accounts, categories, transactions, and other sub-resources.
type PlanDetails struct {
	Plan
	Payees                   []Payee                   `json:"payees"`
	PayeeLocations           []PayeeLocation           `json:"payee_locations"`
	CategoryGroups           []CategoryGroup           `json:"category_groups"`
	Categories               []Category                `json:"categories"`
	Months                   []Month                   `json:"months"`
	Transactions             []Transaction             `json:"transactions"`
	Subtransactions          []Subtransaction          `json:"subtransactions"`
	ScheduledTransactions    []ScheduledTransaction    `json:"scheduled_transactions"`
	ScheduledSubtransactions []ScheduledSubtransaction `json:"scheduled_subtransactions"`
}

// GET Methods using plans

func (c *Client) GetPlans(ctx context.Context, includeAccounts bool) ([]Plan, error) {
	q := url.Values{}
	if includeAccounts {
		q.Set("include_accounts", "true")
	}
	var result planSummaryData
	if err := c.get(ctx, "plans", q, &result); err != nil {
		return nil, err
	}
	return result.Data.Plans, nil
}

func (c *Client) GetPlan(ctx context.Context, id uuid.UUID, params *ListParams) (*PlanDetails, int64, error) {
	var result planDetailsData
	if err := c.get(ctx, fmt.Sprintf("plans/%s", id), buildListParams(params), &result); err != nil {
		return nil, -1, err
	}
	return &result.Data.Plan, result.Data.ServerKnowledge, nil
}

func (c *Client) GetLastUsedPlan(ctx context.Context) (*PlanDetails, error) {
	var result planDetailsData
	if err := c.get(ctx, "plans/last-used", nil, &result); err != nil {
		return nil, err
	}
	return &result.Data.Plan, nil
}

func (c *Client) GetPlanSettings(ctx context.Context, id uuid.UUID) (*PlanSettings, error) {
	var result planSettingsData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/settings", id), nil, &result); err != nil {
		return nil, err
	}
	return &result.Data.Settings, nil
}
