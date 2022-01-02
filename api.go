package bitkub

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ChanasinP/bitkub-go/internal"
	"github.com/ChanasinP/bitkub-go/internal/model"
)

const (
	baseURL         = "https://api.bitkub.com"
	OrderTypeLimit  = "limit"
	OrderTypeMarket = "market"
	OrderSideBuy    = "buy"
	OrderSideSell   = "sell"
)

type bitkubApi struct {
	Timeout   time.Duration
	ApiKey    string
	ApiSecret string
}

func NewBitkub(key, secret string, timeout ...time.Duration) *bitkubApi {
	if len(timeout) > 0 {
		return &bitkubApi{ApiKey: key, ApiSecret: secret, Timeout: timeout[0]}
	}
	return &bitkubApi{ApiKey: key, ApiSecret: secret, Timeout: time.Second * 10}
}

func (b *bitkubApi) sigPayload(payload map[string]interface{}) ([]byte, error) {
	payload["ts"] = time.Now().Unix() * 1000
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	h := hmac.New(sha256.New, []byte(b.ApiSecret))
	if _, err := h.Write(body); err != nil {
		return nil, err
	}
	payload["sig"] = hex.EncodeToString(h.Sum(nil))
	body, err = json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func formatFloatWithoutZeroTrail(num float64) string {
	var (
		zero, dot = "0", "."
		str       = fmt.Sprintf("%f", num)
	)

	return strings.TrimRight(strings.TrimRight(str, zero), dot)
}

// GetServerStatus Get endpoint status. When status is not ok, it is highly recommended to wait until the status changes back to ok.
func (b *bitkubApi) GetServerStatus() ([]model.ServerStatus, error) {
	url := baseURL + "/api/status"
	resp, err := internal.Get(url, map[string]string{}, b.Timeout)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("response status code is %d", resp.StatusCode())
	}

	ret := []model.ServerStatus{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return nil, err
	}

	return ret, nil
}

// GetServerTime Get server timestamp.
func (b *bitkubApi) GetServerTime() (time.Time, error) {
	url := baseURL + "/api/servertime"
	resp, err := internal.Get(url, map[string]string{}, b.Timeout)
	if err != nil {
		return time.Time{}, err
	}

	if resp.StatusCode() != 200 {
		return time.Time{}, fmt.Errorf("response status code is %d", resp.StatusCode())
	}

	i, err := strconv.Atoi(string(resp.Body()))
	if err != nil {
		return time.Time{}, err
	}
	fmt.Printf("%s=%d\n", resp.Body(), i)

	return time.Time{}, nil
}

// GetMarketSymbols List all available symbols.
func (b *bitkubApi) GetMarketSymbols() ([]model.MarketSymbol, error) {
	url := baseURL + "/api/market/symbols"
	resp, err := internal.Get(url, map[string]string{}, b.Timeout)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("response status code is %d", resp.StatusCode())
	}

	ret := model.MarketSymbolResponse{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return nil, err
	}

	if ret.Error != 0 {
		return nil, fmt.Errorf("got server error (%d) : %s", ret.Error, internal.GetErrorMessage(ret.Error))
	}

	return ret.Result, nil
}

// GetMarketTickers Get ticker information.
func (b *bitkubApi) GetMarketTickers(symbol string) (map[string]model.MarketTicker, error) {
	url := baseURL + "/api/market/ticker"
	if symbol != "" {
		url += "?sym=" + symbol
	}
	resp, err := internal.Get(url, map[string]string{}, b.Timeout)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("response status code is %d", resp.StatusCode())
	}

	ret := map[string]model.MarketTicker{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return nil, err
	}

	return ret, nil
}

// GetMarketTrades List recent trades.
func (b *bitkubApi) GetMarketTrades(symbol string, limit int) ([]model.MarketTrade, error) {
	url := baseURL + "/api/market/trades"
	params := []string{}
	if symbol != "" {
		params = append(params, "sym="+symbol)
	}
	if limit > 0 {
		params = append(params, fmt.Sprintf("lmt=%d", limit))
	}
	url += "?" + strings.Join(params, "&")
	resp, err := internal.Get(url, map[string]string{}, b.Timeout)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("response status code is %d", resp.StatusCode())
	}
	ret := model.DefaultResponse{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return nil, err
	}
	if ret.Error != 0 {
		return nil, fmt.Errorf("got server error (%d) : %s", ret.Error, internal.GetErrorMessage(ret.Error))
	}

	trades := []model.MarketTrade{}
	for _, item := range ret.Result.([]interface{}) {
		t := item.([]interface{})
		trades = append(trades, model.MarketTrade{
			Timestamp: int(t[0].(float64)),
			Rate:      t[1].(float64),
			Amount:    t[2].(float64),
			Side:      t[3].(string),
		})
	}

	return trades, nil
}

// GetMarketBids List open buy orders.
func (b *bitkubApi) GetMarketBids(symbol string, limit int) ([]model.MarketBidAndAsk, error) {
	url := baseURL + "/api/market/bids"
	params := []string{}
	if symbol != "" {
		params = append(params, "sym="+symbol)
	}
	if limit > 0 {
		params = append(params, fmt.Sprintf("lmt=%d", limit))
	}
	url += "?" + strings.Join(params, "&")
	resp, err := internal.Get(url, map[string]string{}, b.Timeout)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("response status code is %d", resp.StatusCode())
	}
	ret := model.DefaultResponse{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return nil, err
	}
	if ret.Error != 0 {
		return nil, fmt.Errorf("got server error (%d) : %s", ret.Error, internal.GetErrorMessage(ret.Error))
	}

	bids := []model.MarketBidAndAsk{}
	for _, item := range ret.Result.([]interface{}) {
		b := item.([]interface{})
		bids = append(bids, model.MarketBidAndAsk{
			OrderID:   int64(b[0].(float64)),
			Timestamp: int(b[1].(float64)),
			Volumn:    b[2].(float64),
			Rate:      b[3].(float64),
			Amount:    b[4].(float64),
		})
	}

	return bids, nil
}

