//go:build integration

package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"

	"go-ynab/ynab"
)

func requireEnv(t *testing.T, key string) string {
	t.Helper()
	v := os.Getenv(key)
	if v == "" {
		t.Skipf("%s not set", key)
	}
	return v
}

func setup(t *testing.T) (*ynab.Client, uuid.UUID) {
	t.Helper()
	token := requireEnv(t, "YNAB_TOKEN")
	planIDStr := requireEnv(t, "YNAB_TEST_PLAN_ID")
	planID, err := uuid.Parse(planIDStr)
	if err != nil {
		t.Fatalf("invalid YNAB_TEST_PLAN_ID: %v", err)
	}
	return ynab.NewClient(token), planID
}

func firstAccountID(t *testing.T, client *ynab.Client, planID uuid.UUID) uuid.UUID {
	t.Helper()
	accounts, _, err := client.GetAccounts(context.Background(), planID, nil)
	if err != nil {
		t.Fatalf("GetAccounts: %v", err)
	}
	if len(accounts) == 0 {
		t.Fatal("no accounts in test plan")
	}
	return accounts[0].ID
}

func TestGetTransactions_Smoke(t *testing.T) {
	client, planID := setup(t)

	txs, sk, err := client.GetTransactions(context.Background(), planID, nil)
	if err != nil {
		t.Fatalf("GetTransactions: %v", err)
	}
	if sk <= 0 {
		t.Errorf("expected server_knowledge > 0, got %d", sk)
	}
	t.Logf("fetched %d transactions, server_knowledge=%d", len(txs), sk)
}

func TestTransaction_CreateGetDelete(t *testing.T) {
	client, planID := setup(t)
	ctx := context.Background()
	accountID := firstAccountID(t, client, planID)

	memo := "integration test transaction"
	created, err := client.CreateTransaction(ctx, planID, ynab.SaveTransaction{
		AccountID: accountID,
		Date:      ynab.NewDate(time.Now().Date()),
		Amount:    1000,
		Memo:      &memo,
		Cleared:   ynab.ClearedStatusUncleared,
	})
	if err != nil {
		t.Fatalf("CreateTransaction: %v", err)
	}
	txID := created.Transaction.ID
	t.Logf("created transaction %s", txID)
	t.Cleanup(func() {
		if _, err := client.DeleteTransaction(ctx, planID, txID); err != nil {
			t.Logf("cleanup: DeleteTransaction %s: %v", txID, err)
		}
	})

	fetched, err := client.GetTransaction(ctx, planID, txID)
	if err != nil {
		t.Fatalf("GetTransaction: %v", err)
	}
	if fetched.ID != txID {
		t.Errorf("got ID %s, want %s", fetched.ID, txID)
	}
	if fetched.Amount != 1000 {
		t.Errorf("got Amount %d, want 1000", fetched.Amount)
	}
	if fetched.Memo == nil || *fetched.Memo != memo {
		t.Errorf("got Memo %v, want %q", fetched.Memo, memo)
	}
}

func TestTransaction_UpdateSingle(t *testing.T) {
	client, planID := setup(t)
	ctx := context.Background()
	accountID := firstAccountID(t, client, planID)

	memo := "integration update test"
	created, err := client.CreateTransaction(ctx, planID, ynab.SaveTransaction{
		AccountID: accountID,
		Date:      ynab.NewDate(time.Now().Date()),
		Amount:    1000,
		Memo:      &memo,
		Cleared:   ynab.ClearedStatusUncleared,
	})
	if err != nil {
		t.Fatalf("CreateTransaction: %v", err)
	}
	txID := created.Transaction.ID
	t.Cleanup(func() {
		if _, err := client.DeleteTransaction(ctx, planID, txID); err != nil {
			t.Logf("cleanup: DeleteTransaction %s: %v", txID, err)
		}
	})

	updatedMemo := "integration update test (updated)"
	updated, err := client.UpdateTransaction(ctx, planID, txID, ynab.UpdateTransaction{
		ID:        txID,
		AccountID: accountID,
		Date:      ynab.NewDate(time.Now().Date()),
		Amount:    2000,
		Memo:      &updatedMemo,
		Cleared:   ynab.ClearedStatusUncleared,
		Approved:  &created.Transaction.Approved,
	})
	if err != nil {
		t.Fatalf("UpdateTransaction: %v", err)
	}
	if updated.Transaction.Amount != 2000 {
		t.Errorf("got Amount %d, want 2000", updated.Transaction.Amount)
	}
	if updated.Transaction.Memo == nil || *updated.Transaction.Memo != updatedMemo {
		t.Errorf("got Memo %v, want %q", updated.Transaction.Memo, updatedMemo)
	}
}

