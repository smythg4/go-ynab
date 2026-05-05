package ynab

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

const monthFixture = `{"month":"2024-03-01","note":null,"income":500000,"budgeted":450000,"activity":-300000,"to_be_budgeted":50000,"age_of_money":null,"deleted":false,"categories":[]}`

const monthListFixture = `{"data":{"months":[` + monthFixture + `],"server_knowledge":9}}`
const monthSingleFixture = `{"data":{"month":` + monthFixture + `}}`

func TestGetMonths(t *testing.T) {
	client, _ := newTestClient(monthListFixture, 200)

	months, serverKnowledge, err := client.GetMonths(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(months) != 1 {
		t.Fatalf("expected 1 month, got %d", len(months))
	}

	if months[0].Income != 500000 {
		t.Errorf("got Income %v, want 500000", months[0].Income)
	}

	if months[0].ToBeBudgeted != 50000 {
		t.Errorf("got ToBeBudgeted %v, want 50000", months[0].ToBeBudgeted)
	}

	if serverKnowledge != 9 {
		t.Errorf("got server_knowledge %v, want 9", serverKnowledge)
	}
}

func TestGetMonth(t *testing.T) {
	client, _ := newTestClient(monthSingleFixture, 200)

	month := Date{time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)}
	result, err := client.GetMonth(context.Background(), uuid.New(), month)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Budgeted != 450000 {
		t.Errorf("got Budgeted %v, want 450000", result.Budgeted)
	}

	if result.Activity != -300000 {
		t.Errorf("got Activity %v, want -300000", result.Activity)
	}
}
