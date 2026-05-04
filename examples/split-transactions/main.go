// split-transaction creates a single transaction, broken up into two
// sub-transactions against the first account in the first plan for
// the authenticated user, then prints the result.
//
// Note: The sum value of the sub-transactions must equal the total
// transaction value.
//
// Usage:
//
//	YNAB_TOKEN=your_token go run ./examples/split-transactions
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/smythg4/go-ynab/ynab"
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

	memo := "dummy transaction"
	submemo1 := "subtran 1 " + memo
	stx1 := ynab.SaveSubtransaction{
		Amount: 2000,
		Memo:   &submemo1,
	}
	submemo2 := "subtran 2 " + memo
	stx2 := ynab.SaveSubtransaction{
		Amount: 3000,
		Memo:   &submemo2,
	}

	tx := ynab.SaveTransaction{
		AccountID:       accountId,
		Date:            ynab.NewDate(time.Now().Date()),
		Amount:          5000,
		Memo:            &memo,
		Cleared:         ynab.ClearedStatusUncleared,
		Subtransactions: []ynab.SaveSubtransaction{stx1, stx2},
	}
	resp, err := client.CreateTransaction(ctx, planId, tx)
	if err != nil {
		log.Fatalf("failed to create transaction: %v", err)
	}
	fmt.Println("Created Transaction")
	fmt.Println()
	fmt.Printf("Plan: %s\n", plans[0].Name)
	fmt.Printf("   %-10s %s\n", "ID:", resp.Transaction.ID)
	fmt.Printf("   %-10s %s\n", "Account:", resp.Transaction.AccountName)
	fmt.Printf("   %-10s %s\n", "Date:", resp.Transaction.Date)
	fmt.Printf("   %-10s $%.2f\n", "Amount:", ynab.MilliunitsToAmount(resp.Transaction.Amount))
	displayMemo := ""
	if resp.Transaction.Memo != nil {
		displayMemo = *resp.Transaction.Memo
	}
	fmt.Printf("   %-10s %s\n", "Memo:", displayMemo)
	fmt.Printf("\n   Subtransactions (%d)\n", len(resp.Transaction.Subtransactions))
	for _, stx := range resp.Transaction.Subtransactions {
		stxMemo := ""
		if stx.Memo != nil {
			stxMemo = *stx.Memo
		}
		fmt.Printf("      %-10s $%.2f\n", "Amount:", ynab.MilliunitsToAmount(stx.Amount))
		fmt.Printf("      %-10s %s\n\n", "Memo:", stxMemo)
	}
}
