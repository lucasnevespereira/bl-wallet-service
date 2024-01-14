package models

type Wallet struct {
	ID            string  `json:"id"`
	UserID        string  `json:"userID"`
	Balance       float64 `json:"balance"`
	WalletVersion int     `json:"walletVersion"`
}
