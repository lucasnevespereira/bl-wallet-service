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
	log.Printf("Wallet to insert: %v \n", wallet)
	err := s.db.Create(&wallet).Error
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Wallet created in database")
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
	err := s.db.First(&rowTransaction, "id = ?", transactionID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return rowTransaction, nil
}

func (s *WalletStorage) CreateTransaction(transaction *RowTransaction) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		log.Printf("inserting transaction: %v \n", transaction)
		err := tx.Model(&RowTransaction{}).Create(&transaction).Error
		if err != nil {
			return err
		}

		log.Printf("transaction inserted: %v \n", transaction)
		return nil
	})
}

func (s *WalletStorage) UpdateBalance(wallet *RowWallet, transactionID string) error {
	log.Println("updating balance")
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

		var isolatedWallet *RowWallet
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? and wallet_version = ?", wallet.ID, wallet.WalletVersion).First(&isolatedWallet).Error
		if err != nil {
			log.Printf("tx.Clauses locking err: %v \n", err)
			return err
		}

		// check transaction type
		if transaction.TransactionType == models.CreditTransactionType {
			log.Println("crediting funds")
			updateValues := map[string]interface{}{
				"balance":        isolatedWallet.Balance + transaction.Amount,
				"wallet_version": isolatedWallet.WalletVersion + 1,
			}
			err := tx.Model(&RowWallet{}).Where("user_id = ?", isolatedWallet.UserID).Updates(updateValues).Error
			if err != nil {
				updatedTransaction.Status = models.FailedTransactionStatus
				err := s.db.Model(&RowTransaction{}).Where("id = ?", transaction.ID).Updates(updatedTransaction).Error
				if err != nil {
					return err
				}
				return err
			}
			log.Println("funds credited")
		} else if transaction.TransactionType == models.DebitTransactionType {
			if wallet.Balance < transaction.Amount {
				log.Println("debiting funds")
				updatedTransaction.Status = models.FailedTransactionStatus
				err := s.db.Model(&RowTransaction{}).Where("id = ?", transaction.ID).Updates(updatedTransaction).Error
				if err != nil {
					return err
				}
				return errors.New(models.INSUFFICIENT_FUNDS_ERROR)
			}

			updateValues := map[string]interface{}{
				"balance":        isolatedWallet.Balance - transaction.Amount,
				"wallet_version": isolatedWallet.WalletVersion + 1,
			}
			err := tx.Model(&RowWallet{}).Where("user_id = ?", isolatedWallet.UserID).Updates(updateValues).Error
			if err != nil {
				updatedTransaction.Status = models.FailedTransactionStatus
				err := s.db.Model(&RowTransaction{}).Where("id = ?", transaction.ID).Updates(updatedTransaction).Error
				if err != nil {
					return err
				}
				return err
			}
			log.Println("funds debited")
		}

		updatedTransaction.Status = models.SuccessTransactionStatus
		err = s.db.Model(&RowTransaction{}).Where("id = ?", transaction.ID).Updates(&updatedTransaction).Error
		if err != nil {
			return err
		}

		log.Println("balance updated")
		return nil
	})
}
