package model

type Limit struct {
	Deposit  float64 `json:"deposit"`  // value equivalent
	Withdraw float64 `json:"withdraw"` // value equivalent
}

// Limits limitations by kyc level
type Limits struct {
	Crypto Limit `json:"crypto"` // BTC limit
	Fiat   Limit `json:"fiat"`   // THB limit
}

type CryptoUsage struct {
	Deposit               float64 `json:"deposit"`  // BTC value equivalent
	Withdraw              float64 `json:"withdraw"` // BTC value equivalent
	DepositPercentage     float64 `json:"deposit_percentage"`
	WithdrawPercentage    float64 `json:"withdraw_percentage"`
	DepositTHBEquivalent  float64 `json:"deposit_thb_equivalent"`  // THB value equivalent
	WithdrawTHBEquivalent float64 `json:"withdraw_thb_equivalent"` // THB value equivalent
}

type FiatUsage struct {
	Deposit            float64 `json:"deposit"`  // THB value equivalent
	Withdraw           float64 `json:"withdraw"` // THB value equivalent
	DepositPercentage  float64 `json:"deposit_percentage"`
	WithdrawPercentage float64 `json:"withdraw_percentage"`
}

type UserLimits struct {
	Limits Limits `json:"limits"`
	Usage  struct {
		Crypto CryptoUsage `json:"crypto"`
		Fiat   FiatUsage   `json:"fiat"`
	} `json:"usage"`
	Rate float64 `json:"rate"` // current THB rate used to calculate
}

type UserLimitsResponse struct {
	Error  int        `json:"error"`
	Result UserLimits `json:"result"`
}
