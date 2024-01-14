package models

type CreateWalletResponse struct {
	UserID  string `json:"userID"`
	Message string `json:"message"`
}

type WalletFundsResponse struct {
	UserID  string `json:"userID"`
	Message string `json:"message"`
}
