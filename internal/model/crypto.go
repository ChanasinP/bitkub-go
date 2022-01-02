package model

type CryptoAddress struct {
	Currency  string `json:"currency"`
	Address   string `json:"address"`
	Tag       int    `json:"tag"`
	Timestamp int64  `json:"time"`
}

type CryptoAddressResponse struct {
	Error      int             `json:"error"`
	Result     []CryptoAddress `json:"result"`
	Pagination Pagination      `json:"pagination"`
}

type CryptoWithdraw struct {
	TxnID     string  `json:"txn"` // local transaction id
	Address   string  `json:"adr"` // address
	Memo      string  `json:"mem"` // memo
	Currency  string  `json:"cur"` // currency
	Amount    float64 `json:"amt"` // withdraw amount
	Fee       float64 `json:"fee"` // withdraw fee
	Timestamp int64   `json:"ts"`  // timestamp
}

type CryptoWithdrawResponse struct {
	Error  int            `json:"error"`
	Result CryptoWithdraw `json:"result"`
}

type CryptoWithdrawHistoryResponse struct {
	Error      int              `json:"error"`
	Result     []CryptoWithdraw `json:"result"`
	Pagination Pagination       `json:"pagination"`
}

type CryptoDeposit struct {
	Hash          string  `json:"hash"`     // transaction hash
	Currency      string  `json:"currency"` // currency
	Amount        float64 `json:"amount"`   // amount
	FromAddress   string  `json:"from_address"`
	ToAddress     string  `json:"to_address"`
	Confirmations int     `json:"confirmations"`
	Status        string  `json:"status"`
	Timestamp     int64   `json:"time"`
}

type CryptoDepositResponse struct {
	Error      int             `json:"error"`
	Result     []CryptoDeposit `json:"result"`
	Pagination Pagination      `json:"pagination"`
}

type CryptoGenerateAddress struct {
	Currency string `json:"currency"`
	Address  string `json:"address"`
	Memo     string `json:"mem"`
}

type CryptoGenerateAddressResponse struct {
	Error  int                     `json:"error"`
	Result []CryptoGenerateAddress `json:"result"`
}