func TestTransactions_CreateBatchAndUpdateBatch(t *testing.T) {
	client, planID := setup(t)
	ctx := context.Background()
	accountID := firstAccountID(t, client, planID)

	today := ynab.NewDate(time.Now().Date())
	memo1, memo2 := "integration batch tx 1", "integration batch tx 2"

	created, err := client.CreateTransactions(ctx, planID, []ynab.SaveTransaction{
		{AccountID: accountID, Date: today, Amount: 500, Memo: &memo1, Cleared: ynab.ClearedStatusUncleared},
		{AccountID: accountID, Date: today, Amount: 750, Memo: &memo2, Cleared: ynab.ClearedStatusUncleared},
	})
	if err != nil {
		t.Fatalf("CreateTransactions: %v", err)
	}
	if len(created.TransactionIDs) != 2 {
		t.Fatalf("expected 2 transaction IDs, got %d", len(created.TransactionIDs))
	}
	t.Logf("created %d transactions", len(created.TransactionIDs))
	for _, id := range created.TransactionIDs {
		id := id
		t.Cleanup(func() {
			if _, err := client.DeleteTransaction(ctx, planID, id); err != nil {
				t.Logf("cleanup: DeleteTransaction %s: %v", id, err)
			}
		})
	}

	updates := make([]ynab.UpdateTransaction, len(created.Transactions))
	for i, tx := range created.Transactions {
		updatedMemo := fmt.Sprintf("%s (updated)", *tx.Memo)
		updates[i] = ynab.UpdateTransaction{
			ID:        tx.ID,
			AccountID: tx.AccountID,
			Date:      tx.Date,
			Amount:    tx.Amount,
			Memo:      &updatedMemo,
			Cleared:   tx.Cleared,
			Approved:  &tx.Approved,
		}
	}

	batchUpdated, err := client.UpdateTransactions(ctx, planID, updates)
	if err != nil {
		t.Fatalf("UpdateTransactions: %v", err)
	}
	if len(batchUpdated.Transactions) != 2 {
		t.Errorf("expected 2 updated transactions, got %d", len(batchUpdated.Transactions))
	}
}

func TestTransaction_SplitTransaction(t *testing.T) {
	client, planID := setup(t)
	ctx := context.Background()
	accountID := firstAccountID(t, client, planID)

	memo := "integration split transaction"
	sub1Memo, sub2Memo := "split leg 1", "split leg 2"

	created, err := client.CreateTransaction(ctx, planID, ynab.SaveTransaction{
		AccountID: accountID,
		Date:      ynab.NewDate(time.Now().Date()),
		Amount:    5000,
		Memo:      &memo,
		Cleared:   ynab.ClearedStatusUncleared,
		Subtransactions: []ynab.SaveSubtransaction{
			{Amount: 2000, Memo: &sub1Memo},
			{Amount: 3000, Memo: &sub2Memo},
		},
	})
	if err != nil {
		t.Fatalf("CreateTransaction (split): %v", err)
	}
	txID := created.Transaction.ID
	t.Cleanup(func() {
		if _, err := client.DeleteTransaction(ctx, planID, txID); err != nil {
			t.Logf("cleanup: DeleteTransaction %s: %v", txID, err)
		}
	})

	if len(created.Transaction.Subtransactions) != 2 {
		t.Errorf("expected 2 subtransactions, got %d", len(created.Transaction.Subtransactions))
	}
	var total int64
	for _, stx := range created.Transaction.Subtransactions {
		total += stx.Amount
	}
	if total != created.Transaction.Amount {
		t.Errorf("subtransaction total %d does not match parent amount %d", total, created.Transaction.Amount)
	}
}

