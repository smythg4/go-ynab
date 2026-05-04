package ynab

import "math"

type DateFormat struct {
	Format string `json:"format"`
}

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

func MilliunitsToAmount(m int64) float64 {
	return float64(m) / 1000
}

func AmountToMilliunits(a float64) int64 {
	return int64(math.Round(a * 1000))
}
