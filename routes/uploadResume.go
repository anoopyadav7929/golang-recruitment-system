package routes

import (
	"bytes"
	"encoding/json"
	"io"
	"path/filepath"
	"strings"

	"golang-project/models"
	"golang-project/utils"
	constants "golang-project/utils"
	"net/http"

	"github.com/gin-gonic/gin"
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

	if user.UserType != constants.Applicant {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only applicant can use this API"})
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

	// Call the third-party API to extract resume data
	req, err := http.NewRequest("POST", constants.ApiUrl, bytes.NewBuffer(fileContent))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create API request"})
		return
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("apikey", constants.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send API request"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract resume data External API error"})
		return
	}

	// Parse the JSON response from the API
	var responseData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse external API response"})
		return
	}

	skills, _ := json.Marshal(responseData["skills"])
	education, _ := json.Marshal(responseData["education"])
	experience, _ := json.Marshal(responseData["experience"])

	profile := models.Profile{
		UserId:     user.Id,
		ResumeId:   resume.Id,
		Skills:     string(skills),
		Education:  string(education),
		Experience: string(experience),
		Name:       responseData["name"].(string),
		Email:      responseData["email"].(string),
		Phone:      responseData["phone"].(string),
	}

	var existingProfile models.Profile
	profileResult := db.First(&existingProfile, "user_id = ?", user.Id)
	if profileResult.Error != nil && profileResult.Error != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check profile"})
		return
	}

	// Update or create the profile entry
	if profileResult.RowsAffected > 0 {
		// Update the existing profile
		if err := db.Table("profiles").Where("user_id = ?", existingProfile.UserId).Updates(profile).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// Create a new profile
		if err := db.Create(&profile).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create profile"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Resume uploaded successfully"})
}
