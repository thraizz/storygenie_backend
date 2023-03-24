package controller

import (
	"storygenie-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PublicController struct {
	Database *gorm.DB
}

func (c *PublicController) SeedDatabase(context *gin.Context) {
	user_id := context.MustGet("user_id").(string)

	c.Database.Migrator().DropTable(&models.Product{})
	c.Database.Migrator().DropTable(&models.Story{})
	c.Database.Migrator().DropTable(&models.Feedback{})

	c.Database.AutoMigrate(&models.Product{})
	c.Database.AutoMigrate(&models.Story{})
	c.Database.AutoMigrate(&models.Feedback{})

	firstProduct := models.Product{
		Name:        "Github",
		Description: "Github is a website for hosting and collaborating on code",
		IsExample:   true,
		UserID:      user_id,
	}
	c.Database.Create(&firstProduct)

	firstStory := models.Story{
		Headline:           "Rearrange Alert Colors",
		UserStory:          "As a user, I want to see the right colors in the alerts as well as for the occurences of the secondary color.",
		AcceptanceCriteria: datatypes.JSONType[[]string]{[]string{"Every secondary color ocurrence should be checkend in Storygenie", "Alert colors align with the colors in the design"}},
		ProductID:          firstProduct.UID,
		UserID:             user_id,
	}
	c.Database.Create(&firstStory)
	secondStory := models.Story{
		Headline:           "Add code search across organization",
		UserStory:          "As a user, I want to search for code across my organization.",
		AcceptanceCriteria: datatypes.JSONType[[]string]{[]string{"Code should be searchable across organization", "Code results should be ordered by relevance"}},
		ProductID:          firstProduct.UID,
		UserID:             user_id,
	}
	c.Database.Create(&secondStory)

	secondProduct := models.Product{
		Name:        "Google",
		Description: "Google is a search engine",
		IsExample:   true,
		UserID:      user_id,
	}
	c.Database.Create(&secondProduct)
	thirdStory := models.Story{
		Headline:           "Add image search",
		UserStory:          "As a user, I want to search for images.",
		AcceptanceCriteria: datatypes.JSONType[[]string]{[]string{"Images should be searchable", "Images should be ordered by relevance"}},
		ProductID:          secondProduct.UID,
		UserID:             user_id,
	}
	c.Database.Create(&thirdStory)
	fourthStory := models.Story{
		Headline:           "Add video search",
		UserStory:          "As a user, I want to search for videos.",
		AcceptanceCriteria: datatypes.JSONType[[]string]{[]string{"Videos should be searchable", "Videos should be ordered by relevance"}},
		ProductID:          secondProduct.UID,
		UserID:             user_id,
	}
	c.Database.Create(&fourthStory)

}
