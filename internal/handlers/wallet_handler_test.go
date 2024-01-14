package handlers

import (
	"bl-wallet-service/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_validateFundsRequest(t *testing.T) {
	type args struct {
		request models.WalletFundsRequest
	}
	testCases := []struct {
		name       string
		args       args
		wantErr    bool
		errMessage string
	}{
		{
			name: "Request with empty user id",
			args: args{request: models.WalletFundsRequest{
				UserID: "",
				Amount: 0,
			}},
			wantErr:    true,
			errMessage: "You should provide a user id",
		},
		{
			name: "Request with negative amount",
			args: args{request: models.WalletFundsRequest{
				UserID: "1",
				Amount: -1,
			}},
			wantErr:    true,
			errMessage: "Amount should not be negative",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateFundsRequest(tc.args.request)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
