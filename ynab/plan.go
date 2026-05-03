package ynab

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PlanSummaryData struct {
	Data struct {
		Plans       []Plan `json:"plans"`
		DefaultPlan *Plan  `json:"default_plan"`
	} `json:"data"`
}

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
func (c *Client) GetPlans(ctx context.Context) ([]Plan, error) {
	var result PlanSummaryData
	if err := c.get(ctx, "plans", nil, &result); err != nil {
		return nil, err
	}
	return result.Data.Plans, nil
}

func (c *Client) GetPlan(ctx context.Context, id uuid.UUID) (*PlanDetails, error) {
	var result planDetailsData
	if err := c.get(ctx, fmt.Sprintf("plans/%s", id), nil, &result); err != nil {
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
