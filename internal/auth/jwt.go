package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("change-this-secret")

type Claims struct {
	DeviceID string `json:"device_id"`
	UserID   int64  `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(deviceID string, userID int64) (string, error) {
	claims := Claims{
		DeviceID: deviceID,
		UserID:   userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
