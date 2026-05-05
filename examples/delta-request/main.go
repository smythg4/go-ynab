// delta-request fetches all the transactions for the first plan
// for the authenticated user, it then creates a transaction, then fetches
// all transactions again using the `ServerLastKnowledge` search parameter
// the response should only return that newly created transaction
//
// Usage:
//
//	YNAB_TOKEN=your_token go run ./examples/delta-request
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/smythg4/go-ynab/examples/internal/display"
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

	transactions, sk, err := client.GetTransactions(ctx, planId, nil)
	if err != nil {
		log.Fatalf("failed to get transactions: %v", err)
	}
	if len(transactions) == 0 {
		fmt.Println("no transactions found")
		return
	}

	fmt.Printf("Initial fetch (server knowledge: %d)\n", sk)
	fmt.Printf("Plan: %s  (%d transactions)\n\n", plans[0].Name, len(transactions))
	display.PrintTransactionTable(transactions)

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

	tx := ynab.SaveTransaction{
		AccountID: accountId,
		Date:      today,
		Amount:    1000,
		Memo:      &memo1,
		Cleared:   ynab.ClearedStatusUncleared,
	}

	resp, err := client.CreateTransaction(ctx, planId, tx)
	if err != nil {
		log.Fatalf("failed to create transaction: %v", err)
	}

	fmt.Printf("\nCreated transaction: %s\n\n", resp.Transaction.ID)

	params := ynab.TransactionListParams{
		LastKnowledgeOfServer: &sk,
	}
	prevSk := sk
	transactions, sk, err = client.GetTransactions(ctx, planId, &params)
	if err != nil {
		log.Fatalf("failed to get transactions: %v", err)
	}

	fmt.Printf("Delta fetch (since knowledge: %d)\n", prevSk)
	fmt.Printf("Plan: %s  (%d transactions)\n\n", plans[0].Name, len(transactions))

	if len(transactions) == 0 {
		fmt.Printf("No changes since server knowledge %d\n", prevSk)
		return
	}
	display.PrintTransactionTable(transactions)
}
