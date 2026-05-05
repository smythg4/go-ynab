package ynab

import "math"

// DateFormat holds the date display format preference for a plan.
type DateFormat struct {
	Format string `json:"format"`
}

// CurrencyFormat holds the currency display preferences for a plan.
type CurrencyFormat struct {
	IsoCode          string `json:"iso_code"`
	ExampleFormat    string `json:"example_format"`
	DecimalDigits    int    `json:"decimal_digits"`
	DecimalSeparator string `json:"decimal_separator"`
	SymbolFirst      bool   `json:"symbol_first"`
	GroupSeparator   string `json:"group_separator"`
	CurrencySymbol   string `json:"currency_symbol"`
	DisplaySymbol    bool   `json:"display_symbol"`
}

// MilliunitsToAmount converts a YNAB milliunit value to a decimal amount (e.g. 10000 → 10.00).
func MilliunitsToAmount(m int64) float64 {
	return float64(m) / 1000
}

// AmountToMilliunits converts a decimal amount to YNAB milliunits (e.g. 10.00 → 10000).
func AmountToMilliunits(a float64) int64 {
	return int64(math.Round(a * 1000))
}
