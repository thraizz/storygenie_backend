package controller

import (
	"fmt"
	"net/http"
	"storygenie-backend/api"
	"storygenie-backend/models"

	"github.com/gin-gonic/gin"
)

// This endpoint handles feedback from the user for a given story
func (c *PublicController) CreateFeedback(context *gin.Context) {
	var input api.AddFeedbackForStoryJSONRequestBody
	if err := context.ShouldBindJSON(&input); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var product = models.Feedback{
		StoryID: input.StoryId,
		Liked:   input.Liked,
		UserID:  context.MustGet("user_id").(string),
	}
	c.Database.Create(&product)
	context.JSON(http.StatusOK, gin.H{"data": product})
}

func (c *PublicController) GetFeedbackForStory(context *gin.Context) {
	var feedback = models.Feedback{}
	productId := context.Param("storyId")
	result := c.Database.First(&feedback, "story_id = ? AND user_id = ?", productId, context.MustGet("user_id").(string))

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		if result.Error.Error() == "record not found" {
			context.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Story not found"})
			return
		}
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}
	if result.RowsAffected == 0 {
		context.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Story not found"})
		return
	}

	response := api.Feedback{
		Liked:   feedback.Liked,
		StoryId: feedback.StoryID,
	}

	context.JSON(http.StatusOK, response)
}
