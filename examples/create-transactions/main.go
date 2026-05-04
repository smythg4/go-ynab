// create-transactions creates multiple transactions in a single request against
// the first account in the first plan for the authenticated user, then prints
// a summary of each created transaction.
//
// Usage:
//
//	YNAB_TOKEN=your_token go run ./examples/create-transactions
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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

	accounts, _, err := client.GetAccounts(ctx, planId, nil)
	if err != nil {
		log.Fatalf("failed to get accounts: %v", err)
	}
	if len(accounts) == 0 {
		fmt.Println("no accounts found")
		return
	}
	accountId := accounts[0].ID

	today := ynab.NewDate(time.Now().Date())
	memo1 := "dummy transaction 1"
	memo2 := "dummy transaction 2"
	memo3 := "dummy transaction 3"

	txs := []ynab.SaveTransaction{
		{
			AccountID: accountId,
			Date:      today,
			Amount:    1000,
			Memo:      &memo1,
			Cleared:   ynab.ClearedStatusUncleared,
		},
		{
			AccountID: accountId,
			Date:      today,
			Amount:    2500,
			Memo:      &memo2,
			Cleared:   ynab.ClearedStatusUncleared,
		},
		{
			AccountID: accountId,
			Date:      today,
			Amount:    -500,
			Memo:      &memo3,
			Cleared:   ynab.ClearedStatusUncleared,
		},
	}

	resp, err := client.CreateTransactions(ctx, planId, txs)
	if err != nil {
		log.Fatalf("failed to create transactions: %v", err)
	}

	fmt.Printf("Created %d transactions\n\n", len(resp.Transactions))
	fmt.Printf("Plan: %s\n", plans[0].Name)
	fmt.Printf("\n   %-12s  %-10s  %s\n", "Amount", "Date", "Memo")
	fmt.Printf("   %-12s  %-10s  %s\n", "------------", "----------", "----")
	for _, tx := range resp.Transactions {
		memo := ""
		if tx.Memo != nil {
			memo = *tx.Memo
		}
		fmt.Printf("   $%-11.2f  %-10s  %s\n",
			ynab.MilliunitsToAmount(tx.Amount),
			tx.Date,
			memo,
		)
	}

	if len(resp.DuplicateImportIDs) > 0 {
		fmt.Printf("\n%d duplicate(s) skipped\n", len(resp.DuplicateImportIDs))
	}
}
