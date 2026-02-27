package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT creates a new JSON Web Token for authenticated users.
// It uses dynamic expiration time based on the configuration.
func GenerateJWT(userID string, role string, secret string, expireStr string) (string, error) {
	// 1. Ubah teks "24h" dari .env menjadi format waktu Golang
	expireDuration, err := time.ParseDuration(expireStr)
	if err != nil {
		return "", fmt.Errorf("invalid expire duration format: %w", err)
	}

	// 2. Set custom claims dengan waktu kedaluwarsa yang dinamis
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(expireDuration).Unix(), // Menggunakan hasil parse
		"iat":     time.Now().Unix(),
	}

	// 3. Buat dan sign tokennya
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
