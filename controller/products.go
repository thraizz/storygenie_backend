package controller

import (
	"net/http"
	"storygenie-backend/api"
	"storygenie-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Get product by id
func (c *PublicController) GetProductById(context *gin.Context) {
	var user_id = context.MustGet("user_id").(string)

	productId := uuid.MustParse(context.Param("productId"))

	product := models.Product{
		UserID: user_id,
		UID:    productId,
	}

	if err := c.Database.Where("user_id = ?", user_id).Preload("Story").Where("uid = ?", productId).Find(&product).Error; err != nil {
		// If the record is not found, return an NotFound error.
		if err.Error() == "record not found" {
			context.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Internal server error"})
		return
	}
	stories := []api.Story{}
	for _, product := range product.Story {
		stories = append(stories, api.Story{
			CreatedAt:          product.CreatedAt,
			UpdatedAt:          product.UpdatedAt,
			Id:                 product.UID,
			AcceptanceCriteria: product.AcceptanceCriteria.Data,
			Headline:           product.Headline,
			UserStory:          product.UserStory,
			ProductId:          product.Product.UID,
		})
	}

	var response = api.ProductWithStories{
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
		Description: product.Description,
		IsExample:   product.IsExample,
		Id:          product.UID,
		Name:        product.Name,
		Stories:     &stories,
	}
	context.JSON(http.StatusOK, response)
}

// Get all products from this user
func (c *PublicController) GetProducts(context *gin.Context) {
	var products = []models.Product{}
	if err := c.Database.Find(&products, "user_id = ?", context.MustGet("user_id").(string)).Error; err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Internal server error"})
		return
	}
	response := []api.Product{}
	for _, product := range products {
		response = append(response, api.Product{
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
			DeletedAt:   &product.DeletedAt.Time,
			Description: product.Description,
			IsExample:   product.IsExample,
			Id:          product.UID,
			Name:        product.Name,
		})
	}

	context.JSON(http.StatusOK, response)
}

// Create a new product
func (c *PublicController) CreateProduct(context *gin.Context) {
	var input models.ProductInput
	if err := context.ShouldBindJSON(&input); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var product = models.Product{
		Name:        input.Name,
		Description: input.Description,
		UserID:      context.MustGet("user_id").(string),
	}
	c.Database.Create(&product)
	response := api.Product{
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
		DeletedAt:   &product.DeletedAt.Time,
		Description: product.Description,
		IsExample:   product.IsExample,
		Id:          product.UID,
		Name:        product.Name,
	}
	context.JSON(http.StatusOK, response)
}

// Delete a product by id
func (c *PublicController) DeleteProduct(context *gin.Context) {
	var uid = context.MustGet("user_id").(string)
	var product models.Product
	if err := c.Database.Where("uid = ? AND user_id = ?", context.Param("id"), uid).First(&product).Error; err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.Database.Delete(&product)
	context.JSON(http.StatusOK, gin.H{"data": true})
}
