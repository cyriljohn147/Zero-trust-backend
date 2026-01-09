package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/cyriljohn147/zero-trust-backend/internal/services"
)

type ChallengeRequest struct {
	DeviceID string `json:"device_id" binding:"required"`
}

func GenerateChallengeHandler(c *gin.Context) {
	var req ChallengeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	deviceID, err := uuid.Parse(req.DeviceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device_id"})
		return
	}

	challenge, err := services.GenerateChallenge(
		c.Request.Context(),
		deviceID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate challenge"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"challenge_id": challenge.ChallengeID,
		"challenge":    challenge.Challenge,
		"expires_at":   challenge.ExpiresAt,
	})
}
