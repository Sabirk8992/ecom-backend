package model

type Payment struct {
	ID        int     `json:"id"`
	OrderID   int     `json:"order_id"`
	Amount    float64 `json:"amount"`
	Status    string  `json:"status"`
	Method    string  `json:"method"`
	CreatedAt string  `json:"created_at"`
}

type PaymentRequest struct {
	OrderID int    `json:"order_id"`
	Method  string `json:"method"` // "card", "upi", "wallet"
	// SimulateFailure lets us test failure cases
	SimulateFailure bool `json:"simulate_failure"`
}

type PaymentResponse struct {
	Payment *Payment `json:"payment"`
	Message string   `json:"message"`
}
