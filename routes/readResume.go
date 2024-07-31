package routes

import (
	"golang-project/models"
	"golang-project/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserGetResume(c *gin.Context) {
    db := utils.InitializeDB()

    // Retrieve the token from the Authorization header
    token := c.GetHeader("Authorization")
    if token == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization token is missing"})
        return
    }

    // Extract email from the token
    email, err := utils.ExtractEmailFromToken(token)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid authorization token"})
        return
    }

    // Retrieve user data from the database
    user, err := utils.UserDataFromDB(db, email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user data"})
        return
    }

    // Find the resume for the user
    var resume models.Resume
    result := db.First(&resume, "user_id = ?", user.Id)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Resume not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve resume"})
        return
    }

    // Set the appropriate content type and file extension based on docType
    var contentType string
    switch resume.DocType {
    case "pdf":
        contentType = "application/pdf"
    case "docx":
        contentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
    default:
        contentType = "application/octet-stream"
    }

    // Return the resume content
    c.Header("Content-Type", contentType)
    c.Header("Content-Disposition", "attachment; filename=resumefile."+resume.DocType)
    c.Data(http.StatusOK, contentType, resume.DocContent)
}




//////////////////////////////////////////////
// for admin , pass user_id or email of applicant to view his resume 

type AdminGetResumeRequest struct {
    UserID int64 `json:"user_id"`
    Email  string `json:"email"`
}

func AdminGetResume(c *gin.Context) {
    db := utils.InitializeDB()

    var request AdminGetResumeRequest
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    if request.UserID <0 && request.Email == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User ID or email must be provided"})
        return
    }

    var user models.User
    var result *gorm.DB

    if request.UserID >0 {
        // Find user by ID
        result = db.First(&user, request.UserID)
    } else if request.Email != "" {
        // Find user by email
        result = db.Where("email = ?", request.Email).First(&user)
    }

    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
        return
    }

    // Find the resume for the user
    var resume models.Resume
    result = db.First(&resume, "user_id = ?", user.Id)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Resume not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve resume"})
        return
    }

    // Set the appropriate Content-Type based on document type
    var contentType string
    switch resume.DocType {
    case "pdf":
        contentType = "application/pdf"
    case "docx":
        contentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
    default:
        contentType = "application/octet-stream"
    }

    // Return the resume content
    c.Header("Content-Type", contentType)
    c.Header("Content-Disposition", "inline; filename=resumefile."+resume.DocType)
    c.Data(http.StatusOK, contentType, resume.DocContent)
}