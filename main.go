package main

import (
    "github.com/gin-gonic/gin"
    "golang-project/routes"
)

func main() {
    router := gin.Default()

    // Users routes
    router.POST("/signup", routes.Signup)    // create or return logic
    router.POST("/auth/login", routes.Login)    // auth/endpoint - jwt required
    // router.GET("auth/jobs/apply", routes.ApplyForJob)    // applicant


    // did this step as i wanted to use only One Database ,image uploading directly is paid
    // this uploads resume in DB table resume as binary BLOB
    router.POST("/auth/upload-resume", routes.UploadResume)   

    // this re-convert binary blob to image 
    // (for user - requires only valid auth token)
    router.GET("/auth/read-resume", routes.UserGetResume)
    // (for admin - we cant use user's auth token so pass email or user_id in body)
    router.GET("/admin/read-resume", routes.AdminGetResume)


    router.POST("/admin/job", routes.CreateJob)   // admin auth create job

    router.GET("/jobs", routes.GetAllJobs)   // auth user/admin  get all job






    // GET /admin/job/{job_id}: Authenticated API for fetching information regarding a job
    // opening. Returns details about the job opening and a list of applicants. Only Admin type
    // users can access this API.
    // router.GET("admin/job/:job_id", routes.GetJob) // admin 


    // GET /admin/applicants: Authenticated API for fetching a list of all users in the system. Only
    // Admin type users can access this API.
    // router.GET("/admin/applicants", routes.GetAllApplicants)  // admin auth 


    // GET /admin/applicant/{applicant_id}: Authenticated API for fetching extracted data of an
    // applicant. Only Admin type users can access this API.
    // router.GET("/admin/applicant/:applicant_id", routes.GetApplicantDetails)


//     GET /jobs/apply?job_id={job_id}: Authenticated API for applying to a particular job. Only
// Applicant users are allowed to apply for jobs.
    // router.GET("/applicant/:applicant_id", routes.GetApplicantDetails)


    router.Run(":8080")  
}
