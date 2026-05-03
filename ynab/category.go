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

type monthCategoryData struct {
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

type CategoryGroup struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	Hidden     bool       `json:"hidden"`
	Deleted    bool       `json:"deleted"`
	Categories []Category `json:"categories"`
}

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

type GoalType string

const (
	GoalTypeTargetBalance       GoalType = "TB"
	GoalTypeTargetBalanceByDate GoalType = "TBD"
	GoalTypePlanYourSpending    GoalType = "NEED"
	GoalTypeMonthlyFunding      GoalType = "MF"
	GoalTypeDebt                GoalType = "DEBT"
)

// GET Methods using categories
func (c *Client) GetCategories(ctx context.Context, planId uuid.UUID) ([]CategoryGroup, error) {
	// TODO: Consider how to return the `ServerKnowledge` retrieved from the query
	var result categoriesData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/categories", planId), nil, &result); err != nil {
		return nil, err
	}
	return result.Data.CategoryGroups, nil
}

func (c *Client) GetCategory(ctx context.Context, planId uuid.UUID, categoryId uuid.UUID) (*Category, error) {
	var result categoryData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/categories/%s", planId, categoryId), nil, &result); err != nil {
		return nil, err
	}
	return &result.Data.Category, nil
}

func (c *Client) GetCategoryForMonth(ctx context.Context, planId uuid.UUID, month Date, categoryId uuid.UUID) (*Category, error) {
	var result categoryData
	if err := c.get(ctx, fmt.Sprintf("plans/%s/months/%s/categories/%s", planId, month, categoryId), nil, &result); err != nil {
		return nil, err
	}
	return &result.Data.Category, nil
}

// POST Methods and infrastructure using categories
type SaveCategory struct {
	CategoryGroupID      uuid.UUID `json:"category_group_id"`
	Name                 string    `json:"name"`
	Note                 *string   `json:"note,omitempty"`
	GoalNeedsWholeAmount *bool     `json:"goal_needs_whole_amount,omitempty"`
	GoalTarget           *int64    `json:"goal_target,omitempty"`
	GoalTargetDate       *Date     `json:"goal_target_date,omitempty"`
}

type SaveCategoryWrapper struct {
	Category SaveCategory `json:"category"`
}

func (c *Client) CreateCategory(ctx context.Context, planId uuid.UUID, sc SaveCategory) (*Category, error) {
	var result categoryData
	err := c.post(ctx, fmt.Sprintf("plans/%s/categories", planId), SaveCategoryWrapper{sc}, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data.Category, nil
}

type SaveCategoryGroup struct {
	Name string `json:"name"`
}

type SaveCategoryGroupWrapper struct {
	CategoryGroup SaveCategoryGroup `json:"category_group"`
}

func (c *Client) CreateCategoryGroup(ctx context.Context, planId uuid.UUID, scg SaveCategoryGroup) (*CategoryGroup, error) {
	var result categoryGroupData
	err := c.post(ctx, fmt.Sprintf("plans/%s/category_groups", planId), SaveCategoryGroupWrapper{scg}, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data.CategoryGroup, nil
}

// PATCH Methods and infrastructure using categories
func (c *Client) UpdateCategory(ctx context.Context, planId, categoryId uuid.UUID, sc SaveCategory) (*Category, error) {
	var result categoryData
	err := c.patch(ctx, fmt.Sprintf("plans/%s/categories/%s", planId, categoryId), SaveCategoryWrapper{sc}, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data.Category, nil
}

type SaveMonthCategory struct {
	Budgeted int64 `json:"budgeted"`
}

type SaveMonthCategoryWrapper struct {
	Category SaveMonthCategory `json:"category"`
}

func (c *Client) UpdateCategoryForMonth(ctx context.Context, planId uuid.UUID, month Date, categoryId uuid.UUID, smc SaveMonthCategory) (*Category, error) {
	var result categoryData
	err := c.patch(ctx, fmt.Sprintf("plans/%s/months/%s/categories/%s", planId, month, categoryId), SaveMonthCategoryWrapper{smc}, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data.Category, nil
}

func (c *Client) UpdateCategoryGroup(ctx context.Context, planId, categoryGroupId uuid.UUID, scg SaveCategoryGroup) (*CategoryGroup, error) {
	var result categoryGroupData
	err := c.patch(ctx, fmt.Sprintf("plans/%s/category_groups/%s", planId, categoryGroupId), SaveCategoryGroupWrapper{scg}, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data.CategoryGroup, nil
}