// GetMarketAsks List open sell orders.
func (b *bitkubApi) GetMarketAsks(symbol string, limit int) ([]model.MarketBidAndAsk, error) {
	url := baseURL + "/api/market/asks"
	params := []string{}
	if symbol != "" {
		params = append(params, "sym="+symbol)
	}
	if limit > 0 {
		params = append(params, fmt.Sprintf("lmt=%d", limit))
	}
	url += "?" + strings.Join(params, "&")
	resp, err := internal.Get(url, map[string]string{}, b.Timeout)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("response status code is %d", resp.StatusCode())
	}
	ret := model.DefaultResponse{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return nil, err
	}
	if ret.Error != 0 {
		return nil, fmt.Errorf("got server error (%d) : %s", ret.Error, internal.GetErrorMessage(ret.Error))
	}

	asks := []model.MarketBidAndAsk{}
	for _, item := range ret.Result.([]interface{}) {
		a := item.([]interface{})
		asks = append(asks, model.MarketBidAndAsk{
			OrderID:   int64(a[0].(float64)),
			Timestamp: int(a[1].(float64)),
			Volumn:    a[2].(float64),
			Rate:      a[3].(float64),
			Amount:    a[4].(float64),
		})
	}

	return asks, nil
}

// GetMarketOrderbook List all open orders.
func (b *bitkubApi) GetMarketBooks(symbol string, limit int) (map[string][]model.MarketBidAndAsk, error) {
	url := baseURL + "/api/market/books"
	params := []string{}
	if symbol != "" {
		params = append(params, "sym="+symbol)
	}
	if limit > 0 {
		params = append(params, fmt.Sprintf("lmt=%d", limit))
	}
	url += "?" + strings.Join(params, "&")
	resp, err := internal.Get(url, map[string]string{}, b.Timeout)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("response status code is %d", resp.StatusCode())
	}
	ret := model.DefaultResponse{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return nil, err
	}
	if ret.Error != 0 {
		return nil, fmt.Errorf("got server error (%d) : %s", ret.Error, internal.GetErrorMessage(ret.Error))
	}

	books := map[string][]model.MarketBidAndAsk{}
	for key, item := range ret.Result.(map[string]interface{}) {
		bidOrAsk := item.([]interface{})
		books[key] = []model.MarketBidAndAsk{}
		for _, b := range bidOrAsk {
			b := b.([]interface{})
			books[key] = append(books[key], model.MarketBidAndAsk{
				OrderID:   int64(b[0].(float64)),
				Timestamp: int(b[1].(float64)),
				Volumn:    b[2].(float64),
				Rate:      b[3].(float64),
				Amount:    b[4].(float64),
			})
		}
	}

	return books, nil
}

// GetTradingViewHistory Get historical data for TradingView chart.
func (b *bitkubApi) GetTradingViewHistory(symbol, resolution string, from, to int) (map[string]interface{}, error) {
	url := baseURL + "/tradingview/history"
	params := []string{}
	if symbol != "" {
		params = append(params, "sym="+symbol)
	}
	if resolution != "" {
		params = append(params, "resolution="+resolution)
	}
	if from > 0 {
		params = append(params, fmt.Sprintf("from=%d", from))
	}
	if to > 0 {
		params = append(params, fmt.Sprintf("to=%d", to))
	}
	url += "?" + strings.Join(params, "&")
	resp, err := internal.Get(url, map[string]string{}, b.Timeout)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("response status code is %d", resp.StatusCode())
	}
	ret := map[string]interface{}{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// GetMarketDepth Get depth information.
func (b *bitkubApi) GetMarketDepth(symbol string, limit int) (map[string][]model.MarketDepth, error) {
	url := baseURL + "/api/market/depth"
	params := []string{}
	if symbol != "" {
		params = append(params, "sym="+symbol)
	}
	if limit > 0 {
		params = append(params, fmt.Sprintf("lmt=%d", limit))
	}
	url += "?" + strings.Join(params, "&")
	resp, err := internal.Get(url, map[string]string{}, b.Timeout)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("response status code is %d", resp.StatusCode())
	}
	ret := map[string]interface{}{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return nil, err
	}

	depth := map[string][]model.MarketDepth{}
	for key, itemList := range ret {
		depth[key] = []model.MarketDepth{}
		for _, item := range itemList.([]interface{}) {
			item := item.([]interface{})
			depth[key] = append(depth[key], model.MarketDepth{
				Price:  item[0].(float64),
				Volumn: item[1].(float64),
			})
		}
	}
	return depth, nil
}

// GetWallet Get user available balances (for both available and reserved balances please use GetBalances)
func (b *bitkubApi) GetWallet() (map[string]interface{}, error) {
	if b.ApiKey == "" {
		return nil, fmt.Errorf("api key is empty")
	}
	if b.ApiSecret == "" {
		return nil, fmt.Errorf("api secret is empty")
	}

	url := baseURL + "/api/market/wallet"
	payload, err := b.sigPayload(map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	resp, err := internal.Post(url, map[string]string{"X-BTK-APIKEY": b.ApiKey}, payload, b.Timeout)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("got server error (%d) : %s", resp.StatusCode(), resp.String())
	}

	ret := model.DefaultResponse{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return nil, err
	}
	if ret.Error != 0 {
		return nil, fmt.Errorf("got server error (%d) : %s", ret.Error, internal.GetErrorMessage(ret.Error))
	}

	return ret.Result.(map[string]interface{}), nil
}

