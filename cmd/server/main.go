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
	
	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})
	
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"service": "CoScribe",
			"version": "1.0.0",
		})
	})

	r.GET("/ws/room", ws.RoomHandler)
	r.GET("/ws/document", ws.DocumentHandler)
	
	r.GET("/api/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "API is working"})
	})
	
	r.POST("/api/document/:id", func(c *gin.Context) {
		docID := c.Param("id")
		var requestBody struct {
			Content string `json:"content"`
		}
		
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON"})
			return
		}
		
		doc := ws.GlobalDocumentManager.GetDocument(docID)
		doc.SetContent(requestBody.Content)
		
		if err := ws.GlobalDocumentManager.SaveDocument(docID); err != nil {
			c.JSON(500, gin.H{"error": "Failed to save document"})
			return
		}
		
		c.JSON(200, gin.H{
			"success": true,
			"version": doc.GetVersion(),
			"content": doc.GetContent(),
		})
	})

	log.Fatal(r.Run(":8080"))
}