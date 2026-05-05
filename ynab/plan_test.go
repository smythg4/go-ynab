package ynab

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
)

const currencyFormatFixture = `"date_format":{"format":"MM/DD/YYYY"},"currency_format":{"iso_code":"USD","example_format":"123,456.78","decimal_digits":2,"decimal_separator":".","symbol_first":true,"group_separator":",","currency_symbol":"$","display_symbol":true}`

const planFixture1 = `{"id":"` + testID1 + `","name":"Bunko Budget","last_modified_on":"2024-03-01T00:00:00Z","first_month":"2024-01-01","last_month":"2024-12-01",` + currencyFormatFixture + `,"accounts":[]}`
const planFixture2 = `{"id":"` + testID2 + `","name":"Side Hustle Budget","last_modified_on":"2024-04-01T00:00:00Z","first_month":"2024-02-01","last_month":"2024-12-01",` + currencyFormatFixture + `,"accounts":[]}`

const planListFixture = `{"data":{"plans":[` + planFixture1 + `,` + planFixture2 + `],"default_plan":null}}`
const planDetailsFixture = `{"data":{"plan":{"id":"` + testID1 + `","name":"Bunko Budget","last_modified_on":"2024-03-01T00:00:00Z","first_month":"2024-01-01","last_month":"2024-12-01",` + currencyFormatFixture + `,"accounts":[],"payees":[],"payee_locations":[],"category_groups":[],"categories":[],"months":[],"transactions":[],"subtransactions":[],"scheduled_transactions":[],"scheduled_subtransactions":[]},"server_knowledge":42}}`
const planSettingsFixture = `{"data":{"settings":{` + currencyFormatFixture + `}}}`

func TestGetPlans(t *testing.T) {
	client, _ := newTestClient(planListFixture, 200)

	plans, err := client.GetPlans(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(plans) != 2 {
		t.Fatalf("expected 2 plans, got %d", len(plans))
	}

	idWant := uuid.MustParse(testID1)
	if plans[0].ID != idWant {
		t.Errorf("got ID %v, want %v", plans[0].ID, idWant)
	}

	if plans[0].Name != "Bunko Budget" {
		t.Errorf("got Name %v, want Bunko Budget", plans[0].Name)
	}
}

func TestGetPlan(t *testing.T) {
	client, _ := newTestClient(planDetailsFixture, 200)

	plan, err := client.GetPlan(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := uuid.MustParse(testID1)
	if plan.ID != want {
		t.Errorf("got ID %v, want %v", plan.ID, want)
	}

	if plan.Name != "Bunko Budget" {
		t.Errorf("got Name %v, want Bunko Budget", plan.Name)
	}
}

func TestGetLastUsedPlan(t *testing.T) {
	client, transport := newTestClient(planDetailsFixture, 200)

	plan, err := client.GetLastUsedPlan(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := uuid.MustParse(testID1)
	if plan.ID != want {
		t.Errorf("got ID %v, want %v", plan.ID, want)
	}

	if !strings.HasSuffix(transport.lastReq.URL.Path, "/plans/last-used") {
		t.Errorf("unexpected path %v", transport.lastReq.URL.Path)
	}
}

func TestGetPlanSettings(t *testing.T) {
	client, _ := newTestClient(planSettingsFixture, 200)

	settings, err := client.GetPlanSettings(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if settings.DateFormat.Format != "MM/DD/YYYY" {
		t.Errorf("got DateFormat.Format %v, want MM/DD/YYYY", settings.DateFormat.Format)
	}

	if settings.CurrencyFormat.IsoCode != "USD" {
		t.Errorf("got CurrencyFormat.IsoCode %v, want USD", settings.CurrencyFormat.IsoCode)
	}
}
