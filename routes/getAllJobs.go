package routes

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "golang-project/models"
    "golang-project/utils"
)

func GetAllJobs(c *gin.Context) {
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

    // Fetch user data from the database
    _, err = utils.UserDataFromDB(db, email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user data"})
        return
    }

    // Retrieve all jobs from the database
    var jobs []models.Job
    if err := db.Find(&jobs).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve jobs"})
        return
    }

    // Return the list of jobs
    c.JSON(http.StatusOK, gin.H{"jobs": jobs})
}
