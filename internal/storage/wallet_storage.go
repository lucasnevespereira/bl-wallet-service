package storage

import (
	"bl-wallet-service/internal/models"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

type IWalletStorage interface {
	AutoMigrate() error
	CreateWallet(wallet *RowWallet) error
	GetWalletByUserID(userID string) (*RowWallet, error)
	UpdateBalance(wallet *RowWallet, transactionID string) error
	GetTransaction(transactionID string) (*RowTransaction, error)
	CreateTransaction(transaction *RowTransaction) error
}

type RowWallet struct {
	ID      string
	UserID  string
	Balance float64
}

func (RowWallet) TableName() string {
	return "wallets"
}

type RowTransaction struct {
	ID              string
	UserID          string
	WalletID        string
	Amount          float64
	TransactionType string
	Status          string
	CreatedAt       string
	UpdatedAt       string
}

func (RowTransaction) TableName() string {
	return "transactions"
}

type WalletStorage struct {
	db *gorm.DB
}

func (s *WalletStorage) AutoMigrate() error {
	return s.db.AutoMigrate(&RowWallet{})
}

type WalletStorageConfig struct {
	DbHost     string
	DbPort     int
	DbUser     string
	DbPassword string
	DbName     string
	DbSsl      string
}

func NewWalletStorage(config WalletStorageConfig) (*WalletStorage, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.DbHost,
		config.DbPort,
		config.DbUser,
		config.DbPassword,
		config.DbName,
		config.DbSsl,
	)

	database, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, errors.Wrapf(err, "could not create postgres client")
	}

	internalDB, errInternalDB := database.DB()
	if errInternalDB != nil {
		return nil, errors.Wrapf(errInternalDB, "could not get internal db")
	}

	if errPing := internalDB.Ping(); errPing != nil {
		return nil, errors.Wrapf(errPing, "could not ping database")
	}

	log.Println("Wallet Storage started")
	return &WalletStorage{db: database}, nil
}

func (s *WalletStorage) CreateWallet(wallet *RowWallet) error {
	result := s.db.Create(wallet)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *WalletStorage) GetWalletByUserID(userId string) (*RowWallet, error) {
	var rowWallet *RowWallet
	result := s.db.First(&rowWallet, "user_id = ?", userId)
	if result.Error != nil {
		return nil, result.Error
	}
	return rowWallet, nil
}

func (s *WalletStorage) GetTransaction(transactionID string) (*RowTransaction, error) {
	var rowTransaction *RowTransaction
	result := s.db.First(&rowTransaction, "id = ?", transactionID)
	if result.Error != nil {
		return nil, result.Error
	}
	return rowTransaction, nil
}

func (s *WalletStorage) CreateTransaction(transaction *RowTransaction) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return s.db.Model(&RowTransaction{}).Create(&transaction).Error
	})
}

func (s *WalletStorage) AddFunds(userID string, amount float64) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return s.db.Model(&RowWallet{}).Where("user_id = ?", userID).Update("balance", gorm.Expr("balance + ?", amount)).Error
	})
}

func (s *WalletStorage) RemoveFunds(userID string, amount float64) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return s.db.Model(&RowWallet{}).Where("user_id = ?", userID).Update("balance", gorm.Expr("balance - ?", amount)).Error
	})
}

func (s *WalletStorage) UpdateBalance(wallet *RowWallet, transactionID string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now().Format(time.RFC3339Nano)
		var updatedTransaction = RowTransaction{
			UpdatedAt: now,
		}

		//get transaction by id
		transaction, err := s.GetTransaction(transactionID)
		if err != nil {
			return err
		}

		// check if transaction is pending
		if transaction.Status != models.PendingTransactionStatus {
			return errors.New(models.TRANSACTION_ALREADY_PROCESSED_ERROR)
		}

		// check transaction type
		if transaction.TransactionType == models.CreditTransactionType {
			err := s.AddFunds(transaction.UserID, transaction.Amount)
			if err != nil {
				updatedTransaction.Status = models.FailedTransactionStatus
				err := s.db.Model(&RowTransaction{}).Where("id = ?", transaction.ID).Updates(updatedTransaction).Error
				if err != nil {
					return err
				}
				return err
			}

		} else if transaction.TransactionType == models.DebitTransactionType {
			if wallet.Balance < transaction.Amount {
				updatedTransaction.Status = models.FailedTransactionStatus
				err := s.db.Model(&RowTransaction{}).Where("id = ?", transaction.ID).Updates(updatedTransaction).Error
				if err != nil {
					return err
				}
				return errors.New(models.INSUFFICIENT_FUNDS_ERROR)
			}

			err := s.RemoveFunds(transaction.UserID, transaction.Amount)
			if err != nil {
				updatedTransaction.Status = models.FailedTransactionStatus
				err := s.db.Model(&RowTransaction{}).Where("id = ?", transaction.ID).Updates(updatedTransaction).Error
				if err != nil {
					return err
				}
				return err
			}

		}

		updatedTransaction.Status = models.SuccessTransactionStatus
		return s.db.Model(&RowTransaction{}).Where("id = ?", transaction.ID).Updates(updatedTransaction).Error
	})
}