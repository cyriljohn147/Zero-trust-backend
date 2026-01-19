package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/cyriljohn147/zero-trust-backend/internal/db"
)

func DeviceActiveOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceIDStr, exists := c.Get("device_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "device not found"})
			return
		}

		deviceID, err := uuid.Parse(deviceIDStr.(string))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid device"})
			return
		}

		device, err := db.GetDeviceByDeviceID(c.Request.Context(), deviceID)
		if err != nil || device.Status != "active" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "device revoked"})
			return
		}

		c.Next()
	}
}