// GetBalances Get balances info: this includes both available and reserved balances.
func (b *bitkubApi) GetBalances() (map[string]model.Balance, error) {
	if b.ApiKey == "" {
		return nil, fmt.Errorf("api key is empty")
	}
	if b.ApiSecret == "" {
		return nil, fmt.Errorf("api secret is empty")
	}

	url := baseURL + "/api/market/balances"
	payload, err := b.sigPayload(map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	resp, err := internal.Post(url, map[string]string{"X-BTK-APIKEY": b.ApiKey}, payload, b.Timeout)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("got server error (%d) : %s", resp.StatusCode(), resp.String())
	}

	ret := model.DefaultResponse{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return nil, err
	}
	if ret.Error != 0 {
		return nil, fmt.Errorf("got server error (%d) : %s", ret.Error, internal.GetErrorMessage(ret.Error))
	}

	balances := map[string]model.Balance{}
	for key, item := range ret.Result.(map[string]interface{}) {
		balances[key] = model.Balance{
			Available: item.(map[string]interface{})["available"].(float64),
			Reserved:  item.(map[string]interface{})["reserved"].(float64),
		}
	}
	return balances, nil
}

// PlaceBid Create a buy order.
func (b *bitkubApi) PlaceBid(symbol, bitType string, amount, rate float64, clientID ...string) (*model.Order, error) {
	if b.ApiKey == "" {
		return nil, fmt.Errorf("api key is empty")
	}
	if b.ApiSecret == "" {
		return nil, fmt.Errorf("api secret is empty")
	}
	if symbol == "" {
		return nil, fmt.Errorf("symbol is empty")
	}
	if bitType == "" {
		return nil, fmt.Errorf("bit type is empty")
	}
	if bitType != OrderTypeLimit && bitType != OrderTypeMarket {
		return nil, fmt.Errorf("bit type is invalid")
	}
	if bitType == OrderTypeMarket {
		rate = 0
	}

	url := baseURL + "/api/market/place-bid"

	payload := map[string]interface{}{
		"sym": symbol,
		"typ": bitType,
		"amt": formatFloatWithoutZeroTrail(amount),
		"rat": formatFloatWithoutZeroTrail(rate),
	}
	if len(clientID) > 0 {
		payload["client_id"] = clientID[0]
	}
	payloadBytes, err := b.sigPayload(payload)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{"X-BTK-APIKEY": b.ApiKey, "Content-Type": "application/json", "Accept": "application/json"}
	resp, err := internal.Post(url, headers, payloadBytes, b.Timeout)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("got server error (%d) : %s", resp.StatusCode(), resp.String())
	}

	ret := model.OrderResponse{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return nil, err
	}

	if ret.Error != 0 {
		return nil, fmt.Errorf("got server error (%d) : %s", ret.Error, internal.GetErrorMessage(ret.Error))
	}

	return &ret.Result, nil
}

// PlaceBidTest Test creating a buy order (no balance is deducted).
func (b *bitkubApi) PlaceBidTest(symbol, bitType string, amount, rate float64, clientID ...string) (*model.Order, error) {
	if b.ApiKey == "" {
		return nil, fmt.Errorf("api key is empty")
	}
	if b.ApiSecret == "" {
		return nil, fmt.Errorf("api secret is empty")
	}
	if symbol == "" {
		return nil, fmt.Errorf("symbol is empty")
	}
	if bitType == "" {
		return nil, fmt.Errorf("bit type is empty")
	}
	if bitType != OrderTypeLimit && bitType != OrderTypeMarket {
		return nil, fmt.Errorf("bit type is invalid")
	}
	if bitType == OrderTypeMarket {
		rate = 0
	}

	url := baseURL + "/api/market/place-bid/test"
	payload := map[string]interface{}{
		"sym": symbol,
		"typ": bitType,
		"amt": formatFloatWithoutZeroTrail(amount),
		"rat": formatFloatWithoutZeroTrail(rate),
	}
	if len(clientID) > 0 {
		payload["client_id"] = clientID[0]
	}
	payloadBytes, err := b.sigPayload(payload)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{"X-BTK-APIKEY": b.ApiKey, "Content-Type": "application/json", "Accept": "application/json"}
	resp, err := internal.Post(url, headers, payloadBytes, b.Timeout)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("got server error (%d) : %s", resp.StatusCode(), resp.String())
	}

	ret := model.OrderResponse{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return nil, err
	}

	if ret.Error != 0 {
		return nil, fmt.Errorf("got server error (%d) : %s", ret.Error, internal.GetErrorMessage(ret.Error))
	}

	return &ret.Result, nil
}

