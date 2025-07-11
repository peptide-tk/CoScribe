package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var GlobalHub = NewHub()

type Message struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
	Room    string      `json:"room"`
	User    string      `json:"user"`
	Time    time.Time   `json:"time"`
}

func EchoHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Printf("WebSocket write error: %v", err)
			break
		}
	}
}

// RoomHandler
func RoomHandler(c *gin.Context) {
	roomID := c.Query("room")
	if roomID == "" {
		roomID = "default"
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// クライアント作成
	client := &Client{
		ID:   generateClientID(),
		Conn: conn,
		Room: &Room{ID: roomID},
		Send: make(chan []byte, 256),
	}

	// ハブに登録
	GlobalHub.Register <- client

	// ゴルーチンを開始
	go client.writePump()
	go client.readPump(GlobalHub)
}

// クライアントIDを生成
func generateClientID() string {
	return time.Now().Format("20060102-150405-") + time.Now().Format("000")
}

// クライアントからのメッセージを読み取り
func (c *Client) readPump(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, messageBytes, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		var msg Message
		if err := json.Unmarshal(messageBytes, &msg); err != nil {
			log.Printf("Message unmarshal error: %v", err)
			continue
		}

		msg.Time = time.Now()
		msg.Room = c.Room.ID

		// ブロードキャスト
		messageBytes, _ = json.Marshal(msg)
		hub.BroadcastToRoom(c.Room.ID, messageBytes, c)
	}
}

// メッセージを送信
func (c *Client) writePump() {
	defer c.Conn.Close()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}
		}
	}
}