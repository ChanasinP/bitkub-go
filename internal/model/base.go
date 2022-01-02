package model

type Pagination struct {
	Page int `json:"page"`
	Last int `json:"last"`
}

type ErrorResponse struct {
	Error int `json:"error"`
}

type DefaultResponse struct {
	Error  int         `json:"error"`
	Result interface{} `json:"result"`
}
