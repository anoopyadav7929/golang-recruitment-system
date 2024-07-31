package routes

import (
	"golang-project/models"
	"golang-project/utils"
	constants "golang-project/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"golang-project/serializer"
)

func CreateJob(c *gin.Context) {
	db := utils.InitializeDB()

	// Retrieve the token from the Authorization header
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization token is missing"})
		return
	}

	email, err := utils.ExtractEmailFromToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	// Use the utility function to check the user type
	user, err := utils.UserDataFromDB(db, email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Check if the user is an admin
	if user.UserType != constants.Admin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin users can create jobs"})
		return
	}

	// default value for creating a job
	var job models.Job
	
	var currentTime time.Time
	currentTime = time.Now()
	job.PostedOn = &currentTime

	job.PostedBy = user.Id
	job.TotalApplications = 0

	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := serializer.ValidateJob(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job created", "job": job})
}
