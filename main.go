package main

import (
	"golang-project/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Users routes
	router.POST("/signup", routes.Signup)                 // Create or return logic
	router.POST("/login", routes.Login)                   // JWT create
	router.GET("/jobs", routes.GetAllJobs)                // Get all jobs
	router.POST("/jobs/apply/:job_id", routes.ApplyToJob) // Apply to specific job

	// Resume upload and retrieval
	router.POST("/upload-resume", routes.UploadResume) // Upload resume in DB table as binary BLOB

	// these are extra api , for viewing full resume
	router.GET("/read-resume", routes.UserGetResume)        // Convert binary BLOB to image for applicant user, can view their resume
	router.GET("/admin/read-resume", routes.AdminGetResume) // For admin, pass email or user_id in body

	// Admin routes
	router.POST("/admin/job", routes.CreateJob)                              // Admin - Create job
	router.GET("/admin/applicants", routes.GetAllApplicants)                 // Admin - Get all applicants from user DB
	router.GET("/admin/applicant/:applicant_id", routes.GetApplicantDetails) // Admin - Fetch single applicant data

	// Job and applicants
	router.POST("/admin/job/:job_id", routes.GetJobAndApplicants) // Fetch job and applicants details

	router.Run(":8080")
}
