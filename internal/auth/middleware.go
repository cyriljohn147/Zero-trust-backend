package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/cyriljohn147/zero-trust-backend/internal/db"
)

func ZeroTrustMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			return
		}

		tokenStr := parts[1]

		token, err := jwt.ParseWithClaims(
			tokenStr,
			&Claims{},
			func(token *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			},
		)

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}

		// 🔥 ADD THIS BLOCK (REVOCATION CHECK)
		deviceID, err := uuid.Parse(claims.DeviceID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid device ID"})
			return
		}
		device, err := db.GetDeviceByDeviceID(c.Request.Context(), deviceID)
		if err != nil || device.Status != "active" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "device revoked or inactive",
			})
			return
		}

		// Attach identity to request context
		c.Set("device_id", claims.DeviceID)
		c.Set("user_id", claims.UserID)

		c.Next()
	}
}
