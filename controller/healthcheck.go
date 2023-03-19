package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Simple healthcheck endpoint
func (c *PublicController) HealthCheck(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"status": "ok"})
}
