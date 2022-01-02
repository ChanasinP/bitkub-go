package bitkub_test

import (
	"testing"

	"github.com/ChanasinP/bitkub-go"
)

const (
	API_KEY    = ""
	API_SECRET = ""
)

func TestGetStatus(t *testing.T) {
	api := bitkub.NewBitkub("", "")
	status, err := api.GetServerStatus()
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range status {
		t.Logf("Status name: %+s, status: %s, message: %s", s.Name, s.Status, s.Message)
	}
}

func TestGetServerTime(t *testing.T) {
	api := bitkub.NewBitkub("", "")
	serverTime, err := api.GetServerTime()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Server Time: %s", serverTime.Format("2006-01-02 15:04:05"))
}

func TestGetMarketSymbols(t *testing.T) {
	api := bitkub.NewBitkub("", "")
	symbols, err := api.GetMarketSymbols()
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range symbols {
		t.Logf("Symbol: id=%d, symbol=%s, info=%s", s.ID, s.Symbol, s.Info)
	}
}

func TestGetMarketTickers(t *testing.T) {
	api := bitkub.NewBitkub("", "")
	tickers, err := api.GetMarketTickers("")
	if err != nil {
		t.Fatal(err)
	}
	for sym, ticker := range tickers {
		t.Logf("Ticker of %s sym : %+v", sym, ticker)
	}
}

func TestGetMarketTrades(t *testing.T) {
	api := bitkub.NewBitkub("", "")
	trades, err := api.GetMarketTrades("THB_BTC", 10)
	if err != nil {
		t.Fatal(err)
	}
	for _, trade := range trades {
		t.Logf("Trade timestamp %d, Rate %.03f, Amount %f, Side %s", trade.Timestamp, trade.Rate, trade.Amount, trade.Side)
	}
}

func TestGetMarketBids(t *testing.T) {
	api := bitkub.NewBitkub("", "")
	bids, err := api.GetMarketBids("THB_BTC", 10)
	if err != nil {
		t.Fatal(err)
	}
	for _, bid := range bids {
		t.Logf("Bid Order ID %d, timestamp %d, Valumn %f, Rate %.03f, Amount %f", bid.OrderID, bid.Timestamp, bid.Volumn, bid.Rate, bid.Amount)
	}
}

func TestGetMarketAsks(t *testing.T) {
	api := bitkub.NewBitkub("", "")
	asks, err := api.GetMarketAsks("THB_BTC", 10)
	if err != nil {
		t.Fatal(err)
	}
	for _, ask := range asks {
		t.Logf("Ask Order ID %d, timestamp %d, Valumn %f, Rate %.03f, Amount %f", ask.OrderID, ask.Timestamp, ask.Volumn, ask.Rate, ask.Amount)
	}
}

func TestGetMarketBooks(t *testing.T) {
	api := bitkub.NewBitkub("", "")
	books, err := api.GetMarketBooks("THB_BTC", 10)
	if err != nil {
		t.Fatal(err)
	}
	for key, book := range books {
		for _, item := range book {
			t.Logf("%s Order ID %d, timestamp %d, Valumn %f, Rate %.03f, Amount %f", key, item.OrderID, item.Timestamp, item.Volumn, item.Rate, item.Amount)
		}
	}
}

func TestGetMarketDepth(t *testing.T) {
	api := bitkub.NewBitkub("", "")
	depth, err := api.GetMarketDepth("THB_BTC", 10)
	if err != nil {
		t.Fatal(err)
	}
	for key, items := range depth {
		for _, item := range items {
			t.Logf("Depth %s, Price %.0f, Valumn %f", key, item.Price, item.Volumn)
		}
	}
}

func TestGetWallet(t *testing.T) {
	api := bitkub.NewBitkub(API_KEY, API_SECRET)
	wallet, err := api.GetWallet()
	if err != nil {
		t.Fatal(err)
	}

	for sym, volumn := range wallet {
		t.Logf("%s: %f", sym, volumn)
	}
}

func TestGetBalances(t *testing.T) {
	api := bitkub.NewBitkub(API_KEY, API_SECRET)
	wallet, err := api.GetBalances()
	if err != nil {
		t.Fatal(err)
	}

	for sym, balance := range wallet {
		t.Logf("%s: Available=%f, Reserved=%f", sym, balance.Available, balance.Reserved)
	}
}