// PlaceAsk Create a sell order.
func (b *bitkubApi) PlaceAsk(symbol, bitType string, amount, rate float64, clientID ...string) (*model.Order, error) {
	if b.ApiKey == "" {
		return nil, fmt.Errorf("api key is empty")
	}
	if b.ApiSecret == "" {
		return nil, fmt.Errorf("api secret is empty")
	}
	if symbol == "" {
		return nil, fmt.Errorf("symbol is empty")
	}
	if bitType == "" {
		return nil, fmt.Errorf("bit type is empty")
	}
	if bitType != OrderTypeLimit && bitType != OrderTypeMarket {
		return nil, fmt.Errorf("bit type is invalid")
	}
	if bitType == OrderTypeMarket {
		rate = 0
	}

	url := baseURL + "/api/market/place-ask"
	payload := map[string]interface{}{
		"sym": symbol,
		"typ": bitType,
		"amt": formatFloatWithoutZeroTrail(amount),
		"rat": formatFloatWithoutZeroTrail(rate),
	}
	if len(clientID) > 0 {
		payload["client_id"] = clientID[0]
	}
	payloadBytes, err := b.sigPayload(payload)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{"X-BTK-APIKEY": b.ApiKey, "Content-Type": "application/json", "Accept": "application/json"}
	resp, err := internal.Post(url, headers, payloadBytes, b.Timeout)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("got server error (%d) : %s", resp.StatusCode(), resp.String())
	}

	ret := model.OrderResponse{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return nil, err
	}

	if ret.Error != 0 {
		return nil, fmt.Errorf("got server error (%d) : %s", ret.Error, internal.GetErrorMessage(ret.Error))
	}

	return &ret.Result, nil
}

// PlaceAskTest Test creating a sell order (no balance is deducted).
func (b *bitkubApi) PlaceAskTest(symbol, bitType string, amount, rate float64, clientID ...string) (*model.Order, error) {
	if b.ApiKey == "" {
		return nil, fmt.Errorf("api key is empty")
	}
	if b.ApiSecret == "" {
		return nil, fmt.Errorf("api secret is empty")
	}
	if symbol == "" {
		return nil, fmt.Errorf("symbol is empty")
	}
	if bitType == "" {
		return nil, fmt.Errorf("bit type is empty")
	}
	if bitType != OrderTypeLimit && bitType != OrderTypeMarket {
		return nil, fmt.Errorf("bit type is invalid")
	}
	if bitType == OrderTypeMarket {
		rate = 0
	}

	url := baseURL + "/api/market/place-ask/test"
	payload := map[string]interface{}{
		"sym": symbol,
		"typ": bitType,
		"amt": formatFloatWithoutZeroTrail(amount),
		"rat": formatFloatWithoutZeroTrail(rate),
	}
	if len(clientID) > 0 {
		payload["client_id"] = clientID[0]
	}
	payloadBytes, err := b.sigPayload(payload)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{"X-BTK-APIKEY": b.ApiKey, "Content-Type": "application/json", "Accept": "application/json"}
	resp, err := internal.Post(url, headers, payloadBytes, b.Timeout)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("got server error (%d) : %s", resp.StatusCode(), resp.String())
	}

	ret := model.OrderResponse{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return nil, err
	}

	if ret.Error != 0 {
		return nil, fmt.Errorf("got server error (%d) : %s", ret.Error, internal.GetErrorMessage(ret.Error))
	}

	return &ret.Result, nil
}

// PlaceAskByFiat Create a sell order by specifying the fiat amount you want to receive (selling amount of cryptocurrency is automatically calculated). If order type is market, currrent highest bid will be used as rate.
func (b *bitkubApi) PlaceAskByFiat(symbol, bitType string, amount, rate float64) (*model.Order, error) {
	if b.ApiKey == "" {
		return nil, fmt.Errorf("api key is empty")
	}
	if b.ApiSecret == "" {
		return nil, fmt.Errorf("api secret is empty")
	}
	if symbol == "" {
		return nil, fmt.Errorf("symbol is empty")
	}
	if bitType == "" {
		return nil, fmt.Errorf("bit type is empty")
	}
	if bitType != OrderTypeLimit && bitType != OrderTypeMarket {
		return nil, fmt.Errorf("bit type is invalid")
	}
	if bitType == OrderTypeMarket {
		rate = 0
	}

	url := baseURL + "/api/market/place-ask-by-fiat"
	payload := map[string]interface{}{
		"sym": symbol,
		"typ": bitType,
		"amt": formatFloatWithoutZeroTrail(amount),
		"rat": formatFloatWithoutZeroTrail(rate),
	}

	payloadBytes, err := b.sigPayload(payload)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{"X-BTK-APIKEY": b.ApiKey, "Content-Type": "application/json"}
	resp, err := internal.Post(url, headers, payloadBytes, b.Timeout)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("got server error (%d) : %s", resp.StatusCode(), resp.String())
	}

	ret := model.OrderResponse{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return nil, err
	}

	if ret.Error != 0 {
		return nil, fmt.Errorf("got server error (%d) : %s", ret.Error, internal.GetErrorMessage(ret.Error))
	}

	return &ret.Result, nil
}

// CancelOrder Cancel an open order.
func (b *bitkubApi) CancelOrder(symbol, side, hash string, id int) error {
	if b.ApiKey == "" {
		return fmt.Errorf("api key is empty")
	}
	if b.ApiSecret == "" {
		return fmt.Errorf("api secret is empty")
	}
	var payload map[string]interface{}
	if hash == "" {
		if symbol == "" {
			return fmt.Errorf("symbol is empty")
		}
		if side == "" {
			return fmt.Errorf("side is empty")
		}
		if id == 0 {
			return fmt.Errorf("id is empty")
		}
		if side != OrderSideBuy && side != OrderSideSell {
			return fmt.Errorf("side is invalid")
		}
		payload = map[string]interface{}{
			"sym": symbol,
			"id":  fmt.Sprintf("%d", id),
			"sd":  side,
		}
	} else {
		payload = map[string]interface{}{"hash": hash}
	}

	url := baseURL + "/api/market/cancel-order"
	payloadBytes, err := b.sigPayload(payload)
	if err != nil {
		return err
	}

	headers := map[string]string{"X-BTK-APIKEY": b.ApiKey, "Content-Type": "application/json", "Accept": "application/json"}
	resp, err := internal.Post(url, headers, payloadBytes, b.Timeout)
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("got server error (%d) : %s", resp.StatusCode(), resp.String())
	}

	ret := model.ErrorResponse{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return err
	}

	if ret.Error != 0 {
		return fmt.Errorf("got server error (%d) : %s", ret.Error, internal.GetErrorMessage(ret.Error))
	}

	return nil
}

