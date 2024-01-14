package models

type CreateWalletRequest struct {
	UserID string `json:"userID"`
}

type WalletFundsRequest struct {
	UserID string  `json:"userID"`
	Amount float64 `json:"amount"`
}