func TestPlaceBid(t *testing.T) {
	api := bitkub.NewBitkub(API_KEY, API_SECRET)
	bid, err := api.PlaceBidTest("THB_BTC", bitkub.OrderTypeMarket, 10.01, 0)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Response : %+v", bid)
}

func TestPlaceAsk(t *testing.T) {
	api := bitkub.NewBitkub(API_KEY, API_SECRET)
	bid, err := api.PlaceAskTest("THB_BTC", bitkub.OrderTypeMarket, 0.001, 0)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Response : %+v", bid)
}

func TestGetOrderHistory(t *testing.T) {
	api := bitkub.NewBitkub(API_KEY, API_SECRET)
	orders, pagination, err := api.GetOrderHistory("THB_BTC", 1, 10, 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	for _, order := range orders {
		t.Logf("Order %+v", order)
	}
	t.Logf("Paging %+v", pagination)
}

func TestGetOrderInfo(t *testing.T) {
	api := bitkub.NewBitkub(API_KEY, API_SECRET)
	order, err := api.GetOrderInfo("THB_BTC", "buy", "fwQ6dnQWQq71S9vZ9PNzX59MF28", 1)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Response : %+v", order)
}

func TestGetCryptoAddresses(t *testing.T) {
	api := bitkub.NewBitkub(API_KEY, API_SECRET)
	addresses, pagination, err := api.GetCryptoAddresses(0, 0)
	if err != nil {
		t.Fatal(err)
	}

	for _, address := range addresses {
		t.Logf("Address %+v", address)
	}
	t.Logf("Paging %+v", pagination)
}

func TestGetCryptoDepositHistory(t *testing.T) {
	api := bitkub.NewBitkub(API_KEY, API_SECRET)
	deposits, pagination, err := api.GetCryptoDepositHistory(0, 0)
	if err != nil {
		t.Fatal(err)
	}

	for _, deposit := range deposits {
		t.Logf("Deposit %+v", deposit)
	}
	t.Logf("Paging %+v", pagination)
}

func TestGetCryptoWithdrawHistory(t *testing.T) {
	api := bitkub.NewBitkub(API_KEY, API_SECRET)
	withdraws, pagination, err := api.GetCryptoWithdrawHistory(0, 0)
	if err != nil {
		t.Fatal(err)
	}

	for _, withdraw := range withdraws {
		t.Logf("Withdraw %+v", withdraw)
	}
	t.Logf("Paging %+v", pagination)
}

func TestGetBankAccounts(t *testing.T) {
	api := bitkub.NewBitkub(API_KEY, API_SECRET)
	accounts, pagination, err := api.GetBankAccounts(0, 0)
	if err != nil {
		t.Fatal(err)
	}

	for _, account := range accounts {
		t.Logf("Account %+v", account)
	}
	t.Logf("Paging %+v", pagination)
}

func TestGetFiatDepositHistory(t *testing.T) {
	api := bitkub.NewBitkub(API_KEY, API_SECRET)
	deposits, pagination, err := api.GetFiatDepositHistory(0, 0)
	if err != nil {
		t.Fatal(err)
	}

	for _, deposit := range deposits {
		t.Logf("Deposit %+v", deposit)
	}
	t.Logf("Paging %+v", pagination)
}

func TestGetFiatWithdrawHistory(t *testing.T) {
	api := bitkub.NewBitkub(API_KEY, API_SECRET)
	withdraws, pagination, err := api.GetFiatWithdrawHistory(0, 0)
	if err != nil {
		t.Fatal(err)
	}

	for _, withdraw := range withdraws {
		t.Logf("Withdraw %+v", withdraw)
	}
	t.Logf("Paging %+v", pagination)
}

func TestGetWebSocketToken(t *testing.T) {
	api := bitkub.NewBitkub(API_KEY, API_SECRET)
	token, err := api.GetWebSocketToken()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Response : %s", token)
}

func TestGetUserLimits(t *testing.T) {
	api := bitkub.NewBitkub(API_KEY, API_SECRET)
	limits, err := api.GetUserLimits()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Response : %+v", limits)
}

func TestGetUserTradingCredits(t *testing.T) {
	api := bitkub.NewBitkub(API_KEY, API_SECRET)
	credits, err := api.GetUserTradingCredits()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Response : %f", credits)
}
