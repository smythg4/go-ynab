package ynab

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

const accountFixture = `{"id":"123e4567-e89b-12d3-a456-426614174000","name":"Checking","type":"checking","on_budget":true,"closed":false,"note":null,"balance":100000,"cleared_balance":95000,"uncleared_balance":5000,"transfer_payee_id":null,"direct_import_linked":false,"direct_import_in_error":false,"last_reconciled_at":null,"deleted":false}`

const accountListFixture = `{"data":{"accounts":[` + accountFixture + `],"server_knowledge":7}}`
const accountSingleFixture = `{"data":{"account":` + accountFixture + `}}`

func TestGetAccounts(t *testing.T) {
	t.Run("returns account list on success", func(t *testing.T) {
		client := newTestClient(accountListFixture, 200)

		accounts, serverKnowledge, err := client.GetAccounts(context.Background(), uuid.New(), nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(accounts) != 1 {
			t.Fatalf("expected 1 account, got %d", len(accounts))
		}

		idWant := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
		if accounts[0].ID != idWant {
			t.Errorf("got ID %v, want %v", accounts[0].ID, idWant)
		}

		if accounts[0].Type != AccountTypeChecking {
			t.Errorf("got Type %v, want %v", accounts[0].Type, AccountTypeChecking)
		}

		if accounts[0].Balance != 100000 {
			t.Errorf("got Balance %v, want 100000", accounts[0].Balance)
		}

		if serverKnowledge != 7 {
			t.Errorf("got server_knowledge %v, want 7", serverKnowledge)
		}
	})
}

func TestGetAccount(t *testing.T) {
	t.Run("returns single account on success", func(t *testing.T) {
		client := newTestClient(accountSingleFixture, 200)

		account, err := client.GetAccount(context.Background(), uuid.New(), uuid.New())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		idWant := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
		if account.ID != idWant {
			t.Errorf("got ID %v, want %v", account.ID, idWant)
		}

		nameWant := "Checking"
		if account.Name != nameWant {
			t.Errorf("got Name %v, want %v", account.Name, nameWant)
		}
	})
}
