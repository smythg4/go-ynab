package ynab

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type moneyMovementsData struct {
	Data struct {
		MoneyMovements  []MoneyMovement `json:"money_movements"`
		ServerKnowledge int64           `json:"server_knowledge"`
	} `json:"data"`
}

// MoneyMovement represents a single movement of money between categories.
type MoneyMovement struct {
	ID                   uuid.UUID  `json:"id"`
	Month                Date       `json:"month"`
	MovedAt              *time.Time `json:"moved_at"`
	Note                 *string    `json:"note"`
	MoneyMovementGroupID *uuid.UUID `json:"money_movement_group_id"`
	PerformedByUserID    *uuid.UUID `json:"performed_by_user_id"`
	FromCategoryID       *uuid.UUID `json:"from_category_id"`
	ToCategoryID         *uuid.UUID `json:"to_category_id"`
	Amount               int64      `json:"amount"`
}

type moneyMovementGroupData struct {
	Data struct {
		MoneyMovementGroups []MoneyMovementGroup `json:"money_movement_groups"`
		ServerKnowledge     int64                `json:"server_knowledge"`
	} `json:"data"`
}

// MoneyMovementGroup represents a group of related money movements.
type MoneyMovementGroup struct {
	ID                uuid.UUID  `json:"id"`
	GroupCreatedAt    time.Time  `json:"group_created_at"`
	Month             Date       `json:"month"`
	Note              *string    `json:"note"`
	PerformedByUserID *uuid.UUID `json:"performed_by_user_id"`
}

// GET Methods using money movements

// GetMoneyMovements returns all money movements for a plan.
// The second return value is server knowledge for delta requests.
func (c *Client) GetMoneyMovements(ctx context.Context, planId uuid.UUID, params *ListParams) ([]MoneyMovement, int64, error) {
	var result moneyMovementsData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/money_movements", planId), buildListParams(params), &result); err != nil {
		return nil, -1, err
	}
	return result.Data.MoneyMovements, result.Data.ServerKnowledge, nil
}

// GetMoneyMovementsByMonth returns money movements for a specific budget month.
func (c *Client) GetMoneyMovementsByMonth(ctx context.Context, planId uuid.UUID, month Date) ([]MoneyMovement, int64, error) {
	var result moneyMovementsData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/months/%s/money_movements", planId, month), nil, &result); err != nil {
		return nil, -1, err
	}
	return result.Data.MoneyMovements, result.Data.ServerKnowledge, nil
}

// GetMoneyMovementGroups returns all money movement groups for a plan.
// The second return value is server knowledge for delta requests.
func (c *Client) GetMoneyMovementGroups(ctx context.Context, planId uuid.UUID, params *ListParams) ([]MoneyMovementGroup, int64, error) {
	var result moneyMovementGroupData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/money_movement_groups", planId), buildListParams(params), &result); err != nil {
		return nil, -1, err
	}
	return result.Data.MoneyMovementGroups, result.Data.ServerKnowledge, nil
}

// GetMoneyMovementGroupsByMonth returns money movement groups for a specific budget month.
func (c *Client) GetMoneyMovementGroupsByMonth(ctx context.Context, planId uuid.UUID, month Date) ([]MoneyMovementGroup, int64, error) {
	var result moneyMovementGroupData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/months/%s/money_movement_groups", planId, month), nil, &result); err != nil {
		return nil, -1, err
	}
	return result.Data.MoneyMovementGroups, result.Data.ServerKnowledge, nil
}
