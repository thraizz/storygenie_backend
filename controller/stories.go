package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"storygenie-backend/api"
	"storygenie-backend/models"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	gogpt "github.com/sashabaranov/go-gpt3"
)

func (c *PublicController) GetStories(context *gin.Context) {
	user_id := context.MustGet("user_id").(string)
	var stories = []models.Story{}
	result := c.Database.Model(&models.Story{}).Joins("Product").Find(&stories, "stories.user_id = ?", user_id)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	var response = []api.Story{}
	for _, story := range stories {
		product := api.Product{
			Name:        story.Product.Name,
			Id:          story.Product.UID,
			Description: story.Product.Description,
			CreatedAt:   story.Product.CreatedAt,
			DeletedAt:   nil,
			IsExample:   story.Product.IsExample,
			UpdatedAt:   story.Product.UpdatedAt,
		}
		response = append(response, api.Story{
			CreatedAt:          story.CreatedAt,
			UpdatedAt:          story.UpdatedAt,
			DeletedAt:          nil,
			Id:                 story.UID,
			Headline:           story.Headline,
			UserStory:          story.UserStory,
			AcceptanceCriteria: story.AcceptanceCriteria,
			Product:            product,
			ProductId:          story.ProductID,
		})
	}

	context.JSON(http.StatusOK, response)
}

func (c *PublicController) GetStoryById(context *gin.Context) {
	var user_id = context.MustGet("user_id").(string)
	storyId := uuid.MustParse(context.Param("storyId"))
	var story = models.Story{
		UserID: user_id,
		UID:    storyId,
	}
	result := c.Database.First(&story, "uid = ? AND user_id = ?", storyId, user_id)

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

	response := api.Story{
		CreatedAt:          story.CreatedAt,
		UpdatedAt:          story.UpdatedAt,
		DeletedAt:          &story.DeletedAt.Time,
		Id:                 story.UID,
		Headline:           story.Headline,
		UserStory:          story.UserStory,
		AcceptanceCriteria: story.AcceptanceCriteria,
	}

	context.JSON(http.StatusOK, response)
}

func (c *PublicController) GenerateScrumStories(context *gin.Context) {
	user_id := context.MustGet("user_id").(string)
	var requestData api.GenerateStoryJSONRequestBody

	if err := context.ShouldBindJSON(&requestData); err != nil {
		// Return a 400 Bad Request if the request body is invalid
		log.Println(err.Error())
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	productId := uuid.Must(requestData.ProductId, nil)

	// Get product from request body productId
	var product models.Product
	if err := c.Database.Select("Description").Where("uid = ? AND user_id = ?", requestData.ProductId, context.MustGet("user_id").(string)).First(&product).Error; err != nil {
		log.Println(err.Error())
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	prompt := fmt.Sprintf(`Generate a scrum story with a headline, userstory and acceptance criteria. It must be in valid JSON format, like this: { "headline": string, "userStory": string, "acceptanceCriteria": string[] }\nHeadline as short as possible. Acceptance criterias as specific as possible. No acceptance criteria beyond the specified input. Acceptance criteria and user story can reference the product description. Product is:%s, Story is:%s`, product.Description, requestData.Description)
	client := gogpt.NewClient(os.Getenv("OPENAI_API_KEY"))

	completionParams := &gogpt.CompletionRequest{
		Prompt:      prompt,
		Model:       "text-davinci-003",
		MaxTokens:   2000,
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
	storyWithoutNewLines := strings.ReplaceAll(story, "\n", "")
	var parsedStory api.Story

	if err := json.Unmarshal([]byte(storyWithoutNewLines), &parsedStory); err != nil {
		log.Println(err.Error())
		// Return a 500 Internal Server Error if the generated story is invalid JSON
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong while generating your story."})
		return
	}

	// Save the story to the database
	newStory := models.Story{
		Headline:           parsedStory.Headline,
		UserStory:          parsedStory.UserStory,
		AcceptanceCriteria: parsedStory.AcceptanceCriteria,
		ProductID:          productId,
		UserID:             user_id,
		UID:                uuid.New(),
	}

	if err := c.Database.Create(&newStory).Error; err != nil {
		log.Println(err.Error())
		// Return a 500 Internal Server Error if the database operation fails
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := api.Story{
		CreatedAt:          newStory.CreatedAt,
		UpdatedAt:          newStory.UpdatedAt,
		DeletedAt:          &newStory.DeletedAt.Time,
		Id:                 newStory.UID,
		Headline:           newStory.Headline,
		UserStory:          newStory.UserStory,
		AcceptanceCriteria: newStory.AcceptanceCriteria,
		ProductId:          productId,
	}
	// Return the newly created story
	context.JSON(http.StatusOK, response)
}

func (c *PublicController) CreateStory(context *gin.Context) {
	var input models.Story
	if err := context.ShouldBindJSON(&input); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.UserID = context.MustGet("user_id").(string)

	result := c.Database.Create(&input)
	if result.Error != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	context.JSON(http.StatusOK, input)
}

func (c *PublicController) DeleteStory(context *gin.Context) {
	storyId := context.Param("storyId")
	var story models.Story
	result := c.Database.First(&story, "uid = ? AND user_id = ?", storyId, context.MustGet("user_id").(string))

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

	result = c.Database.Delete(&story)
	if result.Error != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"data": "Story deleted"})
}
