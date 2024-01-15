package storage

import (
	"bl-wallet-service/internal/models"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"time"
)

type IWalletStorage interface {
	AutoMigrate() error
	CreateWallet(wallet *RowWallet) error
	GetWalletByUserID(userID string) (*RowWallet, error)
	UpdateBalance(wallet *RowWallet, transaction *models.TransactionRequest) error
	GetTransaction(transactionID string) (*RowTransaction, error)
	CreateTransaction(transaction *RowTransaction) error
}

type RowWallet struct {
	ID            string
	UserID        string
	Balance       float64
	WalletVersion int
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
	DB *gorm.DB
}

func (s *WalletStorage) AutoMigrate() error {
	return s.DB.AutoMigrate(&RowWallet{})
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
	return &WalletStorage{DB: database}, nil
}

func (s *WalletStorage) CreateWallet(wallet *RowWallet) error {
	log.Printf("Wallet to insert: %v \n", wallet)
	if s.DB == nil {
		log.Println("database nil")
	}
	err := s.DB.Create(&wallet).Error
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Wallet created in database")
	return nil
}

func (s *WalletStorage) GetWalletByUserID(userId string) (*RowWallet, error) {
	var rowWallet *RowWallet
	result := s.DB.First(&rowWallet, "user_id = ?", userId)
	if result.Error != nil {
		return nil, result.Error
	}
	return rowWallet, nil
}

func (s *WalletStorage) GetTransaction(transactionID string) (*RowTransaction, error) {
	var rowTransaction *RowTransaction
	err := s.DB.First(&rowTransaction, "id = ?", transactionID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return rowTransaction, nil
}

func (s *WalletStorage) CreateTransaction(transaction *RowTransaction) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		log.Printf("inserting transaction: %v \n", transaction)
		err := tx.Model(&RowTransaction{}).Create(&transaction).Error
		if err != nil {
			return err
		}

		log.Printf("transaction inserted: %v \n", transaction)
		return nil
	})
}

func (s *WalletStorage) UpdateBalance(wallet *RowWallet, transaction *models.TransactionRequest) error {
	log.Println("updating wallet balance")
	return s.DB.Transaction(func(tx *gorm.DB) error {

		// locking wallet
		var isolatedWallet *RowWallet
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? and wallet_version = ?", wallet.ID, wallet.WalletVersion).Take(&isolatedWallet).Error
		if err != nil {
			log.Printf("tx.Clauses locking err: %v \n", err)
			return err
		}

		newBalance := isolatedWallet.Balance

		// check transaction type
		if transaction.TransactionType == models.CreditTransactionType {
			newBalance += transaction.Amount
		} else if transaction.TransactionType == models.DebitTransactionType {
			if wallet.Balance < transaction.Amount {
				return models.ErrInsufficientFunds
			}
			newBalance -= transaction.Amount
		}

		newWalletVersion := isolatedWallet.WalletVersion + 1
		if err := tx.Model(&isolatedWallet).Updates(RowWallet{Balance: newBalance, WalletVersion: newWalletVersion}).Error; err != nil {
			return err
		}

		now := time.Now().Format(time.RFC3339Nano)
		transactionRow := RowTransaction{
			ID:              transaction.TransactionID,
			UserID:          transaction.UserID,
			WalletID:        isolatedWallet.ID,
			Amount:          transaction.Amount,
			TransactionType: transaction.TransactionType,
			Status:          models.SuccessTransactionStatus,
			CreatedAt:       now,
			UpdatedAt:       now,
		}
		err = tx.Create(&transactionRow).Error
		if err != nil {
			return err
		}

		log.Println("wallet balance updated")
		return nil
	})
}
