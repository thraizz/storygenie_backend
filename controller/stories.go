package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"storygenie-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"

	gogpt "github.com/sashabaranov/go-gpt3"
)

func (c *PublicController) GetStories(context *gin.Context) {
	var story = []models.Story{}
	result := c.Database.Find(&story, "user_id = ?", context.MustGet("uid").(string))
	if result.Error != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": story})
}

func (c *PublicController) GetStoryById(context *gin.Context) {
	storyId := context.Param("id")

	var story = models.Story{}
	result := c.Database.First(&story, "id = ? AND user_id = ?", storyId, context.MustGet("uid").(string))

	if result.Error != nil {
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

	context.JSON(http.StatusOK, gin.H{"data": story})
}

func (c *PublicController) GenerateScrumStories(context *gin.Context) {
	client := gogpt.NewClient(os.Getenv("OPENAI_API_KEY"))
	var requestData struct {
		StoryDescription string `json:"description" binding:"required"`
		ProjectId        string `json:"projectId" binding:"required"`
	}

	if err := context.ShouldBindJSON(&requestData); err != nil {
		// Return a 400 Bad Request if the request body is invalid
		log.Println(err.Error())
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Get project from request body projectId
	var project models.Project
	if err := c.Database.Select("Description").Where("id = ? AND user_id = ?", requestData.ProjectId, context.MustGet("uid").(string)).First(&project).Error; err != nil {
		log.Println(err.Error())
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	prompt := fmt.Sprintf(`Generate a scrum story with a headline, userstory and acceptance criteria. It must be in this json format and parseable with JSON.parse: { "headline": string, "userStory": string, "acceptanceCriteria": string[] }\nHeadline as short as possible. Acceptance criterias as specific as possible. No acceptance criteria beyond the specified input. Acceptance criteria and user story can reference the project description. Project is:\n%s\nStory is:\n%s`, project.Description, requestData.StoryDescription)

	completionParams := &gogpt.CompletionRequest{
		Prompt:      prompt,
		Model:       "text-davinci-003",
		MaxTokens:   1000,
		Temperature: 0,
	}

	result, err := client.CreateCompletion(context, *completionParams)
	if err != nil {
		log.Println(err.Error())
		// Return a 500 Internal Server Error if the AI API call fails
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	story := result.Choices[0].Text
	if story == "" {
		log.Println(err.Error())
		// Return a 500 Internal Server Error if the AI API call returns an empty response
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong while generating your story."})
		return
	}

	var parsedStory map[string]interface{}
	if err := json.Unmarshal([]byte(story), &parsedStory); err != nil {
		log.Println(err.Error())
		// Return a 500 Internal Server Error if the generated story is invalid JSON
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong while generating your story."})
		return
	}

	// Save the story to the database
	newStory := models.Story{
		Headline:           parsedStory["headline"].(string),
		UserStory:          parsedStory["userStory"].(string),
		AcceptanceCriteria: parsedStory["acceptanceCriteria"].(datatypes.JSON),
	}

	if err := c.Database.Create(&newStory).Error; err != nil {
		log.Println(err.Error())
		// Return a 500 Internal Server Error if the database operation fails
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the ID of the newly created story
	context.JSON(http.StatusOK, gin.H{"data": newStory.ID})
}
