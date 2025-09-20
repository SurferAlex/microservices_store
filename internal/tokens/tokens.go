package tokens

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID   int    `json"user_id"`
	Username string `json"username"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte("your-super-secret-key-change-in-production")

func GenerateJWT(username string) (string, error) {
	claims := Claims {
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			

		}
	}
}