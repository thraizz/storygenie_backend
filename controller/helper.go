package controller

import (
	"context"
	"log"
	"storygenie-backend/helper"
	"storygenie-backend/models"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PublicController struct {
	Database *gorm.DB
}

func verifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	// Verify the ID token first.
	client, err := helper.GetFirebaseApp().Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	return client.VerifyIDToken(ctx, idToken)
}

func (c *PublicController) SeedDatabase(context *gin.Context) {
	c.Database.AutoMigrate(&models.User{})
	c.Database.AutoMigrate(&models.Project{})
	c.Database.AutoMigrate(&models.Prompt{})
	c.Database.AutoMigrate(&models.Story{})
}
