package api

import (
	"net/http"

	"github.com/cyriljohn147/zero-trust-backend/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RevokeDeviceRequest struct {
	DeviceID string `json:"device_id" binding:"required"`
}

func RevokeDeviceHandler(c *gin.Context) {
	var req RevokeDeviceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "invalid_request",
			"message": "Invalid request",
		})
		return
	}

	deviceID, err := uuid.Parse(req.DeviceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "invalid_device",
			"message": "Invalid device ID",
		})
		return
	}

	err = db.RevokeDevice(c.Request.Context(), deviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "failed",
			"message": "Failed to revoke device",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Device revoked successfully",
	})
}
