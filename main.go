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
	client := ynab.NewClient(os.Getenv("YNAB_TOKEN")).WithTimeout(2 * time.Second)

	plans, err := client.GetPlans(context.Background(), true)
	if err != nil {
		log.Fatal(err)
	}

	for _, plan := range plans {
		fmt.Println(plan.Name)
		for _, acct := range plan.Accounts {
			fmt.Printf("   %s\n", acct.Name)
		}
	}
}
