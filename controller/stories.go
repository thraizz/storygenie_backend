package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"storygenie-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	gogpt "github.com/sashabaranov/go-gpt3"
)

func (c *PublicController) GetStories(context *gin.Context) {
	var story = models.Story{}
	result := c.Database.Find(&story, "uid = ?", context.MustGet("uid").(string))

	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			context.JSON(http.StatusNotFound, gin.H{"error": "Story not found"})
			return
		}
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"story": story})
}

func (c *PublicController) GetStoryById(context *gin.Context) {
	storyId := context.Param("id")

	var story = models.Story{}
	result := c.Database.First(&story, "id = ? AND uid = ?", storyId, context.MustGet("uid").(string))

	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			context.JSON(http.StatusNotFound, gin.H{"error": "Story not found"})
			return
		}
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"story": story})
}

func (c *PublicController) CreateStory(context *gin.Context) {
	var input models.Story
	input.UserID = context.MustGet("uiwd").(string)

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	story := c.Database.Create(&input)

	if story.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": story.Error})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"story": input})
}
func (c *PublicController) GenerateScrumStories(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := verifyIDToken(c, "asf"); err != nil {
			// Return a 401 Unauthorized if the user is not authenticated
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		client := gogpt.NewClient("your token")
		ctx := context.Background()
		var requestData struct {
			StoryDescription   string `json:"storyDescription" binding:"required"`
			ProjectId          string `json:"projectId" binding:"required"`
			ProjectDescription string `json:"projectDescription" binding:"required"`
		}

		if err := c.ShouldBindJSON(&requestData); err != nil {
			// Return a 400 Bad Request if the request body is invalid
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		prompt := fmt.Sprintf(`Generate a scrum story with a headline, userstory and acceptance criteria. It must be in this json format and parseable with JSON.parse: { "headline": string, "userStory": string, "acceptanceCriteria": string[] }\nHeadline as short as possible. Acceptance criterias as specific as possible. No acceptance criteria beyond the specified input. Acceptance criteria and user story can reference the project description. Project is:\n%s\nStory is:\n%s`, requestData.ProjectDescription, requestData.StoryDescription)

		completionParams := &gogpt.CompletionRequest{
			Prompt:      prompt,
			Model:       "text-davinci-003",
			MaxTokens:   1000,
			Temperature: 0,
		}

		result, err := client.CreateCompletion(ctx, *completionParams)
		if err != nil {
			// Return a 500 Internal Server Error if the AI API call fails
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		story := result.Choices[0].Text
		if story == "" {
			// Return a 500 Internal Server Error if the AI API call returns an empty response
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong while generating your story."})
			return
		}

		var parsedStory map[string]interface{}
		if err := json.Unmarshal([]byte(story), &parsedStory); err != nil {
			// Return a 500 Internal Server Error if the generated story is invalid JSON
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong while generating your story."})
			return
		}

		// Save the story to the database
		newStory := models.Story{
			Headline:           parsedStory["headline"].(string),
			UserStory:          parsedStory["userStory"].(string),
			AcceptanceCriteria: parsedStory["acceptanceCriteria"].(datatypes.JSON),
		}

		if err := db.Create(&newStory).Error; err != nil {
			// Return a 500 Internal Server Error if the database operation fails
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return the ID of the newly created story
		c.JSON(http.StatusOK, gin.H{"id": newStory.ID})
	}
}
