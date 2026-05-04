// list-plans retrieves all YNAB plans for the authenticated user and prints
// each plan's ID, name, and currency. The plan ID is required by most other
// API endpoints, so this is a good starting point for any integration.
//
// Usage:
//
//	YNAB_TOKEN=your_token go run ./examples/list-plans
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

	plans, err := client.GetPlans(context.Background())
	if err != nil {
		log.Fatalf("failed to get plans: %v", err)
	}

	if len(plans) == 0 {
		fmt.Println("no plans found")
		return
	}

	fmt.Printf("%-36s  %-30s  %s\n", "ID", "Name", "Currency")
	fmt.Printf("%-36s  %-30s  %s\n", "------------------------------------", "------------------------------", "--------")

	for _, plan := range plans {
		name := plan.Name
		if len(name) > 30 {
			name = name[:27] + "..."
		}
		fmt.Printf("%-36s  %-30s  %s\n", plan.ID, name, plan.CurrencyFormat.IsoCode)
	}
}
