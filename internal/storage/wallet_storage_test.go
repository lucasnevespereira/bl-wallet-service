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
		wallet        *storage.RowWallet
		transactionID string
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
				wallet:        &storage.RowWallet{ID: "1", UserID: "100", Balance: 100, WalletVersion: 1},
				transactionID: "credit_transaction_id",
			},
			wantErr: false,
		},
		{
			name: "SuccessfulDebitTransaction",
			args: args{
				wallet:        &storage.RowWallet{ID: "2", UserID: "200", Balance: 100, WalletVersion: 1},
				transactionID: "debit_transaction_id",
			},
			wantErr: false,
		},
		{
			name: "FailedDebitTransactionInsufficientFunds",
			args: args{
				wallet:        &storage.RowWallet{ID: "3", UserID: "300", Balance: 50, WalletVersion: 1},
				transactionID: "insufficient_funds_transaction_id",
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
				mockStorage.On("UpdateBalance", tt.args.wallet, tt.args.transactionID).Return(tt.expectedErr)
			} else {
				mockStorage.On("UpdateBalance", tt.args.wallet, tt.args.transactionID).Return(nil)
			}

			err := mockStorage.UpdateBalance(tt.args.wallet, tt.args.transactionID)
			assert.Equal(t, tt.expectedErr, err)

			mockStorage.AssertExpectations(t)
		})
	}
}
