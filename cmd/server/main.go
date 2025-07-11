package main

import (
	"fmt"
	"log"
	"net/http"

	"coscribe/internal/ws"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting CoScribe Server...")

	go ws.GlobalHub.Run()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World from CoScribe!",
			"version": "1.0.0"
			,
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/ws/echo", ws.EchoHandler)
	r.GET("/ws/room", ws.RoomHandler)

	fmt.Println("Server starting on :8080")
	log.Fatal(r.Run(":8080"))
}