package main

import (
	"log"

	"github.com/cyriljohn147/zero-trust-backend/internal/api"
	"github.com/cyriljohn147/zero-trust-backend/internal/db"
	"github.com/gin-gonic/gin"
)

func main() {
	db.Connect()
	defer db.Close()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.POST("/devices/register", api.RegisterDeviceHandler)

	log.Fatal(r.Run(":8080"))
}
