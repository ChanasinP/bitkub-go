package model

type MarketBidAndAsk struct {
	OrderID   int64   `json:"order_id"`
	Timestamp int     `json:"timestamp"`
	Volumn    float64 `json:"volumn"`
	Rate      float64 `json:"rate"`
	Amount    float64 `json:"amount"`
}
