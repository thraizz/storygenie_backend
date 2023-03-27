// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package api

import (
	"time"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"gorm.io/datatypes"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Error defines model for Error.
type Error struct {
	Code    int32   `json:"code"`
	Hint    *string `json:"hint,omitempty"`
	Message string  `json:"message"`
}

// Feedback defines model for Feedback.
type Feedback struct {
	// Liked Whether the user liked the story or not.
	Liked bool `json:"liked"`

	// StoryId The story id
	StoryId openapi_types.UUID `json:"storyId"`
}

// Product defines model for Product.
type Product struct {
	// CreatedAt The date the product was created
	CreatedAt time.Time `json:"createdAt"`

	// DeletedAt The date the product was deleted
	DeletedAt *time.Time `json:"deletedAt"`

	// Description The description of the product
	Description string `json:"description"`

	// Id The id of the product
	Id openapi_types.UUID `json:"id"`

	// IsExample Whether the product is an example product
	IsExample bool `json:"isExample"`

	// Name The name of the product
	Name string `json:"name"`

	// UpdatedAt The date the product was updated
	UpdatedAt time.Time `json:"updatedAt"`
}

// ProductWithStories defines model for ProductWithStories.
type ProductWithStories struct {
	// CreatedAt The date the product was created
	CreatedAt time.Time `json:"createdAt"`

	// DeletedAt The date the product was deleted
	DeletedAt *time.Time `json:"deletedAt"`

	// Description The description of the product
	Description string `json:"description"`

	// Id The id of the product
	Id openapi_types.UUID `json:"id"`

	// IsExample Whether the product is an example product
	IsExample bool `json:"isExample"`

	// Name The name of the product
	Name    string   `json:"name"`
	Stories *[]Story `json:"stories,omitempty"`

	// UpdatedAt The date the product was updated
	UpdatedAt time.Time `json:"updatedAt"`
}

// Story defines model for Story.
type Story struct {
	// AcceptanceCriteria The acceptance criteria of the story
	AcceptanceCriteria datatypes.JSON `json:"acceptanceCriteria"`

	// CreatedAt The date the product was created
	CreatedAt time.Time `json:"createdAt"`

	// DeletedAt The date the product was deleted
	DeletedAt *time.Time `json:"deletedAt"`

	// Headline The headline of the story
	Headline string `json:"headline"`

	// Id The id of the story
	Id      openapi_types.UUID `json:"id"`
	Product Product            `json:"product"`

	// ProductId The id of the product
	ProductId openapi_types.UUID `json:"productId"`

	// UpdatedAt The date the product was updated
	UpdatedAt time.Time `json:"updatedAt"`

	// UserStory The user story of the story
	UserStory string `json:"userStory"`
}

// StoryInput defines model for StoryInput.
type StoryInput struct {
	// Description The description is a general idea of the story.
	Description string `json:"description"`

	// ProductId The product id
	ProductId openapi_types.UUID `json:"productId"`
}

// BadRequest defines model for BadRequest.
type BadRequest = Error

// NotFound defines model for NotFound.
type NotFound = Error

// GetAllProductsParams defines parameters for GetAllProducts.
type GetAllProductsParams struct {
	// Limit The numbers of items to return
	Limit *int `form:"limit,omitempty" json:"limit,omitempty"`
}

// GetAllStoriesParams defines parameters for GetAllStories.
type GetAllStoriesParams struct {
	// Limit The numbers of items to return
	Limit *int `form:"limit,omitempty" json:"limit,omitempty"`
}

// CreateProductJSONRequestBody defines body for CreateProduct for application/json ContentType.
type CreateProductJSONRequestBody = Product

// CreateStoryJSONRequestBody defines body for CreateStory for application/json ContentType.
type CreateStoryJSONRequestBody = Story

// GenerateStoryJSONRequestBody defines body for GenerateStory for application/json ContentType.
type GenerateStoryJSONRequestBody = StoryInput

// AddFeedbackForStoryJSONRequestBody defines body for AddFeedbackForStory for application/json ContentType.
type AddFeedbackForStoryJSONRequestBody = Feedback
