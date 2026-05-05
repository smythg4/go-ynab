package ynab

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
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
	client, _ := newTestClient(txListFixture, 200)

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
}

func TestGetTransactionsFilteredSinceDate(t *testing.T) {
	client, transport := newTestClient(txListFixture, 200)

	d := NewDate(time.Now().Date())
	params := TransactionListParams{SinceDate: &d}
	_, _, err := client.GetTransactions(context.Background(), uuid.New(), &params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(transport.lastReq.URL.RawQuery, "since_date=") {
		t.Errorf("expected since_date in query string, got %q", transport.lastReq.URL.RawQuery)
	}
}

func TestGetTransactionsFilteredServerKnowledge(t *testing.T) {
	client, transport := newTestClient(txListFixture, 200)

	sk := int64(1)
	params := TransactionListParams{LastKnowledgeOfServer: &sk}
	_, _, err := client.GetTransactions(context.Background(), uuid.New(), &params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(transport.lastReq.URL.RawQuery, "last_knowledge_of_server=") {
		t.Errorf("expected last_knowledge_of_server in query string, got %q", transport.lastReq.URL.RawQuery)
	}
}

func TestGetTransactionsFilteredType(t *testing.T) {
	client, transport := newTestClient(txListFixture, 200)

	ttype := "unapproved"
	params := TransactionListParams{Type: &ttype}
	_, _, err := client.GetTransactions(context.Background(), uuid.New(), &params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(transport.lastReq.URL.RawQuery, "type=") {
		t.Errorf("expected type in query string, got %q", transport.lastReq.URL.RawQuery)
	}
}

func TestGetTransaction(t *testing.T) {
	client, _ := newTestClient(txSingleFixture, 200)

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
}

func TestGetTransactionsByAccount(t *testing.T) {
	client, _ := newTestClient(txListFixture, 200)

	txs, serverKnowledge, err := client.GetTransactionsByAccount(context.Background(), uuid.New(), uuid.New(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(txs) != 1 {
		t.Fatalf("expected 1 transaction, got %d", len(txs))
	}

	if serverKnowledge != 5 {
		t.Errorf("got server_knowledge %v, want 5", serverKnowledge)
	}
}

func TestGetTransactionsByCategory(t *testing.T) {
	client, _ := newTestClient(txListFixture, 200)

	txs, _, err := client.GetTransactionsByCategory(context.Background(), uuid.New(), uuid.New(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(txs) != 1 {
		t.Fatalf("expected 1 transaction, got %d", len(txs))
	}
}

func TestGetTransactionsByPayee(t *testing.T) {
	client, _ := newTestClient(txListFixture, 200)

	txs, _, err := client.GetTransactionsByPayee(context.Background(), uuid.New(), uuid.New(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(txs) != 1 {
		t.Fatalf("expected 1 transaction, got %d", len(txs))
	}
}

func TestGetTransactionsByMonth(t *testing.T) {
	client, _ := newTestClient(txListFixture, 200)

	month := Date{time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)}
	txs, _, err := client.GetTransactionsByMonth(context.Background(), uuid.New(), month, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(txs) != 1 {
		t.Fatalf("expected 1 transaction, got %d", len(txs))
	}
}

func TestGetScheduledTransactions(t *testing.T) {
	client, _ := newTestClient(scheduledTxListFixture, 200)

	txs, serverKnowledge, err := client.GetScheduledTransactions(context.Background(), uuid.New(), nil)
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
}

func TestGetScheduledTransaction(t *testing.T) {
	client, _ := newTestClient(scheduledTxSingleFixture, 200)

	tx, err := client.GetScheduledTransaction(context.Background(), uuid.New(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	idWant := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	if tx.ID != idWant {
		t.Errorf("got ID %v, want %v", tx.ID, idWant)
	}
}

func TestDeleteTransaction(t *testing.T) {
	client, transport := newTestClient(txSingleFixture, 200)

	tx, err := client.DeleteTransaction(context.Background(), uuid.New(), "abc-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if transport.lastReq.Method != http.MethodDelete {
		t.Errorf("got method %v, want DELETE", transport.lastReq.Method)
	}

	if tx.ID != "abc-123" {
		t.Errorf("got ID %v, want abc-123", tx.ID)
	}
}

func TestDeleteScheduledTransaction(t *testing.T) {
	client, transport := newTestClient(scheduledTxSingleFixture, 200)

	tx, err := client.DeleteScheduledTransaction(context.Background(), uuid.New(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if transport.lastReq.Method != http.MethodDelete {
		t.Errorf("got method %v, want DELETE", transport.lastReq.Method)
	}

	idWant := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	if tx.ID != idWant {
		t.Errorf("got ID %v, want %v", tx.ID, idWant)
	}
}

func TestUpdateTransaction(t *testing.T) {
	fixture := `{"data":{"transaction_ids":["abc-123"],"transaction":` + txFixture +
		`,"duplicate_import_ids":[],"server_knowledge":10}}`
	client, transport := newTestClient(fixture, 200)

	resp, err := client.UpdateTransaction(context.Background(), uuid.New(), "abc-123", UpdateTransaction{
		ID:        "abc-123",
		AccountID: uuid.New(),
		Date:      Date{time.Now()},
		Amount:    -15000,
		Cleared:   ClearedStatusCleared,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if transport.lastReq.Method != http.MethodPut {
		t.Errorf("got method %v, want PUT", transport.lastReq.Method)
	}

	if resp.Transaction.ID != "abc-123" {
		t.Errorf("got transaction ID %v, want abc-123", resp.Transaction.ID)
	}

	if resp.ServerKnowledge != 10 {
		t.Errorf("got server_knowledge %v, want 10", resp.ServerKnowledge)
	}

	var payload updateTransactionWrapper
	if err := json.Unmarshal(transport.lastBody, &payload); err != nil {
		t.Fatalf("could not unmarshal request body: %v", err)
	}
	if payload.Transactions.ID != "abc-123" {
		t.Errorf("got payload ID %v, want abc-123", payload.Transactions.ID)
	}
}

func TestUpdateScheduledTransaction(t *testing.T) {
	client, transport := newTestClient(scheduledTxSingleFixture, 200)

	tx, err := client.UpdateScheduledTransaction(context.Background(), uuid.New(), uuid.New(),
		SaveScheduledTransaction{
			AccountID: uuid.New(),
			Date:      Date{time.Now()},
			Amount:    -50000,
			Frequency: FrequencyMonthly,
		})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if transport.lastReq.Method != http.MethodPut {
		t.Errorf("got method %v, want PUT", transport.lastReq.Method)
	}

	idWant := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	if tx.ID != idWant {
		t.Errorf("got ID %v, want %v", tx.ID, idWant)
	}

	var payload saveScheduledTransactionWrapper
	if err := json.Unmarshal(transport.lastBody, &payload); err != nil {
		t.Fatalf("could not unmarshal request body: %v", err)
	}
	if payload.Transaction.Frequency != FrequencyMonthly {
		t.Errorf("got payload frequency %v, want %v", payload.Transaction.Frequency, FrequencyMonthly)
	}
}

func TestCreateTransaction(t *testing.T) {
	fixture := `{"data":{"transaction_ids":["abc-123"],"transaction":` + txFixture + `,"duplicate_import_ids":[],"server_knowledge":6}}`
	client, transport := newTestClient(fixture, 201)

	resp, err := client.CreateTransaction(context.Background(), uuid.New(), SaveTransaction{
		AccountID: uuid.New(),
		Date:      Date{time.Now()},
		Amount:    -15000,
		Cleared:   ClearedStatusCleared,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if transport.lastReq.Method != http.MethodPost {
		t.Errorf("got method %v, want POST", transport.lastReq.Method)
	}

	if resp.Transaction.ID != "abc-123" {
		t.Errorf("got transaction ID %v, want abc-123", resp.Transaction.ID)
	}

	if resp.ServerKnowledge != 6 {
		t.Errorf("got server_knowledge %v, want 6", resp.ServerKnowledge)
	}

	var payload saveTransactionWrapper
	if err := json.Unmarshal(transport.lastBody, &payload); err != nil {
		t.Fatalf("could not unmarshal request body: %v", err)
	}
	if payload.Transaction.Amount != -15000 {
		t.Errorf("got payload amount %v, want -15000", payload.Transaction.Amount)
	}
}

func TestCreateTransactions(t *testing.T) {
	fixture := `{"data":{"transaction_ids":["abc-123","abc-456"],"transactions":[` + txFixture + `,` + txFixture + `],"duplicate_import_ids":[],"server_knowledge":7}}`
	client, transport := newTestClient(fixture, 201)

	resp, err := client.CreateTransactions(context.Background(), uuid.New(), []SaveTransaction{
		{AccountID: uuid.New(), Date: Date{time.Now()}, Amount: -15000, Cleared: ClearedStatusCleared},
		{AccountID: uuid.New(), Date: Date{time.Now()}, Amount: -25000, Cleared: ClearedStatusUncleared},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if transport.lastReq.Method != http.MethodPost {
		t.Errorf("got method %v, want POST", transport.lastReq.Method)
	}

	if len(resp.TransactionIDs) != 2 {
		t.Fatalf("expected 2 transaction IDs, got %d", len(resp.TransactionIDs))
	}

	if resp.ServerKnowledge != 7 {
		t.Errorf("got server_knowledge %v, want 7", resp.ServerKnowledge)
	}

	var payload saveTransactionsWrapper
	if err := json.Unmarshal(transport.lastBody, &payload); err != nil {
		t.Fatalf("could not unmarshal request body: %v", err)
	}
	if len(payload.Transactions) != 2 {
		t.Errorf("got %d transactions in payload, want 2", len(payload.Transactions))
	}
}

func TestImportTransactions(t *testing.T) {
	fixture := `{"data":{"transaction_ids":["abc-123"],"server_knowledge":8}}`
	client, transport := newTestClient(fixture, 200)

	resp, err := client.ImportTransactions(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if transport.lastReq.Method != http.MethodPost {
		t.Errorf("got method %v, want POST", transport.lastReq.Method)
	}

	if len(resp.TransactionIDs) != 1 {
		t.Fatalf("expected 1 transaction ID, got %d", len(resp.TransactionIDs))
	}

	if resp.ServerKnowledge != 8 {
		t.Errorf("got server_knowledge %v, want 8", resp.ServerKnowledge)
	}
}

func TestCreateScheduledTransaction(t *testing.T) {
	client, transport := newTestClient(scheduledTxSingleFixture, 201)

	tx, err := client.CreateScheduledTransaction(context.Background(), uuid.New(), SaveScheduledTransaction{
		AccountID: uuid.New(),
		Date:      Date{time.Now()},
		Amount:    -50000,
		Frequency: FrequencyMonthly,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if transport.lastReq.Method != http.MethodPost {
		t.Errorf("got method %v, want POST", transport.lastReq.Method)
	}

	idWant := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	if tx.ID != idWant {
		t.Errorf("got ID %v, want %v", tx.ID, idWant)
	}

	var payload saveScheduledTransactionWrapper
	if err := json.Unmarshal(transport.lastBody, &payload); err != nil {
		t.Fatalf("could not unmarshal request body: %v", err)
	}
	if payload.Transaction.Frequency != FrequencyMonthly {
		t.Errorf("got payload frequency %v, want %v", payload.Transaction.Frequency, FrequencyMonthly)
	}
}

func TestUpdateTransactions(t *testing.T) {
	fixture := `{"data":{"transaction_ids":["abc-123"],"transactions":[` + txFixture + `],"duplicate_import_ids":[],"server_knowledge":9}}`
	client, transport := newTestClient(fixture, 200)

	resp, err := client.UpdateTransactions(context.Background(), uuid.New(), []UpdateTransaction{
		{ID: "abc-123", AccountID: uuid.New(), Date: Date{time.Now()}, Amount: -15000, Cleared: ClearedStatusCleared},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if transport.lastReq.Method != http.MethodPatch {
		t.Errorf("got method %v, want PATCH", transport.lastReq.Method)
	}

	if len(resp.Transactions) != 1 {
		t.Fatalf("expected 1 transaction, got %d", len(resp.Transactions))
	}

	if resp.ServerKnowledge != 9 {
		t.Errorf("got server_knowledge %v, want 9", resp.ServerKnowledge)
	}

	var payload updateTransactionsWrapper
	if err := json.Unmarshal(transport.lastBody, &payload); err != nil {
		t.Fatalf("could not unmarshal request body: %v", err)
	}
	if len(payload.Transactions) != 1 {
		t.Errorf("got %d transactions in payload, want 1", len(payload.Transactions))
	}
}
