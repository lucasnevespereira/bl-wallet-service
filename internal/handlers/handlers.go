package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary Health Check
// @Description Check the health status of the service.
// @ID health-check
// @Accept json
// @Produce json
// @Success 200
// @Router /health [get]
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "up"})
}

func NoRoute(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "error": "Not Found"})
}
