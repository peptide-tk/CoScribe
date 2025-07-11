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
	
	r.GET("/ws/echo", ws.EchoHandler)
	r.GET("/ws/room", ws.RoomHandler)

	fmt.Println("Server starting on :8080")
	log.Fatal(r.Run(":8080"))
}