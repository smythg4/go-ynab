package display

import (
	"fmt"

	"github.com/smythg4/go-ynab/ynab"
)

func truncate(s string, max int) string {
	if len(s) > max {
		return s[:max-3] + "..."
	}
	return s
}

// PrintTransactionTable prints a formatted table of transactions.
func PrintTransactionTable(transactions []ynab.Transaction) {
	fmt.Printf("%-12s  %-20s  %-25s  %13s\n", "Date", "Account", "Payee", "Amount")
	fmt.Printf("%-12s  %-20s  %-25s  %13s\n", "------------", "--------------------",
		"-------------------------", "-------------")

	for _, tx := range transactions {
		payee := ""
		if tx.PayeeName != nil {
			payee = truncate(*tx.PayeeName, 25)
		}
		fmt.Printf("%-12s  %-20s  %-25s  %13s\n",
			tx.Date,
			truncate(tx.AccountName, 20),
			payee,
			fmt.Sprintf("$%.2f", ynab.MilliunitsToAmount(tx.Amount)),
		)
	}
}
