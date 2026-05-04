// create-transaction creates a single transaction against the first account
// in the first plan for the authenticated user, then prints the result.
//
// Usage:
//
//	YNAB_TOKEN=your_token go run ./examples/create-transaction
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

	memo := "dummy transaction"
	tx := ynab.SaveTransaction{
		AccountID: accountId,
		Date:      ynab.NewDate(time.Now().Date()),
		Amount:    1000,
		Memo:      &memo,
		Cleared:   ynab.ClearedStatusUncleared,
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
	fmt.Printf("   %-10s %s\n", "Memo:", *resp.Transaction.Memo)
}
