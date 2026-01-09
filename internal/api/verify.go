package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/cyriljohn147/zero-trust-backend/internal/services"
)

type VerifyRequest struct {
	ChallengeID string `json:"challenge_id" binding:"required"`
	Signature   string `json:"signature" binding:"required"`
}

func VerifyChallengeHandler(c *gin.Context) {
	var req VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	challengeID, err := uuid.Parse(req.ChallengeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid challenge_id"})
		return
	}

	token, err := services.VerifyChallenge(
		c.Request.Context(),
		challengeID,
		req.Signature,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "verification failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
	})
}
