package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/cyriljohn147/zero-trust-backend/internal/db"
	"github.com/cyriljohn147/zero-trust-backend/internal/services"
)

type RegisterDeviceRequest struct {
	PublicKey string `json:"public_key" binding:"required"`
}

func RegisterDeviceHandler(c *gin.Context) {
	// TEMP: user_id (replace later with auth middleware)
	userID := int64(1)

	var req RegisterDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	device, err := services.RegisterDevice(
		c.Request.Context(),
		userID,
		req.PublicKey,
	)
	if err != nil {
		_ = db.CreateAuditLog(
			c.Request.Context(),
			uuid.New(),
			&userID,
			nil,
			"register",
			"failure",
			c.ClientIP(),
			c.Request.UserAgent(),
		)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "device registration failed"})
		return
	}

	_ = db.CreateAuditLog(
		c.Request.Context(),
		uuid.New(),
		&userID,
		&device.DeviceID,
		"register",
		"success",
		c.ClientIP(),
		c.Request.UserAgent(),
	)

	c.JSON(http.StatusCreated, gin.H{
		"device_id": device.DeviceID,
	})
}
