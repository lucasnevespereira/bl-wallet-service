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
	UpdateBalance(wallet *RowWallet, transactionID string) error
	GetTransaction(transactionID string) (*RowTransaction, error)
	CreateTransaction(transaction *RowTransaction) error
	UpdateTransactionStatus(transactionID string, status string) error
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

func (s *WalletStorage) UpdateTransactionStatus(transactionID string, status string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("UPDATE transactions SET status = ?, updated_at = ? WHERE id = ?", status, time.Now().Format(time.RFC3339Nano), transactionID).Error; err != nil {
			return err
		}

		log.Printf("transaction %v updated to status %s \n", transactionID, status)
		return nil
	})
}

func (s *WalletStorage) UpdateBalance(wallet *RowWallet, transactionID string) error {
	log.Println("updating wallet balance")
	return s.DB.Transaction(func(tx *gorm.DB) error {

		transaction, err := s.GetTransaction(transactionID)
		if err != nil {
			return err
		}

		if transaction == nil {
			return models.ErrTransactionNotFound
		}

		// check if transaction is pending
		if transaction.Status != models.PendingTransactionStatus {
			return models.ErrTransactionAlreadyProcessed
		}

		// locking wallet
		var isolatedWallet *RowWallet
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? and wallet_version = ?", wallet.ID, wallet.WalletVersion).Take(&isolatedWallet).Error
		if err != nil {
			log.Printf("tx.Clauses locking err: %v \n", err)
			return err
		}

		newBalance := isolatedWallet.Balance
		// check transaction type
		if transaction.TransactionType == models.CreditTransactionType {
			newBalance += transaction.Amount
		} else if transaction.TransactionType == models.DebitTransactionType {
			newBalance = isolatedWallet.Balance - transaction.Amount
		}

		if newBalance < 0 {
			return models.ErrInsufficientFunds
		}

		newWalletVersion := isolatedWallet.WalletVersion + 1
		// update balance
		if err := tx.Exec("UPDATE wallets SET balance = ?, wallet_version = ? WHERE id = ?", newBalance, newWalletVersion, isolatedWallet.ID).Error; err != nil {
			log.Printf("tx wallet update: %v\n", err)
			return err
		}

		now := time.Now().Format(time.RFC3339Nano)
		transactionRow := RowTransaction{
			Status:    models.SuccessTransactionStatus,
			UpdatedAt: now,
		}

		err = tx.Model(&RowTransaction{}).Where("id = ?", transaction.ID).Updates(&transactionRow).Error
		if err != nil {
			return err
		}

		log.Println("wallet balance updated")

		return nil
	})
}
