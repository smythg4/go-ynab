package ynab

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

const payeeSingleFixture = `{"data":{"payee":{"id":"123e4567-e89b-12d3-a456-426614174000","name":"Testing Tom","transfer_account_id":null,"deleted":false}}}`

func TestGetPayees(t *testing.T) {
	t.Run("returns payee list on success", func(t *testing.T) {
		fixture := `{"data": {"payees": [{"id": "123e4567-e89b-12d3-a456-426614174000","name": "Testing Tom","transfer_account_id": null,"deleted": false}],"server_knowledge": 1}}`
		client, _ := newTestClient(fixture, 200)

		planId := uuid.New()
		payees, serverKnowledge, err := client.GetPayees(context.Background(), planId)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(payees) != 1 {
			t.Fatalf("expected 1 payee, got %d", len(payees))
		}

		idWant := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
		if payees[0].ID != idWant {
			t.Errorf("got ID %v, want %v", payees[0].ID, idWant)
		}

		nameWant := "Testing Tom"
		if payees[0].Name != nameWant {
			t.Errorf("got Name %v, want %v", payees[0].Name, nameWant)
		}

		if serverKnowledge != 1 {
			t.Errorf("got server_knowledge %v, want %v", serverKnowledge, 1)
		}
	})
}

func TestGetPayee(t *testing.T) {
	t.Run("returns single payee on success", func(t *testing.T) {
		fixture := `{"data": {"payee": {"id": "123e4567-e89b-12d3-a456-426614174000","name": "Testing Tom","transfer_account_id": null,"deleted": false}}}`
		client, _ := newTestClient(fixture, 200)

		payee, err := client.GetPayee(context.Background(), uuid.New(), uuid.New())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
		if payee.ID != want {
			t.Errorf("got ID %v, want %v", payee.ID, want)
		}
	})
}

func TestGetPayeeLocations(t *testing.T) {
	t.Run("returns payee location list on success", func(t *testing.T) {
		fixture := `{"data": {"payee_locations": [{"id": "223e4567-e89b-12d3-a456-426614174000","payee_id": "123e4567-e89b-12d3-a456-426614174000","latitude": "40.7128","longitude": "-74.0060","deleted": false}]}}`
		client, _ := newTestClient(fixture, 200)

		locations, err := client.GetPayeeLocations(context.Background(), uuid.New())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(locations) != 1 {
			t.Fatalf("expected 1 location, got %d", len(locations))
		}

		want := uuid.MustParse("223e4567-e89b-12d3-a456-426614174000")
		if locations[0].ID != want {
			t.Errorf("got ID %v, want %v", locations[0].ID, want)
		}
	})
}

func TestGetPayeeLocation(t *testing.T) {
	t.Run("returns single payee location on success", func(t *testing.T) {
		fixture := `{"data": {"payee_location": {"id": "223e4567-e89b-12d3-a456-426614174000","payee_id": "123e4567-e89b-12d3-a456-426614174000","latitude": "40.7128","longitude": "-74.0060","deleted": false}}}`
		client, _ := newTestClient(fixture, 200)

		location, err := client.GetPayeeLocation(context.Background(), uuid.New(), uuid.New())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := uuid.MustParse("223e4567-e89b-12d3-a456-426614174000")
		if location.ID != want {
			t.Errorf("got ID %v, want %v", location.ID, want)
		}
	})
}

func TestGetPayeeLocationsByPayee(t *testing.T) {
	t.Run("returns payee locations for payee on success", func(t *testing.T) {
		fixture := `{"data": {"payee_locations": [{"id": "223e4567-e89b-12d3-a456-426614174000","payee_id": "123e4567-e89b-12d3-a456-426614174000","latitude": "40.7128","longitude": "-74.0060","deleted": false}]}}`
		client, _ := newTestClient(fixture, 200)

		locations, err := client.GetPayeeLocationsByPayee(context.Background(), uuid.New(), uuid.New())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(locations) != 1 {
			t.Fatalf("expected 1 location, got %d", len(locations))
		}

		want := uuid.MustParse("223e4567-e89b-12d3-a456-426614174000")
		if locations[0].ID != want {
			t.Errorf("got ID %v, want %v", locations[0].ID, want)
		}
	})
}

func TestCreatePayee(t *testing.T) {
	t.Run("sends POST and returns payee on success", func(t *testing.T) {
		client, transport := newTestClient(payeeSingleFixture, 201)

		payee, err := client.CreatePayee(context.Background(), uuid.New(), PostPayee{Name: "Testing Tom"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if transport.lastReq.Method != http.MethodPost {
			t.Errorf("got method %v, want POST", transport.lastReq.Method)
		}

		idWant := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
		if payee.ID != idWant {
			t.Errorf("got ID %v, want %v", payee.ID, idWant)
		}

		var payload PostPayeeWrapper
		if err := json.Unmarshal(transport.lastBody, &payload); err != nil {
			t.Fatalf("could not unmarshal request body: %v", err)
		}
		if payload.Payee.Name != "Testing Tom" {
			t.Errorf("got payload name %v, want Testing Tom", payload.Payee.Name)
		}
	})
}

func TestUpdatePayee(t *testing.T) {
	t.Run("sends PATCH and returns payee on success", func(t *testing.T) {
		client, transport := newTestClient(payeeSingleFixture, 200)

		payee, err := client.UpdatePayee(context.Background(), uuid.New(), uuid.New(), PostPayee{Name: "Testing Tom"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if transport.lastReq.Method != http.MethodPatch {
			t.Errorf("got method %v, want PATCH", transport.lastReq.Method)
		}

		idWant := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
		if payee.ID != idWant {
			t.Errorf("got ID %v, want %v", payee.ID, idWant)
		}

		var payload PostPayeeWrapper
		if err := json.Unmarshal(transport.lastBody, &payload); err != nil {
			t.Fatalf("could not unmarshal request body: %v", err)
		}
		if payload.Payee.Name != "Testing Tom" {
			t.Errorf("got payload name %v, want Testing Tom", payload.Payee.Name)
		}
	})
}
