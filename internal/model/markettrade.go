package model

type MarketTrade struct {
	Timestamp int     `json:"timestamp"`
	Rate      float64 `json:"rate"`
	Amount    float64 `json:"amount"`
	Side      string  `json:"side"`
}
