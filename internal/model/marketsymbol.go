package model

type MarketSymbol struct {
	ID     int    `json:"id"`
	Symbol string `json:"symbol"`
	Info   string `json:"info"`
}

type MarketSymbolResponse struct {
	Error  int            `json:"error"`
	Result []MarketSymbol `json:"result"`
}
