package controller

import (
	"net/http"
	"storygenie-backend/models"

	"github.com/gin-gonic/gin"
)

// Get project by id
func (c *PublicController) GetProjectById(context *gin.Context) {
	var uid = context.MustGet("uid").(string)
	var project models.Project
	if err := c.Database.Where("id = ? AND uid = ?", context.Param("id"), uid).First(&project).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"data": project})
}

// Get all projects from this user
func (c *PublicController) GetProjects(context *gin.Context) {
	var project = []models.Project{}
	c.Database.Find(&project, "uid = ?", context.MustGet("uid").(string))
	context.JSON(http.StatusOK, gin.H{"data": project})
}

// Create a new project
func (c *PublicController) CreateProject(context *gin.Context) {

	var input models.Project
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.UserID = context.MustGet("uid").(string)
	c.Database.Create(&input)
	context.JSON(http.StatusOK, gin.H{"data": input})
}

// Delete a project by id
func (c *PublicController) DeleteProject(context *gin.Context) {
	var uid = context.MustGet("uid").(string)
	var project models.Project
	if err := c.Database.Where("id = ? AND uid = ?", context.Param("id"), uid).First(&project).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.Database.Delete(&project)
	context.JSON(http.StatusOK, gin.H{"data": true})
}
