package model

type Order struct {
	ID        int64   `json:"id"`   // order id
	Hash      string  `json:"hash"` // order hash
	Type      string  `json:"typ"`  // order type
	Amount    float64 `json:"amt"`  // spending amount
	Rate      float64 `json:"rat"`  // rate
	Fee       float64 `json:"fee"`  // fee
	Credit    float64 `json:"cre"`  // credit used
	Receive   float64 `json:"rec"`  // amount to receive
	Timestamp int64   `json:"ts"`   // timestamp
}

type OrderResponse struct {
	Error  int   `json:"error"`
	Result Order `json:"result"`
}

type OpenOrder struct {
	ID        int     `json:"id"`        // order id
	Hash      string  `json:"hash"`      // order hash
	Side      string  `json:"side"`      // order side: buy or sell
	Type      string  `json:"type"`      // order type
	Rate      float64 `json:"rate"`      // rate
	Fee       float64 `json:"fee"`       // fee
	Credit    float64 `json:"credit"`    // credit used
	Amount    float64 `json:"amount"`    // amount
	Receive   float64 `json:"receive"`   // amount to receive
	ParentID  int     `json:"parent_id"` // parent order id
	SuperID   int     `json:"super_id"`  // super parent order id
	Timestamp int64   `json:"ts"`        // timestamp
}

type OpenOrderResponse struct {
	Error  int         `json:"error"`
	Result []OpenOrder `json:"result"`
}

type OrderHistory struct {
	TxnID         string  `json:"txn_id"`
	OrderID       int     `json:"order_id"`
	Hash          string  `json:"hash"`
	ParentOrderID int     `json:"parent_order_id"`
	SuperOrderID  int     `json:"super_order_id"`
	TakenByMe     bool    `json:"taken_by_me"`
	IsMaker       bool    `json:"is_maker"`
	Side          string  `json:"side"`
	Type          string  `json:"type"`
	Rate          float64 `json:"rate"`
	Fee           float64 `json:"fee"`
	Credit        float64 `json:"credit"`
	Amount        float64 `json:"amount"`
	Receive       float64 `json:"receive"`
}

type OrderHistoryPagination struct {
	Page     int `json:"page"`
	Last     int `json:"last"`
	Next     int `json:"next"`
	Previous int `json:"prev"`
}

type OrderHistoryResponse struct {
	Error      int                    `json:"error"`
	Result     []OrderHistory         `json:"result"`
	Pagination OrderHistoryPagination `json:"pagination"`
}

type OrderInfoHistory struct {
	Amount    float64 `json:"amount"`
	Credit    float64 `json:"credit"`
	Fee       float64 `json:"fee"`
	ID        int64   `json:"id"`
	Rate      float64 `json:"rate"`
	Timestamp int64   `json:"timestamp"`
}

type OrderInfo struct {
	ID            int64              `json:"id"`
	First         int64              `json:"first"`          // first order id
	Parent        int64              `json:"parent"`         // parent order id
	Last          int64              `json:"last"`           // last order id
	Amount        float64            `json:"amount"`         // order amount
	Rate          float64            `json:"rate"`           // order rate
	Fee           float64            `json:"fee"`            // order fee
	Credit        float64            `json:"credit"`         // order fee credit used
	Filled        float64            `json:"filled"`         // filled amount
	Total         float64            `json:"total"`          // total amount
	Status        string             `json:"status"`         // order status: filled, unfilled
	PartialFilled bool               `json:"partial_filled"` // true when order has been partially filled, false when not filled or fully filled
	Remaining     float64            `json:"remaining"`      // remaining amount to be executed
	History       []OrderInfoHistory `json:"history"`        // order history
}

type OrderInfoResponse struct {
	Error  int       `json:"error"`
	Result OrderInfo `json:"result"`
}
