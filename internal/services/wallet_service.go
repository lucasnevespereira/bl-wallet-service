package services

import (
	"bl-wallet-service/internal/models"
	"bl-wallet-service/internal/storage"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"log"
	"time"
)

type IWalletService interface {
	Create(userID string) error
	GetByUserID(userID string) (*models.Wallet, error)
	ProcessTransaction(request *models.TransactionRequest) error
}

type WalletService struct {
	storage storage.IWalletStorage
}

func NewWalletService(walletStorage storage.IWalletStorage) *WalletService {
	return &WalletService{
		storage: walletStorage,
	}
}

func (s *WalletService) GetByUserID(userID string) (*models.Wallet, error) {
	fmt.Printf("Fetching wallet for user %s \n", userID)
	rowWallet, err := s.storage.GetWalletByUserID(userID)
	if err != nil {
		return nil, err
	}
	return &models.Wallet{
		ID:            rowWallet.ID,
		UserID:        rowWallet.UserID,
		Balance:       rowWallet.Balance,
		WalletVersion: rowWallet.WalletVersion,
	}, nil
}
func (s *WalletService) Create(userID string) error {

	userWallet, err := s.storage.GetWalletByUserID(userID)
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return err
	}

	if userWallet != nil {
		return errors.New(fmt.Sprintf("User %s already has a wallet", userID))
	}

	fmt.Printf("Creating wallet for user %s \n", userID)
	err = s.storage.CreateWallet(&storage.RowWallet{
		ID:            uuid.New().String(),
		UserID:        userID,
		Balance:       0.0,
		WalletVersion: 0,
	})
	if err != nil {
		log.Printf("storage.CreateWallet: %v \n", err)
		return err
	}
	return nil
}

func (s *WalletService) ProcessTransaction(request *models.TransactionRequest) error {
	fmt.Printf("Processing transaction of user %s \n", request.UserID)

	// check if transaction exists
	existingTransaction, err := s.storage.GetTransaction(request.TransactionID)
	if err != nil {
		return err
	}
	if existingTransaction != nil {
		return models.ErrTransactionAlreadyProcessed
	}

	userWallet, err := s.storage.GetWalletByUserID(request.UserID)
	if err != nil {
		return err
	}

	if userWallet == nil {
		return models.ErrUserWalletNotFound
	}

	// update wallet balance
	err = s.storage.UpdateBalance(userWallet, request)
	if err != nil && errors.Is(err, models.ErrInsufficientFunds) {
		// record transaction as failed
		now := time.Now().Format(time.RFC3339Nano)
		err = s.storage.CreateTransaction(&storage.RowTransaction{
			ID:              request.TransactionID,
			UserID:          request.UserID,
			WalletID:        userWallet.ID,
			Amount:          request.Amount,
			TransactionType: request.TransactionType,
			Status:          models.FailedTransactionStatus,
			CreatedAt:       now,
			UpdatedAt:       now,
		})
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}
