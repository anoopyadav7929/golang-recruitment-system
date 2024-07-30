package main

import (
    "github.com/gin-gonic/gin"
    "golang-project/routes"
    // "golang-project/server"
    // "golang-project/utils"
)

func main() {
    router := gin.Default()

    // Users routes
    router.POST("/signup", routes.Signup)    // create or return logic
    router.POST("auth/login", routes.Login)    // auth/endpoint - jwt required
    router.GET("/jobs", routes.GetAllJobs)   
    // router.GET("auth/jobs/apply", routes.ApplyForJob)    // applicant

    // did this step as i wanted to use only One Database ,image uploading directly is paid
    // this uploads resume in DB table resume as binary BLOB
    router.POST("auth/upload-resume", routes.UploadResume)   

    // this re-convert binary blob to image 
    // (for user - requires only valid auth token)
    router.GET("auth/read-resume", routes.UserGetResume)
    // (for admin - we cant use user's auth token so pass email or user_id in body)
    router.GET("admin/read-resume", routes.AdminGetResume)



    // Admin specific routes
    // router.GET("/admin/applicants", routes.GetAllApplicants)
    // router.GET("/admin/applicant/:applicant_id", routes.GetApplicantDetails)

    router.POST("/admin/job", routes.CreateJob)
    router.GET("/admin/job/:job_id", routes.GetJob)

    router.Run(":8080")  
}
