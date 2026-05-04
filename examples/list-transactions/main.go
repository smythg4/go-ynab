// list-transactions retrieves all transactions for the most recently
// modified plan and prints a summary table.
//
// Usage:
//
//	YNAB_TOKEN=your_token go run ./examples/list-transactions
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go-ynab/examples/internal/display"
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

	plan := plans[len(plans)-1]

	transactions, _, err := client.GetTransactions(ctx, plan.ID, nil)
	if err != nil {
		log.Fatalf("failed to get transactions: %v", err)
	}

	if len(transactions) == 0 {
		fmt.Println("no transactions found")
		return
	}

	fmt.Printf("Plan: %s  (%d transactions)\n\n", plan.Name, len(transactions))
	display.PrintTransactionTable(transactions)
}
