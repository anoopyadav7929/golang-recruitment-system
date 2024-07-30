package serializer

import (
	"errors"
	"golang-project/models"
	"strings"
)

// to do : add filed check - done
func emailAndPasswordProvided(user *models.User) error {

	// name , profileheadline , usertype

	// Check if Email is not empty and has a valid format
	if strings.TrimSpace(user.Email) == "" {
		return errors.New("mail id is required")
	}

	// both are already in signin.go
	if !isValidEmail(user.Email) {
		return errors.New("invalid email format")
	}

	if !passwordValid(user.Password) {
		return errors.New("password is required (len must be greater than 6, have one num and one character)")
	}

	return nil
}
