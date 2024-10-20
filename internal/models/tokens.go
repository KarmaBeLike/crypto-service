package models

type Token struct {
	ID       string  `json:"id"`
	Symbol   string  `json:"symbol"`
	Name     string  `json:"name"`
	PriceUSD float64 `json:"current_price"`
}
