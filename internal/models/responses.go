package models

type CreateWalletResponse struct {
	UserID  string `json:"userID"`
	Message string `json:"message"`
}

type TransactionResponse struct {
	UserID  string `json:"userID"`
	Message string `json:"message"`
}
