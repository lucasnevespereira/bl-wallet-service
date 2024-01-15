package services_test

import (
	"bl-wallet-service/internal/models"
	"bl-wallet-service/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWalletService_ProcessTransaction(t *testing.T) {
	type args struct {
		request *models.TransactionRequest
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		expectedErr error
	}{
		{
			name: "TransactionAlreadyProcessed",
			args: args{
				request: &models.TransactionRequest{
					TransactionID:   "transaction1",
					UserID:          "user1",
					Amount:          50,
					TransactionType: models.CreditTransactionType,
				},
			},
			wantErr:     true,
			expectedErr: models.ErrTransactionAlreadyProcessed,
		},

		{
			name: "UserWalletNotFound",
			args: args{
				request: &models.TransactionRequest{
					TransactionID:   "transaction2",
					UserID:          "user2",
					Amount:          30,
					TransactionType: models.DebitTransactionType,
				},
			},
			wantErr:     true,
			expectedErr: models.ErrUserWalletNotFound,
		},

		{
			name: "InsufficientFunds",
			args: args{
				request: &models.TransactionRequest{
					TransactionID:   "transaction3",
					UserID:          "user3",
					Amount:          100,
					TransactionType: models.DebitTransactionType,
				},
			},
			wantErr:     true,
			expectedErr: models.ErrInsufficientFunds,
		},

		{
			name: "SuccessfulDepositTransaction",
			args: args{
				request: &models.TransactionRequest{
					TransactionID:   "transaction4",
					UserID:          "user4",
					Amount:          25,
					TransactionType: models.CreditTransactionType,
				},
			},
			wantErr: false,
		},

		{
			name: "SuccessfulWithdrawTransaction",
			args: args{
				request: &models.TransactionRequest{
					TransactionID:   "transaction5",
					UserID:          "user5",
					Amount:          75,
					TransactionType: models.DebitTransactionType,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()
			mockService := mocks.NewIWalletService(t)
			if tt.wantErr {
				mockService.On("ProcessTransaction", tt.args.request).Return(tt.expectedErr)
			} else {
				mockService.On("ProcessTransaction", tt.args.request).Return(nil)
			}
			err := mockService.ProcessTransaction(tt.args.request)
			assert.Equal(t, tt.expectedErr, err)
			mockService.AssertExpectations(t)
		})
	}
}
