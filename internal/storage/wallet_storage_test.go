package storage_test

import (
	"bl-wallet-service/internal/models"
	"bl-wallet-service/internal/storage"
	"bl-wallet-service/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWalletStorage_UpdateBalance(t *testing.T) {
	type args struct {
		wallet      *storage.RowWallet
		transaction *models.TransactionRequest
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		expectedErr error
	}{
		{
			name: "SuccessfulCreditTransaction",
			args: args{
				wallet: &storage.RowWallet{ID: "1", UserID: "100", Balance: 100, WalletVersion: 1},
				transaction: &models.TransactionRequest{
					TransactionID:   "credit_transaction_id",
					UserID:          "100",
					Amount:          20,
					TransactionType: models.CreditTransactionType,
				},
			},
			wantErr: false,
		},
		{
			name: "SuccessfulDebitTransaction",
			args: args{
				wallet: &storage.RowWallet{ID: "2", UserID: "200", Balance: 100, WalletVersion: 1},
				transaction: &models.TransactionRequest{
					TransactionID:   "debit_transaction_id",
					UserID:          "200",
					Amount:          20,
					TransactionType: models.DebitTransactionType,
				},
			},
			wantErr: false,
		},
		{
			name: "FailedDebitTransactionInsufficientFunds",
			args: args{
				wallet: &storage.RowWallet{ID: "3", UserID: "300", Balance: 50, WalletVersion: 1},
				transaction: &models.TransactionRequest{
					TransactionID:   "insufficient_funds_transaction_id",
					UserID:          "300",
					Amount:          60,
					TransactionType: models.DebitTransactionType,
				},
			},
			wantErr:     true,
			expectedErr: models.ErrInsufficientFunds,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()
			mockStorage := mocks.NewIWalletStorage(t)

			if tt.wantErr {
				mockStorage.On("UpdateBalance", tt.args.wallet, tt.args.transaction).Return(tt.expectedErr)
			} else {
				mockStorage.On("UpdateBalance", tt.args.wallet, tt.args.transaction).Return(nil)
			}

			err := mockStorage.UpdateBalance(tt.args.wallet, tt.args.transaction)
			assert.Equal(t, tt.expectedErr, err)

			mockStorage.AssertExpectations(t)
		})
	}
}
