package routes

import (
	"golang-project/models"
	"golang-project/serializer"
	db "golang-project/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Signup handles user registration
func Signup(c *gin.Context) {
	database := db.InitializeDB()

	var user models.User

	// Bind the JSON payload to the user struct
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := serializer.ValidateUser(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the user already exists by email
	var existingUser models.User
	result := database.Where("email = ?", user.Email).First(&existingUser)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	user.PasswordHash = string(hashedPassword)

	sql := `INSERT INTO users (name, email, address, user_type, password_hash, profile_headline)
			VALUES (?, ?, ?, ?, ?, ?)`
	result = database.Exec(sql, user.Name, user.Email, user.Address, user.UserType, user.PasswordHash, user.ProfileHeadline)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created"})
}
