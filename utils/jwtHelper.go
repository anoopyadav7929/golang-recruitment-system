package utils

import (
    "fmt"
    "time"
    jwt "github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("your_secret_key")

// GenerateToken creates a new JWT token with the provided email.
func GenerateToken(email string) (string, error) {
    // Create a new token with claims
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "email": email,
        "exp":   time.Now().Add(time.Hour * 124).Unix(),
    })
    return token.SignedString(jwtSecret)
}

// parseToken parses and validates the JWT token.
func parseToken(tokenStr string) (*jwt.Token, error) {
    // Parse the token
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        // Ensure the token method is what we expect
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return jwtSecret, nil
    })
    if err != nil {
        return nil, err
    }
    return token, nil
}

// two purpose , first validate, second extracts the email
func ExtractEmailFromToken(tokenStr string) (string, error) {
    token, err := parseToken(tokenStr)
    if err != nil {
        return "", err
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        email, ok := claims["email"].(string)
        if !ok {
            return "", fmt.Errorf("email claim not found")
        }
        return email, nil
    } else {
        return "", fmt.Errorf("invalid token")
    }
}
