package ynab

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type categoryData struct {
	Data struct {
		Category        Category `json:"category"`
		ServerKnowledge int64    `json:"server_knowledge"`
	} `json:"data"`
}

type categoriesData struct {
	Data struct {
		CategoryGroups  []CategoryGroup `json:"category_groups"`
		ServerKnowledge int64           `json:"server_knowledge"`
	} `json:"data"`
}

type categoryGroupData struct {
	Data struct {
		CategoryGroup   CategoryGroup `json:"category_group"`
		ServerKnowledge int64         `json:"server_knowledge"`
	} `json:"data"`
}

// CategoryGroup represents a group of budget categories.
type CategoryGroup struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	Hidden     bool       `json:"hidden"`
	Deleted    bool       `json:"deleted"`
	Categories []Category `json:"categories"`
}

// Category represents a single budget category with goal and balance information.
// Amounts are in milliunits (divide by 1000 for display).
type Category struct {
	ID                      uuid.UUID  `json:"id"`
	CategoryGroupID         uuid.UUID  `json:"category_group_id"`
	CategoryGroupName       string     `json:"category_group_name"`
	Name                    string     `json:"name"`
	Hidden                  bool       `json:"hidden"`
	OriginalCategoryGroupID *uuid.UUID `json:"original_category_group_id"`
	Note                    *string    `json:"note"`
	Budgeted                int64      `json:"budgeted"`
	Activity                int64      `json:"activity"`
	Balance                 int64      `json:"balance"`
	GoalType                *GoalType  `json:"goal_type"`
	GoalNeedsWholeAmount    *bool      `json:"goal_needs_whole_amount"`
	GoalDay                 *int       `json:"goal_day"`
	GoalCadence             *int       `json:"goal_cadence"`
	GoalCadenceFrequency    *int       `json:"goal_cadence_frequency"`
	GoalCreationMonth       *Date      `json:"goal_creation_month"`
	GoalTarget              *int64     `json:"goal_target"`
	GoalTargetDate          *Date      `json:"goal_target_date"`
	GoalPercentageComplete  *int       `json:"goal_percentage_complete"`
	GoalMonthsToBudget      *int       `json:"goal_months_to_budget"`
	GoalUnderFunded         *int64     `json:"goal_under_funded"`
	GoalOverallFunded       *int64     `json:"goal_overall_funded"`
	GoalOverallLeft         *int64     `json:"goal_overall_left"`
	GoalSnoozedAt           *time.Time `json:"goal_snoozed_at"`
	Deleted                 bool       `json:"deleted"`
}

// GoalType represents the type of savings or spending goal assigned to a category.
type GoalType string

const (
	GoalTypeTargetBalance       GoalType = "TB"
	GoalTypeTargetBalanceByDate GoalType = "TBD"
	GoalTypePlanYourSpending    GoalType = "NEED"
	GoalTypeMonthlyFunding      GoalType = "MF"
	GoalTypeDebt                GoalType = "DEBT"
)

// GET Methods using categories

func (c *Client) GetCategories(ctx context.Context, planID uuid.UUID, params *ListParams) ([]CategoryGroup, int64, error) {
	var result categoriesData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/categories", planID), buildListParams(params), &result); err != nil {
		return nil, -1, err
	}
	return result.Data.CategoryGroups, result.Data.ServerKnowledge, nil
}

func (c *Client) GetCategory(ctx context.Context, planID uuid.UUID, categoryID uuid.UUID) (*Category, error) {
	var result categoryData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/categories/%s", planID, categoryID), nil, &result); err != nil {
		return nil, err
	}
	return &result.Data.Category, nil
}

func (c *Client) GetCategoryForMonth(ctx context.Context, planID uuid.UUID, month Date, categoryID uuid.UUID) (*Category, error) {
	var result categoryData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/months/%s/categories/%s", planID, month, categoryID), nil, &result); err != nil {
		return nil, err
	}
	return &result.Data.Category, nil
}

// POST Methods and infrastructure using categories

// SaveCategory is the request body for creating or updating a category.
// All fields are optional to support partial PATCH updates; omitted fields are not changed server-side.
type SaveCategory struct {
	CategoryGroupID      *uuid.UUID `json:"category_group_id,omitempty"`
	Name                 *string    `json:"name,omitempty"`
	Note                 *string    `json:"note,omitempty"`
	GoalNeedsWholeAmount *bool      `json:"goal_needs_whole_amount,omitempty"`
	GoalTarget           *int64     `json:"goal_target,omitempty"`
	GoalTargetDate       *Date      `json:"goal_target_date,omitempty"`
}

type saveCategoryWrapper struct {
	Category SaveCategory `json:"category"`
}

func (c *Client) CreateCategory(ctx context.Context, planID uuid.UUID, sc SaveCategory) (*Category, int64, error) {
	var result categoryData
	err := c.post(ctx, fmt.Sprintf("plans/%s/categories", planID), saveCategoryWrapper{sc}, &result)
	if err != nil {
		return nil, -1, err
	}
	return &result.Data.Category, result.Data.ServerKnowledge, nil
}

// SaveCategoryGroup is the request body for creating or updating a category group.
type SaveCategoryGroup struct {
	Name string `json:"name"`
}

type saveCategoryGroupWrapper struct {
	CategoryGroup SaveCategoryGroup `json:"category_group"`
}

func (c *Client) CreateCategoryGroup(ctx context.Context, planID uuid.UUID, scg SaveCategoryGroup) (*CategoryGroup, int64, error) {
	var result categoryGroupData
	err := c.post(ctx, fmt.Sprintf("plans/%s/category_groups", planID), saveCategoryGroupWrapper{scg}, &result)
	if err != nil {
		return nil, -1, err
	}
	return &result.Data.CategoryGroup, result.Data.ServerKnowledge, nil
}

// PATCH Methods and infrastructure using categories

func (c *Client) UpdateCategory(ctx context.Context, planID, categoryID uuid.UUID, sc SaveCategory) (*Category, int64, error) {
	var result categoryData
	err := c.patch(ctx, fmt.Sprintf("plans/%s/categories/%s", planID, categoryID), saveCategoryWrapper{sc}, &result)
	if err != nil {
		return nil, -1, err
	}
	return &result.Data.Category, result.Data.ServerKnowledge, nil
}

// SaveMonthCategory is the request body for updating a category's budget for a specific month.
type SaveMonthCategory struct {
	Budgeted int64 `json:"budgeted"`
}

type saveMonthCategoryWrapper struct {
	Category SaveMonthCategory `json:"category"`
}

func (c *Client) UpdateCategoryForMonth(ctx context.Context, planID uuid.UUID, month Date, categoryID uuid.UUID, smc SaveMonthCategory) (*Category, int64, error) {
	var result categoryData
	err := c.patch(ctx, fmt.Sprintf("plans/%s/months/%s/categories/%s", planID, month, categoryID), saveMonthCategoryWrapper{smc}, &result)
	if err != nil {
		return nil, -1, err
	}
	return &result.Data.Category, result.Data.ServerKnowledge, nil
}

func (c *Client) UpdateCategoryGroup(ctx context.Context, planID, categoryGroupID uuid.UUID, scg SaveCategoryGroup) (*CategoryGroup, int64, error) {
	var result categoryGroupData
	err := c.patch(ctx, fmt.Sprintf("plans/%s/category_groups/%s", planID, categoryGroupID), saveCategoryGroupWrapper{scg}, &result)
	if err != nil {
		return nil, -1, err
	}
	return &result.Data.CategoryGroup, result.Data.ServerKnowledge, nil
}
