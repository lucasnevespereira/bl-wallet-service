package models

type CreateWalletRequest struct {
	UserID string `json:"userID"`
}

type TransactionRequest struct {
	TransactionID string  `json:"transactionID"`
	UserID        string  `json:"userID"`
	Amount        float64 `json:"amount"`
}