// GetOpenOrder List all open orders of the given symbol.
func (b *bitkubApi) GetOpenOrder(symbol string) ([]model.OpenOrder, error) {
	if b.ApiKey == "" {
		return nil, fmt.Errorf("api key is empty")
	}
	if b.ApiSecret == "" {
		return nil, fmt.Errorf("api secret is empty")
	}
	if symbol == "" {
		return nil, fmt.Errorf("symbol is empty")
	}

	url := baseURL + "/api/market/my-open-orders"
	payload := map[string]interface{}{"sym": symbol}
	payloadBytes, err := b.sigPayload(payload)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{"X-BTK-APIKEY": b.ApiKey, "Content-Type": "application/json", "Accept": "application/json"}
	resp, err := internal.Post(url, headers, payloadBytes, b.Timeout)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("got server error (%d) : %s", resp.StatusCode(), resp.String())
	}

	ret := model.OpenOrderResponse{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return nil, err
	}

	if ret.Error != 0 {
		return nil, fmt.Errorf("got server error (%d) : %s", ret.Error, internal.GetErrorMessage(ret.Error))
	}

	return ret.Result, nil
}

// GetOrderHistory List all orders that have already matched.
func (b *bitkubApi) GetOrderHistory(symbol string, page, limit int, start, end int64) ([]model.OrderHistory, *model.OrderHistoryPagination, error) {
	if b.ApiKey == "" {
		return nil, nil, fmt.Errorf("api key is empty")
	}
	if b.ApiSecret == "" {
		return nil, nil, fmt.Errorf("api secret is empty")
	}
	if symbol == "" {
		return nil, nil, fmt.Errorf("symbol is empty")
	}

	url := baseURL + "/api/market/my-order-history"
	payload := map[string]interface{}{"sym": symbol}
	if page > 0 {
		payload["p"] = page
	}
	if limit > 0 {
		payload["lmt"] = limit
	}
	if start > 0 {
		payload["start"] = start
	}
	if end > 0 {
		payload["end"] = end
	}

	payloadBytes, err := b.sigPayload(payload)
	if err != nil {
		return nil, nil, err
	}

	headers := map[string]string{"X-BTK-APIKEY": b.ApiKey, "Content-Type": "application/json", "Accept": "application/json"}
	resp, err := internal.Post(url, headers, payloadBytes, b.Timeout)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, nil, fmt.Errorf("got server error (%d) : %s", resp.StatusCode(), resp.String())
	}

	ret := model.OrderHistoryResponse{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return nil, nil, err
	}

	if ret.Error != 0 {
		return nil, nil, fmt.Errorf("got server error (%d) : %s", ret.Error, internal.GetErrorMessage(ret.Error))
	}

	return ret.Result, &ret.Pagination, nil
}

// GetOrderInfo Get information regarding the specified order.
func (b *bitkubApi) GetOrderInfo(symbol, side, hash string, id int) (*model.OrderInfo, error) {
	if b.ApiKey == "" {
		return nil, fmt.Errorf("api key is empty")
	}
	if b.ApiSecret == "" {
		return nil, fmt.Errorf("api secret is empty")
	}
	var payload map[string]interface{}
	if hash == "" {
		if symbol == "" {
			return nil, fmt.Errorf("symbol is empty")
		}
		if side == "" {
			return nil, fmt.Errorf("side is empty")
		}
		if id == 0 {
			return nil, fmt.Errorf("id is empty")
		}
		if side != OrderSideBuy && side != OrderSideSell {
			return nil, fmt.Errorf("side is invalid")
		}
		payload = map[string]interface{}{
			"sym": symbol,
			"id":  fmt.Sprintf("%d", id),
			"sd":  side,
		}
	} else {
		payload = map[string]interface{}{"hash": hash}
	}

	payloadBytes, err := b.sigPayload(payload)
	if err != nil {
		return nil, err
	}

	url := baseURL + "/api/market/order-info"
	headers := map[string]string{"X-BTK-APIKEY": b.ApiKey, "Content-Type": "application/json", "Accept": "application/json"}
	resp, err := internal.Post(url, headers, payloadBytes, b.Timeout)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("got server error (%d) : %s", resp.StatusCode(), resp.String())
	}

	ret := model.OrderInfoResponse{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return nil, err
	}

	if ret.Error != 0 {
		return nil, fmt.Errorf("got server error (%d) : %s", ret.Error, internal.GetErrorMessage(ret.Error))
	}

	return &ret.Result, nil
}

