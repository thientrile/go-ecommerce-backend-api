package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func GetHash(key string) string {

	hash := sha256.New()
	hash.Write([]byte(key))
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func GenerateSalt(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be greater than 0")
	}
	salt := make([]byte, length)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	return hex.EncodeToString(salt), nil
}

func HashPassword(password, salt string) string {
	if password == "" || salt == "" {
		return ""
	}
	hashed := sha256.Sum256([]byte(password + salt))
	return hex.EncodeToString(hashed[:])
}
