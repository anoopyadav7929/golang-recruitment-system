package routes

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "golang-project/models"
    "golang-project/utils"
	constants   "golang-project/utils"
)

func GetAllApplicants(c *gin.Context) {
    // Initialize the database connection
    db := utils.InitializeDB()

    // Retrieve and validate the token
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

	var user models.User
    // Fetch user data from the database
    user, err = utils.UserDataFromDB(db, email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user data"})
        return
    }

	if user.UserType != constants.Admin {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Admin access is Required"})
        return	}


    // Retrieve all user from the database
    var applicant []models.User
    if err := db.Where("user_type = ?", constants.Applicant).Find(&applicant).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user data"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"Applicants": applicant})
}

