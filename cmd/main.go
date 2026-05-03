package main

import (
	"context"
	"fmt"
	"go-ynab/ynab"
	"log"
	"os"
)

func main() {
	key := os.Getenv("YNAB_API_KEY")

	ctx := context.Background()

	client := ynab.NewClient(key).WithRateLimit(200, 10)

	plans, err := client.GetPlans(ctx)

	if err != nil {
		log.Panicln(err)
	}

	pid := plans[0].ID

	plan, err := client.GetPlan(ctx, pid)

	if err != nil {
		log.Panicln(err)
	}

	settings, err := client.GetPlanSettings(ctx, pid)

	if err != nil {
		log.Panicln(err)
	}

	user, err := client.GetUser(ctx)

	if err != nil {
		log.Panicln(err)
	}

	accounts, err := client.GetAccounts(ctx, pid, nil)

	if err != nil {
		log.Panicln(err)
	}

	aid := accounts[0].ID

	account, err := client.GetAccount(ctx, pid, aid)

	if err != nil {
		log.Panicln(err)
	}

	payees, err := client.GetPayees(ctx, pid)
	if err != nil {
		log.Panicln(err)
	}

	payid := payees[0].ID

	payee, err := client.GetPayee(ctx, pid, payid)

	if err != nil {
		log.Panicln(err)
	}

	fmt.Println(plans)
	fmt.Println(plan)
	fmt.Println(settings)
	fmt.Println(user)
	fmt.Println(accounts)
	fmt.Println(account)
	fmt.Println(payees)
	fmt.Println(payee)
}
