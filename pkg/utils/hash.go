package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword encrypts a plain text password using bcrypt.
func HashPassword(password string) (string, error) {
	// bcrypt.DefaultCost is 10. This provides a good balance between security and performance.
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash compares a plain text password with a hashed password from the database.
// It returns true if they match.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
