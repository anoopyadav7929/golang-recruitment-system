package main

import (
    "github.com/gin-gonic/gin"
    "golang-project/routes"
    // "golang-project/server"
    // "golang-project/utils"
)

func main() {
    // local so cant pass anything in db 
    router := gin.Default()

    // Users routes
    router.POST("/signup", routes.Signup)    // create or return logic
    router.POST("/login", routes.Login)      // applicant  

    router.POST("/upload-resume", routes.UploadResume)   

    router.GET("/jobs", routes.GetAllJobs)   

    // router.GET("/jobs/apply", routes.ApplyForJob)    // applicant

    // Admin specific routes
    // router.GET("/admin/applicants", routes.GetAllApplicants)
    // router.GET("/admin/applicant/:applicant_id", routes.GetApplicantDetails)

    router.POST("/admin/job", routes.CreateJob)
    router.GET("/admin/job/:job_id", routes.GetJob)

    router.Run(":8080")  
}
