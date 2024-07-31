package routes

import (
	"golang-project/models"
	"golang-project/utils"
	constants "golang-project/utils"

	"github.com/gin-gonic/gin"

	"net/http"
)

func GetApplicantDetails(c *gin.Context) {
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

	if user.UserType != constants.Admin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin users can access this endpoint"})
		return
	}

	// Retrieve applicant ID from URL parameter
	applicantID := c.Param("applicant_id")
	if applicantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Applicant ID is required"})
		return
	}

	var resume models.Resume
	var profile models.Profile

	if err := db.First(&resume, "user_id = ?", applicantID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resume not found"})
		return
	}

	if err := db.First(&profile, "user_id = ?", applicantID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

    // beautify json
	var skills []string
	if !utils.UnmarshalJSONField(profile.Skills, &skills, c) {
		return
	}
	var education []map[string]interface{}
	if !utils.UnmarshalJSONField(profile.Education, &education, c) {
		return
	}
	var experience []map[string]interface{}
	if !utils.UnmarshalJSONField(profile.Experience, &experience, c) {
		return
	}

	data := gin.H{
		"doc_type": resume.DocType,
		"profile": gin.H{
			"user_id":    profile.UserId,
			"resume_id":  profile.ResumeId,
			"skills":     skills,
			"education":  education,
			"experience": experience,
			"name":       profile.Name,
			"email":      profile.Email,
			"phone":      profile.Phone,
		},
	}

	c.JSON(http.StatusOK, data)
}
