// update-transaction fetches the most recent transaction from the first plan
// for the authenticated user, appends " (updated)" to its memo, and replaces
// it via a PUT request, then prints the before and after state.
//
// UpdateTransaction is a full replacement — all fields must be provided.
// Fetch the existing transaction first to avoid losing data.
//
// Note: running this example multiple times will append " (updated)" repeatedly
// to the memo of the same transaction.
//
// Usage:
//
//	YNAB_TOKEN=your_token go run ./examples/update-transaction
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go-ynab/ynab"
)

func main() {
	token := os.Getenv("YNAB_TOKEN")
	if token == "" {
		log.Fatal("YNAB_TOKEN environment variable is not set")
	}

	client := ynab.NewClient(token)
	ctx := context.Background()

	plans, err := client.GetPlans(ctx)
	if err != nil {
		log.Fatalf("failed to get plans: %v", err)
	}
	if len(plans) == 0 {
		fmt.Println("no plans found")
		return
	}
	planId := plans[0].ID

	txs, _, err := client.GetTransactions(ctx, planId, nil)
	if err != nil {
		log.Fatalf("failed to get transactions: %v", err)
	}
	if len(txs) == 0 {
		fmt.Println("no transactions found")
		return
	}

	original := txs[len(txs)-1]

	oldMemo := ""
	if original.Memo != nil {
		oldMemo = *original.Memo
	}
	newMemo := oldMemo + " (updated)"

	update := ynab.UpdateTransaction{
		ID:         original.ID,
		AccountID:  original.AccountID,
		Date:       original.Date,
		Amount:     original.Amount,
		PayeeID:    original.PayeeID,
		CategoryID: original.CategoryID,
		Cleared:    original.Cleared,
		Approved:   &original.Approved,
		FlagColor:  original.FlagColor,
		Memo:       &newMemo,
	}

	resp, err := client.UpdateTransaction(ctx, planId, original.ID, update)
	if err != nil {
		log.Fatalf("failed to update transaction: %v", err)
	}

	fmt.Printf("Updated Transaction\n\n")
	fmt.Printf("Plan: %s\n\n", plans[0].Name)
	fmt.Printf("   %-10s %s\n", "ID:", resp.Transaction.ID)
	fmt.Printf("   %-10s %s\n", "Account:", resp.Transaction.AccountName)
	fmt.Printf("   %-10s %s\n", "Date:", resp.Transaction.Date)
	fmt.Printf("   %-10s $%.2f\n", "Amount:", ynab.MilliunitsToAmount(resp.Transaction.Amount))
	fmt.Printf("   %-10s %q  ->  %q\n", "Memo:", oldMemo, newMemo)
}
