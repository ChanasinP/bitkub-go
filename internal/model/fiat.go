package model

type BankAccount struct {
	ID        string `json:"id"`
	Bank      string `json:"bank"`
	Name      string `json:"name"`
	Timestamp int64  `json:"time"`
}

type FiatWithdraw struct {
	TxnID     string  `json:"txn"` // local transaction id
	AccountID string  `json:"acc"` // bank account id
	Currency  string  `json:"cur"` // currency
	Amount    float64 `json:"amt"` // withdraw amount
	Fee       float64 `json:"fee"` // withdraw fee
	Receive   float64 `json:"rec"` // amount to receive
	Timestamp int64   `json:"ts"`  // timestamp
}

type FiatDeposit struct {
	TxnID     string  `json:"txn_id"`   // local transaction id
	Currency  string  `json:"currency"` // currency
	Amount    float64 `json:"amount"`   // deposit amount
	Status    string  `json:"status"`
	Timestamp int64   `json:"time"` // timestamp
}

type FiatAccountsResponse struct {
	Error      int           `json:"error"`
	Result     []BankAccount `json:"result"`
	Pagination Pagination    `json:"pagination"`
}

type FiatWithdrawResponse struct {
	Error  int          `json:"error"`
	Result FiatWithdraw `json:"result"`
}

type FiatDepositResponse struct {
	Error      int           `json:"error"`
	Result     []FiatDeposit `json:"result"`
	Pagination Pagination    `json:"pagination"`
}

type FiatWithdrawHistoryResponse struct {
	Error      int            `json:"error"`
	Result     []FiatWithdraw `json:"result"`
	Pagination Pagination     `json:"pagination"`
}
