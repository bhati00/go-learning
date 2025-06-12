package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
func CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
func GenerateJWT(username string) (string, error) {
	// Placeholder for JWT generation logic
	// In a real application, you would use a library like "github.com/dgrijalva/jwt-go"
	// to create a JWT token with the username as a claim.
	return "mocked.jwt.token", nil
}
