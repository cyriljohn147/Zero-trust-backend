package main

import (
	"log"

	"github.com/cyriljohn147/zero-trust-backend/internal/api"
	"github.com/cyriljohn147/zero-trust-backend/internal/auth"
	"github.com/cyriljohn147/zero-trust-backend/internal/db"
	"github.com/gin-gonic/gin"
)

func main() {
	db.Connect()
	defer db.Close()

	r := gin.Default()

	protected := r.Group("/api")
	protected.Use(
		auth.ZeroTrustMiddleware(),
		auth.DeviceActiveOnly(),
	)

	protected.GET("/secure-data", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Zero Trust access granted",
		})
	})
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.POST("/devices/register", api.RegisterDeviceHandler)
	r.POST("/auth/challenge", api.GenerateChallengeHandler)
	r.POST("/auth/verify", api.VerifyChallengeHandler)

	log.Fatal(r.Run(":8080"))
}
