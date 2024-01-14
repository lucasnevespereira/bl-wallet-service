package services

//func TestWalletService_ProcessTransaction(t *testing.T) {
//	mockStorage := mocks.NewIWalletStorage(t)
//	mockCache := mocks.NewITransactionCache(t)
//	service := &WalletService{
//		storage: mockStorage,
//		cache:   mockCache,
//	}
//	testCases := []struct {
//		name            string
//		userID          string
//		amount          float64
//		balance         float64
//		transactionType string
//		wantErr         bool
//		errMessage      string
//	}{
//		{
//			name:            "Sufficient Funds",
//			userID:          "1",
//			amount:          10.0,
//			balance:         20.0,
//			transactionType: "DEBIT",
//			wantErr:         false,
//		},
//		{
//			name:            "Insufficient Funds",
//			userID:          "1",
//			amount:          30.0,
//			balance:         20.0,
//			transactionType: "DEBIT",
//			wantErr:         true,
//			errMessage:      "Insufficient funds",
//		},
//	}
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			mockStorage.On("GetWalletByUserID", tc.userID).Return(&storage.RowWallet{Balance: tc.balance}, nil)
//
//			if tc.balance >= tc.amount {
//				mockStorage.On("RemoveFunds", tc.userID, tc.amount).Return(nil).Once()
//			}
//
//			err := service.ProcessTransaction(tc.userID, tc.amount)
//			if tc.wantErr {
//				assert.Error(t, err)
//				assert.Contains(t, err.Error(), tc.errMessage)
//			} else {
//				assert.NoError(t, err)
//			}
//
//			mockStorage.AssertExpectations(t)
//		})
//	}
//}
