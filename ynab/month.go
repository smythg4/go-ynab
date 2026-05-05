package ynab

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type monthData struct {
	Data struct {
		Month Month `json:"month"`
	} `json:"data"`
}

type monthsData struct {
	Data struct {
		Months          []Month `json:"months"`
		ServerKnowledge int64   `json:"server_knowledge"`
	} `json:"data"`
}

// Month represents a budget month, including all category allocations and activity.
type Month struct {
	Month        Date       `json:"month"`
	Note         *string    `json:"note"`
	Income       int64      `json:"income"`
	Budgeted     int64      `json:"budgeted"`
	Activity     int64      `json:"activity"`
	ToBeBudgeted int64      `json:"to_be_budgeted"`
	AgeOfMoney   *int       `json:"age_of_money"`
	Deleted      bool       `json:"deleted"`
	Categories   []Category `json:"categories"`
}

// GET Methods using months

// GetMonths returns all budget months for a plan.
// The second return value is server knowledge for delta requests.
func (c *Client) GetMonths(ctx context.Context, planId uuid.UUID, params *ListParams) ([]Month, int64, error) {
	var result monthsData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/months", planId), buildListParams(params), &result); err != nil {
		return nil, -1, err
	}
	return result.Data.Months, result.Data.ServerKnowledge, nil
}

// GetMonth returns a single budget month including its category details.
func (c *Client) GetMonth(ctx context.Context, planId uuid.UUID, month Date) (*Month, error) {
	var result monthData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/months/%s", planId, month), nil, &result); err != nil {
		return nil, err
	}
	return &result.Data.Month, nil
}
