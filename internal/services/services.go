package services

import (
	"bl-wallet-service/configs"
	"bl-wallet-service/internal/storage"
	"log"
)

type Services struct {
	WalletService IWalletService
}

func InitServices(config configs.Config) *Services {
	walletStorage, err := storage.NewWalletStorage(storage.WalletStorageConfig{
		DbHost:     config.DbHost,
		DbPort:     config.DbPort,
		DbUser:     config.DbUser,
		DbPassword: config.DbPassword,
		DbName:     config.DbName,
		DbSsl:      config.DbSsl,
	})
	if err != nil {
		log.Printf("could not init walletStorage: %v \n", err)
	}

	walletSvc := NewWalletService(walletStorage)
	return &Services{
		WalletService: walletSvc,
	}
}
