package model

/*
{
    "id": 1,
    "last": 216415.00,
    "lowestAsk": 216678.00,
    "highestBid": 215000.00,
    "percentChange": 1.91,
    "baseVolume": 71.02603946,
    "quoteVolume": 15302897.99,
    "isFrozen": 0,
    "high24hr": 221396.00,
    "low24hr": 206414.00
}
*/
type MarketTicker struct {
	ID            int64   `json:"id"`
	Last          float64 `json:"last"`
	LowestAsk     float64 `json:"lowestAsk"`
	HighestBid    float64 `json:"highestBid"`
	PercentChange float64 `json:"percentChange"`
	BaseVolume    float64 `json:"baseVolume"`
	QuoteVolume   float64 `json:"quoteVolume"`
	IsFrozen      int64   `json:"isFrozen"`
	High24hr      float64 `json:"high24hr"`
	Low24hr       float64 `json:"low24hr"`
}
