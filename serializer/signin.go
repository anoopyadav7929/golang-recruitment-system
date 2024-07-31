package serializer

import (
	"errors"
	"golang-project/models"
	constants "golang-project/utils"
	"regexp"
	"strings"
)

// to do : add filed check - done
func ValidateUser(user *models.User) error {

	// name , profileheadline , usertype
	if strings.TrimSpace(user.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(user.ProfileHeadline) == "" {
		return errors.New("profileHeadline is required")
	}
	if user.UserType != constants.Admin {
		user.UserType = constants.Applicant
	}

	// Check if Email is not empty and has a valid format
	if strings.TrimSpace(user.Email) == "" {
		return errors.New("mail id is required")
	}
	if !isValidEmail(user.Email) {
		return errors.New("invalid email format")
	}

	// password checkadded num and char
	if !passwordValid(user.Password) {
		return errors.New("password is required (len must be greater than 6, have one num and one character)")
	}

	return nil
}

// isValidEmail validates the email format
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// pass check
// len > 6 , one char and one num atleast
func passwordValid(password string) bool {
	password = strings.TrimSpace(password)
	if len(password) <= 6 {
		return false
	}

	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	if !hasLetter {
		return false
	}

	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	if !hasDigit {
		return false
	}

	return true
}