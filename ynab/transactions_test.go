package ynab

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

const txFixture = `{"id":"abc-123","date":"2024-03-15","amount":-15000,"memo":null,"cleared":"cleared","approved":true,"flag_color":null,"flag_name":null,"account_id":"123e4567-e89b-12d3-a456-426614174000","payee_id":null,"account_name":"Checking","payee_name":null,"category_id":null,"category_name":null,"matched_transaction_id":null,"subtransactions":[]}`

const txListFixture = `{"data":{"transactions":[` + txFixture + `],"server_knowledge":5}}`
const txSingleFixture = `{"data":{"transaction":` + txFixture + `}}`

const scheduledTxFixture = `{"id":"123e4567-e89b-12d3-a456-426614174000","date_first":"2024-01-01","date_next":"2024-04-01","frequency":"monthly","amount":-50000,"memo":null,"flag_color":null,"flag_name":null,"account_id":"223e4567-e89b-12d3-a456-426614174000","payee_id":null,"category_id":null,"account_name":"Checking","payee_name":null,"category_name":null,"subtransactions":[],"transfer_account_id":null}`

const scheduledTxListFixture = `{"data":{"scheduled_transactions":[` + scheduledTxFixture + `],"server_knowledge":3}}`
const scheduledTxSingleFixture = `{"data":{"scheduled_transaction":` + scheduledTxFixture + `}}`

func TestGetTransactions(t *testing.T) {
	t.Run("returns transaction list on success", func(t *testing.T) {
		client := newTestClient(txListFixture, 200)

		txs, serverKnowledge, err := client.GetTransactions(context.Background(), uuid.New(), nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(txs) != 1 {
			t.Fatalf("expected 1 transaction, got %d", len(txs))
		}

		if txs[0].ID != "abc-123" {
			t.Errorf("got ID %v, want abc-123", txs[0].ID)
		}

		if txs[0].Amount != -15000 {
			t.Errorf("got Amount %v, want -15000", txs[0].Amount)
		}

		if serverKnowledge != 5 {
			t.Errorf("got server_knowledge %v, want 5", serverKnowledge)
		}
	})
}

func TestGetTransaction(t *testing.T) {
	t.Run("returns single transaction on success", func(t *testing.T) {
		client := newTestClient(txSingleFixture, 200)

		tx, err := client.GetTransaction(context.Background(), uuid.New(), "abc-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if tx.ID != "abc-123" {
			t.Errorf("got ID %v, want abc-123", tx.ID)
		}

		if tx.Cleared != ClearedStatusCleared {
			t.Errorf("got Cleared %v, want %v", tx.Cleared, ClearedStatusCleared)
		}
	})
}

func TestGetTransactionsByAccount(t *testing.T) {
	t.Run("returns transactions for account on success", func(t *testing.T) {
		client := newTestClient(txListFixture, 200)

		txs, serverKnowledge, err := client.GetTransactionsByAccount(context.Background(), uuid.New(), uuid.New())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(txs) != 1 {
			t.Fatalf("expected 1 transaction, got %d", len(txs))
		}

		if serverKnowledge != 5 {
			t.Errorf("got server_knowledge %v, want 5", serverKnowledge)
		}
	})
}

func TestGetTransactionsByCategory(t *testing.T) {
	t.Run("returns transactions for category on success", func(t *testing.T) {
		client := newTestClient(txListFixture, 200)

		txs, _, err := client.GetTransactionsByCategory(context.Background(), uuid.New(), uuid.New())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(txs) != 1 {
			t.Fatalf("expected 1 transaction, got %d", len(txs))
		}
	})
}

func TestGetTransactionsByPayee(t *testing.T) {
	t.Run("returns transactions for payee on success", func(t *testing.T) {
		client := newTestClient(txListFixture, 200)

		txs, _, err := client.GetTransactionsByPayee(context.Background(), uuid.New(), uuid.New())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(txs) != 1 {
			t.Fatalf("expected 1 transaction, got %d", len(txs))
		}
	})
}

func TestGetTransactionsByMonth(t *testing.T) {
	t.Run("returns transactions for month on success", func(t *testing.T) {
		client := newTestClient(txListFixture, 200)

		month := Date{time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)}
		txs, _, err := client.GetTransactionsByMonth(context.Background(), uuid.New(), month)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(txs) != 1 {
			t.Fatalf("expected 1 transaction, got %d", len(txs))
		}
	})
}

func TestGetScheduledTransactions(t *testing.T) {
	t.Run("returns scheduled transaction list on success", func(t *testing.T) {
		client := newTestClient(scheduledTxListFixture, 200)

		txs, serverKnowledge, err := client.GetScheduledTransactions(context.Background(), uuid.New())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(txs) != 1 {
			t.Fatalf("expected 1 scheduled transaction, got %d", len(txs))
		}

		idWant := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
		if txs[0].ID != idWant {
			t.Errorf("got ID %v, want %v", txs[0].ID, idWant)
		}

		if txs[0].Frequency != FrequencyMonthly {
			t.Errorf("got Frequency %v, want %v", txs[0].Frequency, FrequencyMonthly)
		}

		if serverKnowledge != 3 {
			t.Errorf("got server_knowledge %v, want 3", serverKnowledge)
		}
	})
}

func TestGetScheduledTransaction(t *testing.T) {
	t.Run("returns single scheduled transaction on success", func(t *testing.T) {
		client := newTestClient(scheduledTxSingleFixture, 200)

		tx, err := client.GetScheduledTransaction(context.Background(), uuid.New(), uuid.New())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		idWant := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
		if tx.ID != idWant {
			t.Errorf("got ID %v, want %v", tx.ID, idWant)
		}
	})
}
