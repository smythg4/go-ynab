// get-category-balance retrieves the budgeted amount, activity, and balance
// for the first category in the first plan for the plan's most recent month.
//
// Usage:
//
//	YNAB_TOKEN=your_token go run ./examples/get-category-balance
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/smythg4/go-ynab/ynab"

	"github.com/google/uuid"
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

	idx := len(plans) - 1
	planID := plans[idx].ID
	month := plans[idx].LastMonth

	catGroups, _, err := client.GetCategories(ctx, planID)
	if err != nil {
		log.Fatalf("failed to get categories: %v", err)
	}

	var categoryID uuid.UUID
	var categoryName string
	for _, g := range catGroups {
		for _, c := range g.Categories {
			if !c.Deleted && !c.Hidden {
				categoryID = c.ID
				categoryName = c.Name
				break
			}
		}
		if categoryID != (uuid.UUID{}) {
			break
		}
	}

	if categoryID == (uuid.UUID{}) {
		fmt.Println("no categories found")
		return
	}

	cat, err := client.GetCategoryForMonth(ctx, planID, month, categoryID)
	if err != nil {
		log.Fatalf("failed to get category: %v", err)
	}

	fmt.Printf("Plan:     %s\n", plans[idx].Name)
	fmt.Printf("Month:    %s\n", month)
	fmt.Printf("Category: %s\n\n", categoryName)
	fmt.Printf("%-12s  $%12.2f\n", "Budgeted", ynab.MilliunitsToAmount(cat.Budgeted))
	fmt.Printf("%-12s  $%12.2f\n", "Activity", ynab.MilliunitsToAmount(cat.Activity))
	fmt.Printf("%-12s  $%12.2f\n", "Balance", ynab.MilliunitsToAmount(cat.Balance))
}
