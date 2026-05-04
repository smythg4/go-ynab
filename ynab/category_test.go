package ynab

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

const categoryFixture = `{"id":"123e4567-e89b-12d3-a456-426614174000","category_group_id":"223e4567-e89b-12d3-a456-426614174000","category_group_name":"Everyday Expenses","name":"Groceries","hidden":false,"original_category_group_id":null,"note":null,"budgeted":50000,"activity":-30000,"balance":20000,"goal_type":null,"goal_needs_whole_amount":null,"goal_day":null,"goal_cadence":null,"goal_cadence_frequency":null,"goal_creation_month":null,"goal_target":null,"goal_target_date":null,"goal_percentage_complete":null,"goal_months_to_budget":null,"goal_under_funded":null,"goal_overall_funded":null,"goal_overall_left":null,"goal_snoozed_at":null,"deleted":false}`

const categoryGroupFixture = `{"id":"223e4567-e89b-12d3-a456-426614174000","name":"Everyday Expenses","hidden":false,"deleted":false,"categories":[` + categoryFixture + `]}`

const categoryGroupListFixture = `{"data":{"category_groups":[` + categoryGroupFixture + `],"server_knowledge":2}}`
const categorySingleFixture = `{"data":{"category":` + categoryFixture + `}}`

func TestGetCategories(t *testing.T) {
	t.Run("returns category group list on success", func(t *testing.T) {
		client := newTestClient(categoryGroupListFixture, 200)

		groups, serverKnowledge, err := client.GetCategories(context.Background(), uuid.New())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(groups) != 1 {
			t.Fatalf("expected 1 category group, got %d", len(groups))
		}

		idWant := uuid.MustParse("223e4567-e89b-12d3-a456-426614174000")
		if groups[0].ID != idWant {
			t.Errorf("got ID %v, want %v", groups[0].ID, idWant)
		}

		if len(groups[0].Categories) != 1 {
			t.Fatalf("expected 1 category in group, got %d", len(groups[0].Categories))
		}

		if groups[0].Categories[0].Name != "Groceries" {
			t.Errorf("got category name %v, want Groceries", groups[0].Categories[0].Name)
		}

		if serverKnowledge != 2 {
			t.Errorf("got server_knowledge %v, want 2", serverKnowledge)
		}
	})
}

func TestGetCategory(t *testing.T) {
	t.Run("returns single category on success", func(t *testing.T) {
		client := newTestClient(categorySingleFixture, 200)

		category, err := client.GetCategory(context.Background(), uuid.New(), uuid.New())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		idWant := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
		if category.ID != idWant {
			t.Errorf("got ID %v, want %v", category.ID, idWant)
		}

		if category.Balance != 20000 {
			t.Errorf("got Balance %v, want 20000", category.Balance)
		}
	})
}

func TestGetCategoryForMonth(t *testing.T) {
	t.Run("returns category for month on success", func(t *testing.T) {
		client := newTestClient(categorySingleFixture, 200)

		month := Date{time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)}
		category, err := client.GetCategoryForMonth(context.Background(), uuid.New(), month, uuid.New())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		idWant := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
		if category.ID != idWant {
			t.Errorf("got ID %v, want %v", category.ID, idWant)
		}
	})
}
