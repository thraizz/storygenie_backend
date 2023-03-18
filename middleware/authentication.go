package middleware

import (
	"net/http"
	"storygenie-backend/helper"

	"github.com/gin-gonic/gin"
)

func Authentication(c *gin.Context) {
	uid, err := helper.GetUserFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	c.Set("uid", uid)
}
