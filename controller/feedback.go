package controller

import (
	"fmt"
	"net/http"
	"storygenie-backend/api"
	"storygenie-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// This endpoint handles feedback from the user for a given story
func (c *PublicController) CreateFeedback(context *gin.Context) {
	var input api.AddFeedbackForStoryJSONRequestBody
	if err := context.ShouldBindJSON(&input); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var feedback = models.Feedback{
		StoryID: input.StoryId,
		Liked:   input.Liked,
		UserID:  context.MustGet("user_id").(string),
	}
	c.Database.Create(&feedback)
	context.JSON(http.StatusOK, feedback)
}

func (c *PublicController) GetFeedbackForStory(context *gin.Context) {
	var feedback = models.Feedback{}
	productId := uuid.MustParse(context.Param("storyId"))
	result := c.Database.First(&feedback, "story_id = ? AND user_id = ?", productId, context.MustGet("user_id").(string))

	var response = api.Feedback{}
	if result.Error != nil {
		fmt.Println(result.Error.Error())
		if result.Error.Error() == "record not found" {
			// No feedback so far is okay, we just return an empty object
			context.JSON(http.StatusOK, struct {
				StoryId uuid.UUID `json:"storyId"`
			}{
				StoryId: productId,
			})
			return
		}
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	response = api.Feedback{
		Liked:   feedback.Liked,
		StoryId: feedback.StoryID,
	}

	context.JSON(http.StatusOK, response)
}
