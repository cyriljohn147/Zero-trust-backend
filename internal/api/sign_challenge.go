package api

import (
	"net/http"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

type SignRequest struct {
	Challenge string `json:"challenge" binding:"required"`
}

func SignChallengeHandler(c *gin.Context) {
	var req SignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	cmd := exec.Command(
		"/Users/cyriljohn147/Documents/coding/zero-trust-backend/venv/bin/python",
		"/Users/cyriljohn147/Documents/coding/zero-trust-backend/sign_challenge.py",
		req.Challenge,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": string(out), // shows real python error if any
		})
		return
	}

	sig := strings.TrimSpace(string(out))
	if sig == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "signing failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"signature": sig,
	})
}
