package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"os"
	"time"
)

var refreshSecret = []byte(os.Getenv("REFRESH_SECRET"))

// Генерация случайного opauqe токена
func GenerateRefreshOpaque(userID int) (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// Хеширование для хранения в БД
func HashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return base64.RawStdEncoding.EncodeToString(sum[:])
}

func RefreshTTL() time.Duration { return 7 * 24 * time.Hour }
