package main

import (
	"fmt"
	"log"

	"coscribe/internal/database"
	"coscribe/internal/document"
	"coscribe/internal/store"
	"coscribe/internal/ws"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting CoScribe Server...")

	if err := database.Initialize(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.Close()
	documentStore := store.NewDocumentStore()
	
	ws.GlobalDocumentManager = document.NewManager(documentStore)

	go ws.GlobalHub.Run()

	r := gin.Default()
	
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"service": "CoScribe",
			"version": "1.0.0",
		})
	})

	r.GET("/ws/room", ws.RoomHandler)
	r.GET("/ws/document", ws.DocumentHandler)

	fmt.Println("Server starting on :8080")
	log.Fatal(r.Run(":8080"))
}