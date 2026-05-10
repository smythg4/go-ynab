// update-transactions fetches the most recent 2 transactions from the first plan
// for the authenticated user, appends " (updated)" to their memos via a PATCH
// request, then prints the before and after state.
//
// Only the memo field is sent; all other fields are left unchanged by the API.
//
// Note: running this example multiple times will append " (updated)" repeatedly
// to the memo of the same transaction.
//
// Usage:
//
//	YNAB_TOKEN=your_token go run ./examples/update-transactions
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/smythg4/go-ynab/ynab"
)

func main() {
	token := os.Getenv("YNAB_TOKEN")
	if token == "" {
		log.Fatal("YNAB_TOKEN environment variable is not set")
	}

	client := ynab.NewClient(token)
	ctx := context.Background()

	plans, err := client.GetPlans(ctx, false)
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
	if len(txs) < 2 {
		fmt.Println("not enough transactions found")
		return
	}

	originals := txs[len(txs)-2:]
	updates := make([]ynab.PatchTransaction, 0, 2)
	oldmemos := make(map[uuid.UUID]string, 2)
	for _, original := range originals {
		id := original.ID
		oldMemo := ""
		if original.Memo != nil {
			oldMemo = *original.Memo
		}
		newMemo := oldMemo + " (updated)"

		updates = append(updates, ynab.PatchTransaction{
			ID:   &id,
			Memo: &newMemo,
		})
		oldmemos[original.ID] = oldMemo
	}

	resp, err := client.UpdateTransactions(ctx, planId, updates)
	if err != nil {
		log.Fatalf("failed to update transactions: %v", err)
	}

	fmt.Printf("Updated Transactions\n\n")
	fmt.Printf("Plan: %s\n\n", plans[0].Name)
	for _, tx := range resp.Transactions {

		fmt.Printf("   %-10s %s\n", "ID:", tx.ID)
		fmt.Printf("   %-10s %s\n", "Account:", tx.AccountName)
		fmt.Printf("   %-10s %s\n", "Date:", tx.Date)
		fmt.Printf("   %-10s $%.2f\n", "Amount:", ynab.MilliunitsToAmount(tx.Amount))
		fmt.Printf("   %-10s %q  ->  %q\n", "Memo:", oldmemos[tx.ID], *tx.Memo)
		fmt.Println()
	}

}
