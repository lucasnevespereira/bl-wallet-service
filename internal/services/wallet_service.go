package services

import (
	"bl-wallet-service/internal/cache"
	"bl-wallet-service/internal/models"
	"bl-wallet-service/internal/storage"
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

type IWalletService interface {
	Create(userID string) error
	GetByUserID(userID string) (*models.Wallet, error)
	ProcessTransaction(userID, transactionType string, amount float64) error
}

type WalletService struct {
	storage storage.IWalletStorage
	cache   cache.ITransactionCache
	mu      sync.Mutex
}

func NewWalletService(walletStorage storage.IWalletStorage, transactionCache cache.ITransactionCache) *WalletService {
	return &WalletService{
		storage: walletStorage,
		cache:   transactionCache,
	}
}

func (s *WalletService) GetByUserID(userID string) (*models.Wallet, error) {
	fmt.Printf("Fetching wallet for user %s \n", userID)
	rowWallet, err := s.storage.GetWalletByUserID(userID)
	if err != nil {
		return nil, err
	}
	return &models.Wallet{
		ID:      rowWallet.ID,
		UserID:  rowWallet.UserID,
		Balance: rowWallet.Balance,
	}, nil
}
func (s *WalletService) Create(userID string) error {
	fmt.Printf("Creating wallet for user %s \n", userID)
	err := s.storage.CreateWallet(&storage.RowWallet{
		ID:      uuid.New().String(),
		UserID:  userID,
		Balance: 0.0,
	})
	if err != nil {
		return err
	}
	return nil
}
func (s *WalletService) ProcessTransaction(userID, transactionType string, amount float64) error {
	fmt.Printf("Processing transaction of user %s \n", userID)
	s.mu.Lock()
	defer s.mu.Unlock()
	userWallet, err := s.storage.GetWalletByUserID(userID)
	if err != nil {
		return err
	}
	now := time.Now()

	// create transaction
	transactionID := uuid.New().String()
	err = s.storage.CreateTransaction(&storage.RowTransaction{
		ID:              transactionID,
		UserID:          userID,
		WalletID:        userWallet.ID,
		Amount:          amount,
		TransactionType: transactionType,
		Status:          models.PendingTransactionStatus,
		CreatedAt:       now.Format(time.RFC3339Nano),
		UpdatedAt:       "",
	})
	if err != nil {
		return err
	}

	// update wallet balance
	err = s.storage.UpdateBalance(userWallet, transactionID)
	if err != nil {
		return err
	}

	return nil
}
