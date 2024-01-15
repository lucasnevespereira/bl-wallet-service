package handlers

import (
	"bl-wallet-service/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_validateFundsRequest(t *testing.T) {
	type args struct {
		request models.TransactionRequest
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		expectedErr error
	}{
		{
			name: "Request with empty user id",
			args: args{request: models.TransactionRequest{
				TransactionID: "transaction_id",
				UserID:        "",
				Amount:        20.0,
			}},
			wantErr:     true,
			expectedErr: models.ErrEmptyUserID,
		},
		{
			name: "Request with empty transaction id",
			args: args{request: models.TransactionRequest{
				TransactionID: "",
				UserID:        "2",
				Amount:        10.0,
			}},
			wantErr:     true,
			expectedErr: models.ErrEmptyTransactionID,
		},
		{
			name: "Request with negative amount",
			args: args{request: models.TransactionRequest{
				TransactionID: "transaction_id",
				UserID:        "1",
				Amount:        -1,
			}},
			wantErr:     true,
			expectedErr: models.ErrAmountNegative,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTransactionRequest(tt.args.request)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
