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
func (c *Client) GetMonths(ctx context.Context, planId uuid.UUID) ([]Month, error) {
	// TODO: Consider how to return the `ServerKnowledge` retrieved from the query
	var result monthsData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/months", planId), nil, &result); err != nil {
		return nil, err
	}
	return result.Data.Months, nil
}

func (c *Client) GetMonth(ctx context.Context, planId uuid.UUID, month Date) (*Month, error) {
	var result monthData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/months/%s", planId, month), nil, &result); err != nil {
		return nil, err
	}
	return &result.Data.Month, nil
}
