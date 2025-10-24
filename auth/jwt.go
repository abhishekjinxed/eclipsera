package auth

import (
	"lumora/user"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("supersecretkey123") // ⚠️ Replace with env variable in production!

type Claims struct {
	Email    string        `json:"email"`
	Name     string        `json:"name"`
	Friends  []user.Friend `json:"friends"`
	Picture  string        `json:"picture"`
	GoogleID string        `json:"google_id"`

	jwt.RegisteredClaims
}
type Friend struct {
	UserID string `json:"user_id"`
	Status string `json:"status"` // "pending", "accepted", "rejected"
}

// GenerateJWT creates a new token valid for 24 hours
func GenerateJWT(email, name string, friends []user.Friend, picture string, google_id string) (string, error) {
	claims := &Claims{
		Email:   email,
		Name:    name,
		Picture: picture,
		Friends: friends,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		GoogleID: google_id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateJWT verifies a token
func ValidateJWT(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
