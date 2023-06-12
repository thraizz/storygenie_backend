package controller

import (
	"net/http"
	"os"
	"storygenie-backend/api"
	"storygenie-backend/models"

	"github.com/gin-gonic/gin"
)

// Return the jira client secret from the env variable
func (c *PublicController) GetJiraClientSecret(context *gin.Context) {
	context.JSON(200, os.Getenv("JIRA_CLIENT_SECRET"))
}

// Set the jira refresh token for this user
func (c *PublicController) SetJiraRefreshToken(context *gin.Context) {
	// Get the refresh token from the PUT body
	var input api.SetJiraRefreshTokenJSONBody
	if err := context.ShouldBindJSON(&input); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// fail if the refresh token is empty
	if input.RefreshToken == "" {
		context.AbortWithStatusJSON(400, gin.H{"error": "refreshToken is empty"})
		return
	}

	user_id := context.MustGet("user_id").(string)

	var jira = models.JiraRefreshToken{
		UserID:       user_id,
		RefreshToken: input.RefreshToken,
	}

	if err := c.Database.Where("user_id = ?", user_id).Assign(models.JiraRefreshToken{RefreshToken: input.RefreshToken}).FirstOrCreate(&jira).Error; err != nil {
		context.AbortWithStatusJSON(500, gin.H{"error": "Internal server error"})
		return
	}

	context.JSON(200, gin.H{"message": "success"})
}

// Get the jira refresh token for this user
func (c *PublicController) GetJiraRefreshToken(context *gin.Context) {
	user_id := context.MustGet("user_id").(string)

	var jira = models.JiraRefreshToken{
		UserID: user_id,
	}

	// If we find an refresh token, return it. Otherwise, return an empty string
	if err := c.Database.First(&jira).Where("user_id = ?", context.MustGet("user_id").(string)).Error; err != nil {
		context.JSON(200, gin.H{"refreshToken": ""})
		return
	}

	context.JSON(200, gin.H{"refreshToken": jira.RefreshToken})
}
