package models

import "time"

type Token struct {
	ID       string  `json:"id"`
	Symbol   string  `json:"symbol"`
	Name     string  `json:"name"`
	PriceUSD float64 `json:"current_price"`
}

type TokenPriceHistory struct {
	Symbol       string    `json:"symbol"`
	CurrentPrice float64   `json:"price"`
	CreatedAt    time.Time `json:"created_at"`
}
