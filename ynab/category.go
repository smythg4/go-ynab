package ynab

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type categoryData struct {
	Data struct {
		Category Category `json:"category"`
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

// GetCategories returns all category groups and their categories for a plan.
// The second return value is server knowledge for delta requests.
func (c *Client) GetCategories(ctx context.Context, planId uuid.UUID, params *ListParams) ([]CategoryGroup, int64, error) {
	var result categoriesData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/categories", planId), buildListParams(params), &result); err != nil {
		return nil, -1, err
	}
	return result.Data.CategoryGroups, result.Data.ServerKnowledge, nil
}

// GetCategory returns a single category by ID.
func (c *Client) GetCategory(ctx context.Context, planId uuid.UUID, categoryId uuid.UUID) (*Category, error) {
	var result categoryData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/categories/%s", planId, categoryId), nil, &result); err != nil {
		return nil, err
	}
	return &result.Data.Category, nil
}

// GetCategoryForMonth returns a category's data for a specific budget month.
func (c *Client) GetCategoryForMonth(ctx context.Context, planId uuid.UUID, month Date, categoryId uuid.UUID) (*Category, error) {
	var result categoryData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/months/%s/categories/%s", planId, month, categoryId), nil, &result); err != nil {
		return nil, err
	}
	return &result.Data.Category, nil
}

// POST Methods and infrastructure using categories

// SaveCategory is the request body for creating or updating a category.
type SaveCategory struct {
	CategoryGroupID      uuid.UUID `json:"category_group_id"`
	Name                 string    `json:"name"`
	Note                 *string   `json:"note,omitempty"`
	GoalNeedsWholeAmount *bool     `json:"goal_needs_whole_amount,omitempty"`
	GoalTarget           *int64    `json:"goal_target,omitempty"`
	GoalTargetDate       *Date     `json:"goal_target_date,omitempty"`
}

type saveCategoryWrapper struct {
	Category SaveCategory `json:"category"`
}

// CreateCategory creates a new category within a category group.
func (c *Client) CreateCategory(ctx context.Context, planId uuid.UUID, sc SaveCategory) (*Category, error) {
	var result categoryData
	err := c.post(ctx, fmt.Sprintf("plans/%s/categories", planId), saveCategoryWrapper{sc}, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data.Category, nil
}

// SaveCategoryGroup is the request body for creating or updating a category group.
type SaveCategoryGroup struct {
	Name string `json:"name"`
}

type saveCategoryGroupWrapper struct {
	CategoryGroup SaveCategoryGroup `json:"category_group"`
}

// CreateCategoryGroup creates a new category group.
func (c *Client) CreateCategoryGroup(ctx context.Context, planId uuid.UUID, scg SaveCategoryGroup) (*CategoryGroup, error) {
	var result categoryGroupData
	err := c.post(ctx, fmt.Sprintf("plans/%s/category_groups", planId), saveCategoryGroupWrapper{scg}, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data.CategoryGroup, nil
}

// PATCH Methods and infrastructure using categories

// UpdateCategory updates an existing category.
func (c *Client) UpdateCategory(ctx context.Context, planId, categoryId uuid.UUID, sc SaveCategory) (*Category, error) {
	var result categoryData
	err := c.patch(ctx, fmt.Sprintf("plans/%s/categories/%s", planId, categoryId), saveCategoryWrapper{sc}, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data.Category, nil
}

// SaveMonthCategory is the request body for updating a category's budget for a specific month.
type SaveMonthCategory struct {
	Budgeted int64 `json:"budgeted"`
}

type saveMonthCategoryWrapper struct {
	Category SaveMonthCategory `json:"category"`
}

// UpdateCategoryForMonth updates a category's budgeted amount for a specific month.
func (c *Client) UpdateCategoryForMonth(ctx context.Context, planId uuid.UUID, month Date, categoryId uuid.UUID, smc SaveMonthCategory) (*Category, error) {
	var result categoryData
	err := c.patch(ctx, fmt.Sprintf("plans/%s/months/%s/categories/%s", planId, month, categoryId), saveMonthCategoryWrapper{smc}, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data.Category, nil
}

// UpdateCategoryGroup updates an existing category group.
func (c *Client) UpdateCategoryGroup(ctx context.Context, planId, categoryGroupId uuid.UUID, scg SaveCategoryGroup) (*CategoryGroup, error) {
	var result categoryGroupData
	err := c.patch(ctx, fmt.Sprintf("plans/%s/category_groups/%s", planId, categoryGroupId), saveCategoryGroupWrapper{scg}, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data.CategoryGroup, nil
}
