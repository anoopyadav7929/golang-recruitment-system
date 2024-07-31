package routes

import (
	"golang-project/models"
	"golang-project/utils"
	constants "golang-project/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ApplyToJob(c *gin.Context) {
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
		c.JSON(http.StatusForbidden, gin.H{"error": "Only applicant users can apply for jobs"})
		return
	}

	// check user have uploded resume or not
	var resume models.Resume
	if err := db.First(&resume, "user_id = ?", user.Id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Please upload resume first before applying for job"})
		return
	}

	// Retrieve job ID from query parameter
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

	// Check if the applicant has already applied to the job
	var existingApplication models.JobApplication
	if err := db.First(&existingApplication, "job_id = ? AND user_id = ?", jobId, user.Id).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "You have already applied to this job"})
		return
	}

	// add user with the job in new DB, means he applies
	application := models.JobApplication{
		JobId:  job.Id,
		UserId: user.Id,
	}

	if err := db.Create(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot apply for job currently"})
		return
	}

	// Update the job's total applications count
	job.TotalApplications++
	if err := db.Save(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update job applications count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job application successful"})
}
