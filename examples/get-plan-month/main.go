// get-plan-month retrieves budget totals for a specific month from the first
// plan on the account.
//
// Usage:
//
//	YNAB_TOKEN=your_token go run ./examples/get-plan-month
package main

import (
	"context"
	"fmt"
	"log"
	"os"

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

	month := plans[0].LastMonth

	mdata, err := client.GetMonth(ctx, plans[0].ID, month)
	if err != nil {
		log.Fatalf("failed to get month: %v", err)
	}

	fmt.Printf("Plan: %s\n", plans[0].Name)
	fmt.Printf("%-12s  %13s  %13s  %13s\n", "Month", "Income", "Budgeted", "Activity")
	fmt.Printf("%-12s  %13s  %13s  %13s\n", "------------", "-------------", "-------------", "-------------")
	fmt.Printf("%-12s  $%12.2f  $%12.2f  $%12.2f\n",
		mdata.Month,
		ynab.MilliunitsToAmount(mdata.Income),
		ynab.MilliunitsToAmount(mdata.Budgeted),
		ynab.MilliunitsToAmount(mdata.Activity),
	)
}