// TestCategory_CreateAndUpdate creates a category group and category, then updates both.
//
// Note: the YNAB API has no delete endpoint for categories or groups; these will
// persist in the test plan.
func TestCategory_CreateAndUpdate(t *testing.T) {
	client, planID := setup(t)
	ctx := context.Background()

	group, err := client.CreateCategoryGroup(ctx, planID, ynab.SaveCategoryGroup{
		Name: "integration-test-group",
	})
	if err != nil {
		t.Fatalf("CreateCategoryGroup: %v", err)
	}
	t.Logf("created category group %s", group.ID)

	cat, err := client.CreateCategory(ctx, planID, ynab.SaveCategory{
		CategoryGroupID: group.ID,
		Name:            "integration-test-category",
	})
	if err != nil {
		t.Fatalf("CreateCategory: %v", err)
	}
	t.Logf("created category %s", cat.ID)

	updatedCat, err := client.UpdateCategory(ctx, planID, cat.ID, ynab.SaveCategory{
		CategoryGroupID: group.ID,
		Name:            "integration-test-category (updated)",
	})
	if err != nil {
		t.Fatalf("UpdateCategory: %v", err)
	}
	if updatedCat.Name != "integration-test-category (updated)" {
		t.Errorf("got Name %q, want %q", updatedCat.Name, "integration-test-category (updated)")
	}

	updatedGroup, err := client.UpdateCategoryGroup(ctx, planID, group.ID, ynab.SaveCategoryGroup{
		Name: "integration-test-group (updated)",
	})
	if err != nil {
		t.Fatalf("UpdateCategoryGroup: %v", err)
	}
	if updatedGroup.Name != "integration-test-group (updated)" {
		t.Errorf("got Name %q, want %q", updatedGroup.Name, "integration-test-group (updated)")
	}
}

// TestPayee_CreateAndUpdate creates and updates a payee.
//
// Note: the YNAB API has no delete endpoint for payees; these will persist
// in the test plan.
func TestPayee_CreateAndUpdate(t *testing.T) {
	client, planID := setup(t)
	ctx := context.Background()

	payee, err := client.CreatePayee(ctx, planID, ynab.PostPayee{Name: "integration-test-payee"})
	if err != nil {
		t.Fatalf("CreatePayee: %v", err)
	}
	t.Logf("created payee %s", payee.ID)

	updated, err := client.UpdatePayee(ctx, planID, payee.ID, ynab.PostPayee{Name: "integration-test-payee (updated)"})
	if err != nil {
		t.Fatalf("UpdatePayee: %v", err)
	}
	if updated.Name != "integration-test-payee (updated)" {
		t.Errorf("got Name %q, want %q", updated.Name, "integration-test-payee (updated)")
	}
}

func TestScheduledTransaction_CRUD(t *testing.T) {
	client, planID := setup(t)
	ctx := context.Background()
	accountID := firstAccountID(t, client, planID)

	memo := "integration scheduled transaction"
	created, err := client.CreateScheduledTransaction(ctx, planID, ynab.SaveScheduledTransaction{
		AccountID: accountID,
		Date:      ynab.NewDate(time.Now().Date()),
		Amount:    -5000,
		Frequency: ynab.FrequencyMonthly,
		Memo:      &memo,
	})
	if err != nil {
		t.Fatalf("CreateScheduledTransaction: %v", err)
	}
	stxID := created.ID
	t.Logf("created scheduled transaction %s", stxID)
	t.Cleanup(func() {
		if _, err := client.DeleteScheduledTransaction(ctx, planID, stxID); err != nil {
			t.Logf("cleanup: DeleteScheduledTransaction %s: %v", stxID, err)
		}
	})

	fetched, err := client.GetScheduledTransaction(ctx, planID, stxID)
	if err != nil {
		t.Fatalf("GetScheduledTransaction: %v", err)
	}
	if fetched.ID != stxID {
		t.Errorf("got ID %s, want %s", fetched.ID, stxID)
	}
	if fetched.Frequency != ynab.FrequencyMonthly {
		t.Errorf("got Frequency %s, want %s", fetched.Frequency, ynab.FrequencyMonthly)
	}

	updatedMemo := "integration scheduled transaction (updated)"
	updated, err := client.UpdateScheduledTransaction(ctx, planID, stxID, ynab.SaveScheduledTransaction{
		AccountID: accountID,
		Date:      ynab.NewDate(time.Now().Date()),
		Amount:    -7500,
		Frequency: ynab.FrequencyWeekly,
		Memo:      &updatedMemo,
	})
	if err != nil {
		t.Fatalf("UpdateScheduledTransaction: %v", err)
	}
	if updated.Amount != -7500 {
		t.Errorf("got Amount %d, want -7500", updated.Amount)
	}
	if updated.Frequency != ynab.FrequencyWeekly {
		t.Errorf("got Frequency %s, want %s", updated.Frequency, ynab.FrequencyWeekly)
	}
}