// GetCryptoAddresses List all crypto addresses.
func (b *bitkubApi) GetCryptoAddresses(page, limit int) ([]model.CryptoAddress, *model.Pagination, error) {
	if b.ApiKey == "" {
		return nil, nil, fmt.Errorf("api key is empty")
	}
	if b.ApiSecret == "" {
		return nil, nil, fmt.Errorf("api secret is empty")
	}

	payload := map[string]interface{}{}
	if page > 0 {
		payload["p"] = page
	}
	if limit > 0 {
		payload["lmt"] = limit
	}

	payloadBytes, err := b.sigPayload(payload)
	if err != nil {
		return nil, nil, err
	}

	url := baseURL + "/api/crypto/addresses"
	headers := map[string]string{"X-BTK-APIKEY": b.ApiKey, "Content-Type": "application/json", "Accept": "application/json"}
	resp, err := internal.Post(url, headers, payloadBytes, b.Timeout)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, nil, fmt.Errorf("got server error (%d) : %s", resp.StatusCode(), resp.String())
	}

	ret := model.CryptoAddressResponse{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return nil, nil, err
	}

	if ret.Error != 0 {
		return nil, nil, fmt.Errorf("got server error (%d) : %s", ret.Error, internal.GetErrorMessage(ret.Error))
	}

	return ret.Result, &ret.Pagination, nil
}

// CryptoWithdraw Make a withdrawal to a trusted address.
func (b *bitkubApi) CryptoWithdraw(currency, address string, amount float64, memo string) (*model.CryptoWithdraw, error) {
	if b.ApiKey == "" {
		return nil, fmt.Errorf("api key is empty")
	}
	if b.ApiSecret == "" {
		return nil, fmt.Errorf("api secret is empty")
	}
	if currency == "" {
		return nil, fmt.Errorf("currency is empty")
	}
	if address == "" {
		return nil, fmt.Errorf("address is empty")
	}
	if amount <= 0 {
		return nil, fmt.Errorf("amount is invalid")
	}

	payload := map[string]interface{}{
		"cur": currency,
		"adr": address,
		"amt": fmt.Sprintf("%f", amount),
	}
	if memo != "" {
		payload["mem"] = memo
	}

	payloadBytes, err := b.sigPayload(payload)
	if err != nil {
		return nil, err
	}

	url := baseURL + "/api/crypto/withdraw"
	headers := map[string]string{"X-BTK-APIKEY": b.ApiKey, "Content-Type": "application/json", "Accept": "application/json"}
	resp, err := internal.Post(url, headers, payloadBytes, b.Timeout)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("got server error (%d) : %s", resp.StatusCode(), resp.String())
	}

	ret := model.CryptoWithdrawResponse{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return nil, err
	}

	if ret.Error != 0 {
		return nil, fmt.Errorf("got server error (%d) : %s", ret.Error, internal.GetErrorMessage(ret.Error))
	}

	return &ret.Result, nil
}

// CryptoInternalWithdraw Make a withdraw to an internal address. The destination address is not required to be a trusted address. This API is not enabled by default, Only KYB users can request this feature by contacting us via support@bitkub.com
func (b *bitkubApi) CryptoInternalWithdraw(currency, address string, amount float64, memo string) (*model.CryptoWithdraw, error) {
	if b.ApiKey == "" {
		return nil, fmt.Errorf("api key is empty")
	}
	if b.ApiSecret == "" {
		return nil, fmt.Errorf("api secret is empty")
	}
	if currency == "" {
		return nil, fmt.Errorf("currency is empty")
	}
	if address == "" {
		return nil, fmt.Errorf("address is empty")
	}
	if amount <= 0 {
		return nil, fmt.Errorf("amount is invalid")
	}

	payload := map[string]interface{}{
		"cur": currency,
		"adr": address,
		"amt": fmt.Sprintf("%f", amount),
	}
	if memo != "" {
		payload["mem"] = memo
	}

	payloadBytes, err := b.sigPayload(payload)
	if err != nil {
		return nil, err
	}

	url := baseURL + "/api/crypto/internal-withdraw"
	headers := map[string]string{"X-BTK-APIKEY": b.ApiKey, "Content-Type": "application/json", "Accept": "application/json"}
	resp, err := internal.Post(url, headers, payloadBytes, b.Timeout)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("got server error (%d) : %s", resp.StatusCode(), resp.String())
	}

	ret := model.CryptoWithdrawResponse{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return nil, err
	}

	if ret.Error != 0 {
		return nil, fmt.Errorf("got server error (%d) : %s", ret.Error, internal.GetErrorMessage(ret.Error))
	}

	return &ret.Result, nil
}

// GetCryptoDepositHistory List crypto deposit history.
func (b *bitkubApi) GetCryptoDepositHistory(page, limit int) ([]model.CryptoDeposit, *model.Pagination, error) {
	if b.ApiKey == "" {
		return nil, nil, fmt.Errorf("api key is empty")
	}
	if b.ApiSecret == "" {
		return nil, nil, fmt.Errorf("api secret is empty")
	}

	payload := map[string]interface{}{}
	if page > 0 {
		payload["p"] = page
	}
	if limit > 0 {
		payload["lmt"] = limit
	}

	payloadBytes, err := b.sigPayload(payload)
	if err != nil {
		return nil, nil, err
	}

	url := baseURL + "/api/crypto/deposit-history"
	headers := map[string]string{"X-BTK-APIKEY": b.ApiKey, "Content-Type": "application/json", "Accept": "application/json"}
	resp, err := internal.Post(url, headers, payloadBytes, b.Timeout)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, nil, fmt.Errorf("got server error (%d) : %s", resp.StatusCode(), resp.String())
	}

	ret := model.CryptoDepositResponse{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return nil, nil, err
	}

	if ret.Error != 0 {
		return nil, nil, fmt.Errorf("got server error (%d) : %s", ret.Error, internal.GetErrorMessage(ret.Error))
	}

	return ret.Result, &ret.Pagination, nil
}

