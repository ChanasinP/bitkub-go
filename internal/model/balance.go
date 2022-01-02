package model

type Balance struct {
	Available float64 `json:"available"`
	Reserved  float64 `json:"reserved"`
}
