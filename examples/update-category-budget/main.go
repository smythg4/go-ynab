// update-category-budget fetches the first eligible category from the first
// plan for the authenticated user and increases its budgeted amount for the
// current month by $10, then prints the before and after state.
//
// Usage:
//
//	YNAB_TOKEN=your_token go run ./examples/update-category-budget
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/smythg4/go-ynab/ynab"
)

func firstEligibleCategory(groups []ynab.CategoryGroup) *ynab.Category {
	for _, group := range groups {
		// Skip groups that don't allow budget updates
		if group.Hidden || group.Deleted || group.Name == "Internal Master Category" {
			continue
		}
		for i := range group.Categories {
			c := &group.Categories[i]
			if !c.Hidden && !c.Deleted {
				return c
			}
		}
	}
	return nil
}

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

	catGroups, _, err := client.GetCategories(ctx, planId, nil)
	if err != nil {
		log.Fatalf("failed to get categories: %v", err)
	}

	cat := firstEligibleCategory(catGroups)
	if cat == nil {
		fmt.Println("no eligible categories found")
		return
	}

	now := time.Now()
	month := ynab.NewDate(now.Year(), now.Month(), 1)

	current, err := client.GetCategoryForMonth(ctx, planId, month, cat.ID)
	if err != nil {
		log.Fatalf("failed to get category for month: %v", err)
	}

	updated, err := client.UpdateCategoryForMonth(ctx, planId, month, cat.ID, ynab.SaveMonthCategory{
		Budgeted: current.Budgeted + 10000,
	})
	if err != nil {
		log.Fatalf("failed to update category budget: %v", err)
	}

	fmt.Printf("Updated Category Budget\n\n")
	fmt.Printf("Plan:  %s\n", plans[0].Name)
	fmt.Printf("Month: %s\n\n", month)
	fmt.Printf("   %-12s %s\n", "Category:", updated.Name)
	fmt.Printf("   %-12s $%.2f  ->  $%.2f\n", "Budgeted:",
		ynab.MilliunitsToAmount(current.Budgeted),
		ynab.MilliunitsToAmount(updated.Budgeted),
	)
}
