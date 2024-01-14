package router

import (
	"bl-wallet-service/internal/handlers"
	"bl-wallet-service/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine, services *services.Services) {
	router.Use(cors.Default())

	router.GET("/health", handlers.Health)
	router.NoRoute(handlers.NoRoute)

	walletHandler := handlers.NewWalletHandler(services.WalletService)
	users := router.Group("/users")
	users.POST("/wallet", walletHandler.CreateWallet)
	users.GET("/:id/wallet", walletHandler.GetWallet)
	users.POST("/wallet/add", walletHandler.AddFunds)
	users.POST("/wallet/remove", walletHandler.RemoveFunds)
}
