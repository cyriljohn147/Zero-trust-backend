package api

import (
	"net/http"

	"errors"

	"github.com/cyriljohn147/zero-trust-backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type VerifyRequest struct {
	ChallengeID string `json:"challenge_id" binding:"required"`
	Signature   string `json:"signature" binding:"required"`
}

func VerifyChallengeHandler(c *gin.Context) {
	var req VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "invalid_request",
			"message": "Invalid request",
		})
		return
	}

	challengeID, err := uuid.Parse(req.ChallengeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "invalid_challenge",
			"message": "Invalid challenge ID",
		})
		return
	}

	token, err := services.VerifyChallenge(
		c.Request.Context(),
		challengeID,
		req.Signature,
		c.ClientIP(),
	)

	if err != nil {

		if errors.Is(err, services.ErrDeviceRevoked) {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "revoked",
				"message": "Device access has been revoked",
			})
			return
		}

		if errors.Is(err, services.ErrHighRisk) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "high_risk",
				"message": "Access denied due to high risk",
			})
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "failed",
			"message": "Verification failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"access_token": token,
	})
}