// GetCryptoWithdrawHistory List crypto withdrawal history.
func (b *bitkubApi) GetCryptoWithdrawHistory(page, limit int) ([]model.CryptoWithdraw, *model.Pagination, error) {
	if b.ApiKey == "" {
		return nil, nil, fmt.Errorf("api key is empty")
	}
	if b.ApiSecret == "" {
		return nil, nil, fmt.Errorf("api secret is empty")
	}

	payload := map[string]interface{}{}
	if page > 0 {
		payload["p"] = page
	}
	if limit > 0 {
		payload["lmt"] = limit
	}

	payloadBytes, err := b.sigPayload(payload)
	if err != nil {
		return nil, nil, err
	}

	url := baseURL + "/api/crypto/withdraw-history"
	headers := map[string]string{"X-BTK-APIKEY": b.ApiKey, "Content-Type": "application/json", "Accept": "application/json"}
	resp, err := internal.Post(url, headers, payloadBytes, b.Timeout)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, nil, fmt.Errorf("got server error (%d) : %s", resp.StatusCode(), resp.String())
	}

	ret := model.CryptoWithdrawHistoryResponse{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return nil, nil, err
	}

	if ret.Error != 0 {
		return nil, nil, fmt.Errorf("got server error (%d) : %s", ret.Error, internal.GetErrorMessage(ret.Error))
	}

	return ret.Result, &ret.Pagination, nil
}

// CryptoGenerateAddress Generate a new crypto address (will replace existing address; previous address can still be used to received funds)
func (b *bitkubApi) CryptoGenerateAddress(symbol string) ([]model.CryptoGenerateAddress, error) {
	if b.ApiKey == "" {
		return nil, fmt.Errorf("api key is empty")
	}
	if b.ApiSecret == "" {
		return nil, fmt.Errorf("api secret is empty")
	}
	if symbol == "" {
		return nil, fmt.Errorf("symbol is empty")
	}

	payload := map[string]interface{}{"sym": symbol}
	payloadBytes, err := b.sigPayload(payload)
	if err != nil {
		return nil, err
	}

	url := baseURL + "/api/crypto/generate-address"
	headers := map[string]string{"X-BTK-APIKEY": b.ApiKey, "Content-Type": "application/json", "Accept": "application/json"}
	resp, err := internal.Post(url, headers, payloadBytes, b.Timeout)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("got server error (%d) : %s", resp.StatusCode(), resp.String())
	}

	ret := model.CryptoGenerateAddressResponse{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return nil, err
	}

	if ret.Error != 0 {
		return nil, fmt.Errorf("got server error (%d) : %s", ret.Error, internal.GetErrorMessage(ret.Error))
	}

	return ret.Result, nil
}

// GetBankAccounts List all approved bank accounts.
func (b *bitkubApi) GetBankAccounts(page, limit int) ([]model.BankAccount, *model.Pagination, error) {
	if b.ApiKey == "" {
		return nil, nil, fmt.Errorf("api key is empty")
	}
	if b.ApiSecret == "" {
		return nil, nil, fmt.Errorf("api secret is empty")
	}

	payload := map[string]interface{}{}
	if page > 0 {
		payload["p"] = page
	}
	if limit > 0 {
		payload["lmt"] = limit
	}

	payloadBytes, err := b.sigPayload(payload)
	if err != nil {
		return nil, nil, err
	}

	url := baseURL + "/api/fiat/accounts"
	headers := map[string]string{"X-BTK-APIKEY": b.ApiKey, "Content-Type": "application/json", "Accept": "application/json"}
	resp, err := internal.Post(url, headers, payloadBytes, b.Timeout)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, nil, fmt.Errorf("got server error (%d) : %s", resp.StatusCode(), resp.String())
	}

	ret := model.FiatAccountsResponse{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return nil, nil, err
	}

	if ret.Error != 0 {
		return nil, nil, fmt.Errorf("got server error (%d) : %s", ret.Error, internal.GetErrorMessage(ret.Error))
	}

	return ret.Result, &ret.Pagination, nil
}

// FiatWithdraw Make a withdrawal to an approved bank account.
func (b *bitkubApi) FiatWithdraw(bankID string, amount float64) (*model.FiatWithdraw, error) {
	if b.ApiKey == "" {
		return nil, fmt.Errorf("api key is empty")
	}
	if b.ApiSecret == "" {
		return nil, fmt.Errorf("api secret is empty")
	}
	if bankID == "" {
		return nil, fmt.Errorf("bank id is empty")
	}
	if amount <= 0 {
		return nil, fmt.Errorf("amount is invalid")
	}

	payload := map[string]interface{}{"id": bankID, "amount": amount}
	payloadBytes, err := b.sigPayload(payload)
	if err != nil {
		return nil, err
	}

	url := baseURL + "/api/fiat/withdraw"
	headers := map[string]string{"X-BTK-APIKEY": b.ApiKey, "Content-Type": "application/json", "Accept": "application/json"}
	resp, err := internal.Post(url, headers, payloadBytes, b.Timeout)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("got server error (%d) : %s", resp.StatusCode(), resp.String())
	}

	ret := model.FiatWithdrawResponse{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return nil, err
	}

	if ret.Error != 0 {
		return nil, fmt.Errorf("got server error (%d) : %s", ret.Error, internal.GetErrorMessage(ret.Error))
	}

	return &ret.Result, nil
}

