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
	fmt.Printf("%-12s  %-20s  %-25s  %13s\n", "Date", "Account", "Payee", "Amount")
	fmt.Printf("%-12s  %-20s  %-25s  %13s\n", "------------", "--------------------",
		"-------------------------", "-------------")

	for _, tx := range transactions {
		payee := ""
		if tx.PayeeName != nil {
			payee = *tx.PayeeName
			if len(payee) > 25 {
				payee = payee[:22] + "..."
			}
		}
		account := tx.AccountName
		if len(account) > 20 {
			account = account[:17] + "..."
		}
		fmt.Printf("%-12s  %-20s  %-25s  %13s\n",
			tx.Date,
			account,
			payee,
			fmt.Sprintf("$%.2f", ynab.MilliunitsToAmount(tx.Amount)),
		)
	}
}
