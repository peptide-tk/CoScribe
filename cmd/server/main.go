package main

import (
	"fmt"
	"log"

	"coscribe/internal/ws"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting CoScribe Server...")

	go ws.GlobalHub.Run()

	r := gin.Default()
	
	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"service": "CoScribe",
			"version": "1.0.0",
		})
	})

	// WebSocket endpoints
	r.GET("/ws/room", ws.RoomHandler)
	r.GET("/ws/document", ws.DocumentHandler)

	fmt.Println("Server starting on :8080")
	log.Fatal(r.Run(":8080"))
}