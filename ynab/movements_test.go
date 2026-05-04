package ynab

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

const movementFixture = `{"id":"123e4567-e89b-12d3-a456-426614174000","month":"2024-03-01","moved_at":null,"note":null,"money_movement_group_id":null,"performed_by_user_id":null,"from_category_id":"223e4567-e89b-12d3-a456-426614174000","to_category_id":"323e4567-e89b-12d3-a456-426614174000","amount":5000}`

const movementListFixture = `{"data":{"money_movements":[` + movementFixture + `],"server_knowledge":4}}`

const movementGroupFixture = `{"id":"423e4567-e89b-12d3-a456-426614174000","group_created_at":"2024-03-15T10:00:00Z","month":"2024-03-01","note":null,"performed_by_user_id":null}`

const movementGroupListFixture = `{"data":{"money_movement_groups":[` + movementGroupFixture + `],"server_knowledge":6}}`

func TestGetMoneyMovements(t *testing.T) {
	t.Run("returns money movement list on success", func(t *testing.T) {
		client, _ := newTestClient(movementListFixture, 200)

		movements, err := client.GetMoneyMovements(context.Background(), uuid.New())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(movements) != 1 {
			t.Fatalf("expected 1 movement, got %d", len(movements))
		}

		idWant := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
		if movements[0].ID != idWant {
			t.Errorf("got ID %v, want %v", movements[0].ID, idWant)
		}

		if movements[0].Amount != 5000 {
			t.Errorf("got Amount %v, want 5000", movements[0].Amount)
		}
	})
}

func TestGetMoneyMovementsByMonth(t *testing.T) {
	t.Run("returns money movements for month on success", func(t *testing.T) {
		client, _ := newTestClient(movementListFixture, 200)

		month := Date{time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)}
		movements, serverKnowledge, err := client.GetMoneyMovementsByMonth(context.Background(), uuid.New(), month)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(movements) != 1 {
			t.Fatalf("expected 1 movement, got %d", len(movements))
		}

		if serverKnowledge != 4 {
			t.Errorf("got server_knowledge %v, want 4", serverKnowledge)
		}
	})
}

func TestGetMoneyMovementGroups(t *testing.T) {
	t.Run("returns money movement group list on success", func(t *testing.T) {
		client, _ := newTestClient(movementGroupListFixture, 200)

		groups, serverKnowledge, err := client.GetMoneyMovementGroups(context.Background(), uuid.New())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(groups) != 1 {
			t.Fatalf("expected 1 group, got %d", len(groups))
		}

		idWant := uuid.MustParse("423e4567-e89b-12d3-a456-426614174000")
		if groups[0].ID != idWant {
			t.Errorf("got ID %v, want %v", groups[0].ID, idWant)
		}

		if serverKnowledge != 6 {
			t.Errorf("got server_knowledge %v, want 6", serverKnowledge)
		}
	})
}

func TestGetMoneyMovementGroupsByMonth(t *testing.T) {
	t.Run("returns money movement groups for month on success", func(t *testing.T) {
		client, _ := newTestClient(movementGroupListFixture, 200)

		month := Date{time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)}
		groups, serverKnowledge, err := client.GetMoneyMovementGroupsByMonth(context.Background(), uuid.New(), month)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(groups) != 1 {
			t.Fatalf("expected 1 group, got %d", len(groups))
		}

		if serverKnowledge != 6 {
			t.Errorf("got server_knowledge %v, want 6", serverKnowledge)
		}
	})
}
