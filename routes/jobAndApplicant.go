package routes

import (
	"golang-project/models"
	"golang-project/utils"
	constants "golang-project/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetJobAndApplicants(c *gin.Context) {
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
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin users can access this API"})
		return
	}

	// Retrieve job ID from URL parameter
	jobId := c.Param("job_id")
	if jobId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Job ID is required"})
		return
	}

	var job models.Job
	if err := db.First(&job, "id = ?", jobId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	// this will return user_id that applied to the jobs
	var applications []models.JobApplication
	if err := db.Where("job_id = ?", jobId).Find(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve applications"})
		return
	}

	// Get applicants details
	// will extract user data from user table and profile data from profile table
	var applicants []gin.H // map[string]any
	for _, application := range applications {
		var user models.User
		var profile models.Profile
		if err := db.First(&user, "id = ?", application.UserId).Error; err == nil {
			if err := db.First(&profile, "user_id = ?", application.UserId).Error; err == nil {
				applicants = append(applicants, gin.H{
					"user":    user,
					"profile": profile,
				})
			}
		}
	}

	data := gin.H{
		"job details":                           job,
		"applicant, resumeid and their profile": applicants,
	}

	c.JSON(http.StatusOK, data)
}
