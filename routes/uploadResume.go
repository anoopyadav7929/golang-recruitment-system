package routes

import (
	"io"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"golang-project/models"
	"golang-project/utils"
	"net/http"
	"gorm.io/gorm"
)

func UploadResume(c *gin.Context) {
	db := utils.InitializeDB()

	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization token is missing"})
		return
	}
	email, err := utils.ExtractEmailFromToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid authorization token"})
		return
	}
	user, err := utils.UserDataFromDB(db, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user data"})
		return
	}

	// Retrieve the file from the request
	file, fileHeader, err := c.Request.FormFile("resume")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}
	defer file.Close()

	fileName := fileHeader.Filename
	ext := strings.ToLower(filepath.Ext(fileName))
	
	var docType string
	switch ext {
	case ".pdf":
		docType = "pdf"
	case ".docx":
		docType = "docx"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document type only .pdf or .docx allowded"})
		return
	}

	// Read the file content into a byte slice
	fileContent, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	// Check if a resume already exists for the user
	var resume models.Resume
	result := db.First(&resume, "user_id = ?", user.Id)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check resume"})
		return
	}

	// Update or create the resume entry
	if result.RowsAffected > 0 {
		resume.DocContent = fileContent
		resume.DocType = docType
		if err := db.Save(&resume).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update resume"})
			return
		}
	} else {
		// Resume does not exist, create a new one
		resume = models.Resume{
			UserId:     user.Id,
			DocContent: fileContent,
			DocType:    docType,
		}
		if err := db.Create(&resume).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create resume"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Resume uploaded successfully"})
}
