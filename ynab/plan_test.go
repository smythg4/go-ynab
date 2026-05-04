package ynab

import (
	"context"
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
	t.Run("returns plan list on success", func(t *testing.T) {
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

		nameWant := "Bunko Budget"
		if plans[0].Name != nameWant {
			t.Errorf("got Name %v, want %v", plans[0].Name, nameWant)
		}
	})
}

func TestGetPlan(t *testing.T) {
	t.Run("returns plan details on success", func(t *testing.T) {
		client, _ := newTestClient(planDetailsFixture, 200)

		plan, err := client.GetPlan(context.Background(), uuid.New())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := uuid.MustParse(testID1)
		if plan.ID != want {
			t.Errorf("got ID %v, want %v", plan.ID, want)
		}

		nameWant := "Bunko Budget"
		if plan.Name != nameWant {
			t.Errorf("got Name %v, want %v", plan.Name, nameWant)
		}
	})
}

func TestGetPlanSettings(t *testing.T) {
	t.Run("returns plan settings on success", func(t *testing.T) {
		client, _ := newTestClient(planSettingsFixture, 200)

		settings, err := client.GetPlanSettings(context.Background(), uuid.New())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		dateFormatWant := "MM/DD/YYYY"
		if settings.DateFormat.Format != dateFormatWant {
			t.Errorf("got DateFormat.Format %v, want %v", settings.DateFormat.Format, dateFormatWant)
		}

		isoCodeWant := "USD"
		if settings.CurrencyFormat.IsoCode != isoCodeWant {
			t.Errorf("got CurrencyFormat.IsoCode %v, want %v", settings.CurrencyFormat.IsoCode, isoCodeWant)
		}
	})
}
