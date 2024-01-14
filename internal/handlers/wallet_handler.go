package handlers

import (
	"bl-wallet-service/internal/models"
	"bl-wallet-service/internal/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

type WalletHandler struct {
	walletService services.IWalletService
}

func NewWalletHandler(walletService services.IWalletService) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
	}
}

// @Summary Create a Wallet
// @Description Create a wallet for the specified user ID.
// @ID create-wallet
// @Accept json
// @Produce json
// @Param request body models.CreateWalletRequest true "Create Wallet Request"
// @Success 200 {object} models.CreateWalletResponse "OK"
// @Router /users/wallet [post]
func (h *WalletHandler) CreateWallet(c *gin.Context) {
	var request models.CreateWalletRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.walletService.Create(request.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &models.CreateWalletResponse{
		UserID:  request.UserID,
		Message: "wallet created",
	})
}

// @Summary Get Wallet
// @Description Get wallet details for the specified user ID.
// @ID get-wallet
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.Wallet "OK"
// @Router /users/{id}/wallet [get]
func (h *WalletHandler) GetWallet(c *gin.Context) {
	userID := c.Param("id")
	wallet, err := h.walletService.GetByUserID(userID)
	if wallet == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusInternalServerError,
			"message": fmt.Sprintf("Wallet of user %s not found", userID),
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, wallet)
}

// @Summary Add Funds to Wallet
// @Description Add funds to the wallet for the specified user ID.
// @ID add-funds
// @Accept json
// @Produce json
// @Param request body models.WalletFundsRequest true "Wallet Funds Request"
// @Param x-idempotency-key header string false "Idempotency Key"
// @Success 200 {object} models.WalletFundsResponse "OK"
// @Router /users/wallet/add [post]
func (h *WalletHandler) AddFunds(c *gin.Context) {
	var request models.TransactionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validateFundsRequest(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	err := h.walletService.ProcessTransaction(&request, models.CreditTransactionType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &models.WalletFundsResponse{
		UserID:  request.UserID,
		Message: "funds were added",
	})
}

// @Summary Remove Funds from Wallet
// @Description Remove funds from the wallet for the specified user ID.
// @ID remove-funds
// @Accept json
// @Produce json
// @Param request body models.WalletFundsRequest true "Wallet Funds Request"
// @Param x-idempotency-key header string false "Idempotency Key"
// @Success 200 {object} models.WalletFundsResponse "OK"
// @Router /users/wallet/remove [post]
func (h *WalletHandler) RemoveFunds(c *gin.Context) {
	var request models.TransactionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validateFundsRequest(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	err := h.walletService.ProcessTransaction(&request, models.DebitTransactionType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &models.WalletFundsResponse{
		UserID:  request.UserID,
		Message: "funds were removed",
	})
}

func validateFundsRequest(request models.TransactionRequest) error {
	if request.UserID == "" {
		return errors.New("You should provide a user id")
	}

	if request.TransactionID == "" {
		return errors.New("You should provide a transaction id")
	}

	if request.Amount < 0.0 {
		return errors.New("Amount should not be negative")
	}

	return nil
}
