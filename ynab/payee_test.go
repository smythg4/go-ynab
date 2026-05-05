package ynab

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

const payeeFixture = `{"id":"` + testID1 + `","name":"Testing Tom","transfer_account_id":null,"deleted":false}`
const payeeListFixture = `{"data":{"payees":[` + payeeFixture + `],"server_knowledge":1}}`
const payeeSingleFixture = `{"data":{"payee":` + payeeFixture + `}}`
const payeeLocationFixture = `{"id":"` + testID2 + `","payee_id":"` + testID1 + `","latitude":"40.7128","longitude":"-74.0060","deleted":false}`
const payeeLocationListFixture = `{"data":{"payee_locations":[` + payeeLocationFixture + `]}}`
const payeeLocationSingleFixture = `{"data":{"payee_location":` + payeeLocationFixture + `}}`

func TestGetPayees(t *testing.T) {
	client, _ := newTestClient(payeeListFixture, 200)

	payees, serverKnowledge, err := client.GetPayees(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(payees) != 1 {
		t.Fatalf("expected 1 payee, got %d", len(payees))
	}

	idWant := uuid.MustParse(testID1)
	if payees[0].ID != idWant {
		t.Errorf("got ID %v, want %v", payees[0].ID, idWant)
	}

	if payees[0].Name != "Testing Tom" {
		t.Errorf("got Name %v, want Testing Tom", payees[0].Name)
	}

	if serverKnowledge != 1 {
		t.Errorf("got server_knowledge %v, want 1", serverKnowledge)
	}
}

func TestGetPayee(t *testing.T) {
	client, _ := newTestClient(payeeSingleFixture, 200)

	payee, err := client.GetPayee(context.Background(), uuid.New(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := uuid.MustParse(testID1)
	if payee.ID != want {
		t.Errorf("got ID %v, want %v", payee.ID, want)
	}
}

func TestGetPayeeLocations(t *testing.T) {
	client, _ := newTestClient(payeeLocationListFixture, 200)

	locations, err := client.GetPayeeLocations(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(locations) != 1 {
		t.Fatalf("expected 1 location, got %d", len(locations))
	}

	want := uuid.MustParse(testID2)
	if locations[0].ID != want {
		t.Errorf("got ID %v, want %v", locations[0].ID, want)
	}
}

func TestGetPayeeLocation(t *testing.T) {
	client, _ := newTestClient(payeeLocationSingleFixture, 200)

	location, err := client.GetPayeeLocation(context.Background(), uuid.New(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := uuid.MustParse(testID2)
	if location.ID != want {
		t.Errorf("got ID %v, want %v", location.ID, want)
	}
}

func TestGetPayeeLocationsByPayee(t *testing.T) {
	client, _ := newTestClient(payeeLocationListFixture, 200)

	locations, err := client.GetPayeeLocationsByPayee(context.Background(), uuid.New(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(locations) != 1 {
		t.Fatalf("expected 1 location, got %d", len(locations))
	}

	want := uuid.MustParse(testID2)
	if locations[0].ID != want {
		t.Errorf("got ID %v, want %v", locations[0].ID, want)
	}
}

func TestCreatePayee(t *testing.T) {
	client, transport := newTestClient(payeeSingleFixture, 201)

	payee, err := client.CreatePayee(context.Background(), uuid.New(), PostPayee{Name: "Testing Tom"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if transport.lastReq.Method != http.MethodPost {
		t.Errorf("got method %v, want POST", transport.lastReq.Method)
	}

	idWant := uuid.MustParse(testID1)
	if payee.ID != idWant {
		t.Errorf("got ID %v, want %v", payee.ID, idWant)
	}

	var payload postPayeeWrapper
	if err := json.Unmarshal(transport.lastBody, &payload); err != nil {
		t.Fatalf("could not unmarshal request body: %v", err)
	}
	if payload.Payee.Name != "Testing Tom" {
		t.Errorf("got payload name %v, want Testing Tom", payload.Payee.Name)
	}
}

func TestUpdatePayee(t *testing.T) {
	client, transport := newTestClient(payeeSingleFixture, 200)

	payee, err := client.UpdatePayee(context.Background(), uuid.New(), uuid.New(), PostPayee{Name: "Testing Tom"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if transport.lastReq.Method != http.MethodPatch {
		t.Errorf("got method %v, want PATCH", transport.lastReq.Method)
	}

	idWant := uuid.MustParse(testID1)
	if payee.ID != idWant {
		t.Errorf("got ID %v, want %v", payee.ID, idWant)
	}

	var payload postPayeeWrapper
	if err := json.Unmarshal(transport.lastBody, &payload); err != nil {
		t.Fatalf("could not unmarshal request body: %v", err)
	}
	if payload.Payee.Name != "Testing Tom" {
		t.Errorf("got payload name %v, want Testing Tom", payload.Payee.Name)
	}
}
