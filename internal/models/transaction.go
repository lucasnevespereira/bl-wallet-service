package models

import "time"

type Transaction struct {
	ID              string     `json:"id"`
	WalletID        string     `json:"walletID"`
	UserID          string     `json:"userID"`
	Amount          float64    `json:"amount"`
	TransactionType string     `json:"transactionType"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       *time.Time `json:"updatedAt"`
}

const (
	PendingTransactionStatus = "PENDING"
	SuccessTransactionStatus = "SUCCESS"
	FailedTransactionStatus  = "FAILED"
)

const (
	CreditTransactionType = "CREDIT"
	DebitTransactionType  = "DEBIT"
)