// GetFiatDepositHistory List fiat deposit history.
func (b *bitkubApi) GetFiatDepositHistory(page, limit int) ([]model.FiatDeposit, *model.Pagination, error) {
	if b.ApiKey == "" {
		return nil, nil, fmt.Errorf("api key is empty")
	}
	if b.ApiSecret == "" {
		return nil, nil, fmt.Errorf("api secret is empty")
	}

	payload := map[string]interface{}{}
	if page > 0 {
		payload["p"] = page
	}
	if limit > 0 {
		payload["lmt"] = limit
	}

	payloadBytes, err := b.sigPayload(payload)
	if err != nil {
		return nil, nil, err
	}

	url := baseURL + "/api/fiat/deposit-history"
	headers := map[string]string{"X-BTK-APIKEY": b.ApiKey, "Content-Type": "application/json", "Accept": "application/json"}
	resp, err := internal.Post(url, headers, payloadBytes, b.Timeout)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, nil, fmt.Errorf("got server error (%d) : %s", resp.StatusCode(), resp.String())
	}

	ret := model.FiatDepositResponse{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return nil, nil, err
	}

	if ret.Error != 0 {
		return nil, nil, fmt.Errorf("got server error (%d) : %s", ret.Error, internal.GetErrorMessage(ret.Error))
	}

	return ret.Result, &ret.Pagination, nil
}

// GetFiatWithdrawHistory List fiat withdrawal history.
func (b *bitkubApi) GetFiatWithdrawHistory(page, limit int) ([]model.FiatWithdraw, *model.Pagination, error) {
	if b.ApiKey == "" {
		return nil, nil, fmt.Errorf("api key is empty")
	}
	if b.ApiSecret == "" {
		return nil, nil, fmt.Errorf("api secret is empty")
	}

	payload := map[string]interface{}{}
	if page > 0 {
		payload["p"] = page
	}
	if limit > 0 {
		payload["lmt"] = limit
	}

	payloadBytes, err := b.sigPayload(payload)
	if err != nil {
		return nil, nil, err
	}

	url := baseURL + "/api/fiat/withdraw-history"
	headers := map[string]string{"X-BTK-APIKEY": b.ApiKey, "Content-Type": "application/json", "Accept": "application/json"}
	resp, err := internal.Post(url, headers, payloadBytes, b.Timeout)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, nil, fmt.Errorf("got server error (%d) : %s", resp.StatusCode(), resp.String())
	}

	ret := model.FiatWithdrawHistoryResponse{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return nil, nil, err
	}

	if ret.Error != 0 {
		return nil, nil, fmt.Errorf("got server error (%d) : %s", ret.Error, internal.GetErrorMessage(ret.Error))
	}

	return ret.Result, &ret.Pagination, nil
}

// GetWebSocketToken Get the token for websocket authentication
func (b *bitkubApi) GetWebSocketToken() (string, error) {
	if b.ApiKey == "" {
		return "", fmt.Errorf("api key is empty")
	}
	if b.ApiSecret == "" {
		return "", fmt.Errorf("api secret is empty")
	}

	payloadBytes, err := b.sigPayload(map[string]interface{}{})
	if err != nil {
		return "", err
	}

	url := baseURL + "/api/market/wstoken"
	headers := map[string]string{"X-BTK-APIKEY": b.ApiKey, "Content-Type": "application/json", "Accept": "application/json"}
	resp, err := internal.Post(url, headers, payloadBytes, b.Timeout)
	if err != nil {
		return "", err
	}
	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("got server error (%d) : %s", resp.StatusCode(), resp.String())
	}

	ret := model.DefaultResponse{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return "", err
	}

	if ret.Error != 0 {
		return "", fmt.Errorf("got server error (%d) : %s", ret.Error, internal.GetErrorMessage(ret.Error))
	}

	return ret.Result.(string), nil
}

// GetUserLimits Check deposit/withdraw limitations and usage.
func (b *bitkubApi) GetUserLimits() (*model.UserLimits, error) {
	if b.ApiKey == "" {
		return nil, fmt.Errorf("api key is empty")
	}
	if b.ApiSecret == "" {
		return nil, fmt.Errorf("api secret is empty")
	}

	payloadBytes, err := b.sigPayload(map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	url := baseURL + "/api/user/limits"
	headers := map[string]string{"X-BTK-APIKEY": b.ApiKey, "Content-Type": "application/json", "Accept": "application/json"}
	resp, err := internal.Post(url, headers, payloadBytes, b.Timeout)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("got server error (%d) : %s", resp.StatusCode(), resp.String())
	}

	ret := model.UserLimitsResponse{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return nil, err
	}

	if ret.Error != 0 {
		return nil, fmt.Errorf("got server error (%d) : %s", ret.Error, internal.GetErrorMessage(ret.Error))
	}

	return &ret.Result, nil
}

// GetUserTradingCredits Check trading credit balance.
func (b *bitkubApi) GetUserTradingCredits() (float64, error) {
	if b.ApiKey == "" {
		return 0, fmt.Errorf("api key is empty")
	}
	if b.ApiSecret == "" {
		return 0, fmt.Errorf("api secret is empty")
	}

	payloadBytes, err := b.sigPayload(map[string]interface{}{})
	if err != nil {
		return 0, err
	}

	url := baseURL + "/api/user/trading-credits"
	headers := map[string]string{"X-BTK-APIKEY": b.ApiKey, "Content-Type": "application/json", "Accept": "application/json"}
	resp, err := internal.Post(url, headers, payloadBytes, b.Timeout)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode() != 200 {
		return 0, fmt.Errorf("got server error (%d) : %s", resp.StatusCode(), resp.String())
	}

	ret := model.DefaultResponse{}
	if err := json.Unmarshal(resp.Body(), &ret); err != nil {
		return 0, err
	}

	if ret.Error != 0 {
		return 0, fmt.Errorf("got server error (%d) : %s", ret.Error, internal.GetErrorMessage(ret.Error))
	}

	return ret.Result.(float64), nil
}
