package utils

import (
	"errors"
	"fmt"
	"golang-project/models"
	"time"

	"gorm.io/gorm"
)

// UserDataFromDB retrieves user data based on user ID or email from the database
func UserDataFromDB(db *gorm.DB, value interface{}) (models.User, error) {
	var user models.User

	// Convert value to the appropriate type
	var userID int64
	var userEmail string
	switch v := value.(type) {
	case int:
		userID = int64(v)
	case string:
		userEmail = v
	default:
		return user, errors.New("value must be of type int (userID) or string (email)")
	}

	var err error
	if userID > 0 {
		err = db.First(&user, "id = ?", userID).Error
	} else if userEmail != "" {
		err = db.First(&user, "email = ?", userEmail).Error
	} else {
		return user, errors.New("neither ID nor email was provided")
	}

	// Handle the case where the user is not found or other errors
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, fmt.Errorf("user not found")
		}
		return user, err
	}

	return user, nil
}

func GetCurrentTime() (time.Time, string) {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	currentTime := time.Now().In(loc)
	return currentTime, currentTime.Format("2006-01-02 15:04:05")
}
