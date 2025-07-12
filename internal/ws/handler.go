package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"coscribe/internal/document"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	EnableCompression: false,
}

var GlobalHub = NewHub()
var GlobalDocumentManager *document.Manager

type Message struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
	Room    string      `json:"room"`
	User    string      `json:"user"`
	Time    time.Time   `json:"time"`
}

type DocumentMessage struct {
	Type     string         `json:"type"`
	Document string         `json:"document"`
	Edit     *document.Edit `json:"edit,omitempty"`
	Content  string         `json:"content,omitempty"`
	Version  int            `json:"version,omitempty"`
	User     string         `json:"user"`
	Time     time.Time      `json:"time"`
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

	client := &Client{
		ID:   generateClientID(),
		Conn: conn,
		Room: &Room{ID: roomID},
		Send: make(chan []byte, 256),
	}

	GlobalHub.Register <- client

	go client.writePump()
	go client.readPump(GlobalHub)
}

// DocumentHandler
func DocumentHandler(c *gin.Context) {
	docID := c.Query("doc")
	if docID == "" {
		docID = "default"
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		ID:   generateClientID(),
		Conn: conn,
		Room: &Room{ID: docID},
		Send: make(chan []byte, 256),
	}

	GlobalHub.Register <- client

	doc := GlobalDocumentManager.GetDocument(docID)
	initialMsg := DocumentMessage{
		Type:     "document_state",
		Document: docID,
		Content:  doc.GetContent(),
		Version:  doc.GetVersion(),
		Time:     time.Now(),
	}
	
	if msgBytes, err := json.Marshal(initialMsg); err == nil {
		client.Send <- msgBytes
	}

	go client.writePump()
	go client.readDocumentPump(GlobalHub)
}

func generateClientID() string {
	return time.Now().Format("20060102-150405-") + time.Now().Format("000")
}

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

		messageBytes, _ = json.Marshal(msg)
		hub.BroadcastToRoom(c.Room.ID, messageBytes, c)
	}
}

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

// readDocumentPump
func (c *Client) readDocumentPump(hub *Hub) {
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

		var msg DocumentMessage
		if err := json.Unmarshal(messageBytes, &msg); err != nil {
			log.Printf("Document message unmarshal error: %v", err)
			continue
		}

		msg.Time = time.Now()
		msg.Document = c.Room.ID

		switch msg.Type {
		case "edit":
			if msg.Edit != nil {
				msg.Edit.User = msg.User
				err := GlobalDocumentManager.ApplyEdit(c.Room.ID, msg.Edit)
				
				if err != nil {
					errorMsg := DocumentMessage{
						Type:     "error",
						Document: c.Room.ID,
						Content:  err.Error(),
						Time:     time.Now(),
					}
					if errBytes, marshalErr := json.Marshal(errorMsg); marshalErr == nil {
						c.Send <- errBytes
					}
					continue
				}

				msg.Version = GlobalDocumentManager.GetDocument(c.Room.ID).GetVersion()
				messageBytes, _ = json.Marshal(msg)
				hub.BroadcastToRoom(c.Room.ID, messageBytes, c)
			}

		case "request_document":
			doc := GlobalDocumentManager.GetDocument(c.Room.ID)
			stateMsg := DocumentMessage{
				Type:     "document_state",
				Document: c.Room.ID,
				Content:  doc.GetContent(),
				Version:  doc.GetVersion(),
				Time:     time.Now(),
			}
			if stateBytes, err := json.Marshal(stateMsg); err == nil {
				c.Send <- stateBytes
			}

		case "document_update":
			if msg.Content != "" {
				doc := GlobalDocumentManager.GetDocument(c.Room.ID)
				doc.SetContent(msg.Content)
				
				if err := GlobalDocumentManager.SaveDocument(c.Room.ID); err != nil {
					log.Printf("Failed to save document: %v", err)
				} else {
					log.Printf("Document saved successfully")
				}
				
				successMsg := DocumentMessage{
					Type:     "document_updated",
					Document: c.Room.ID,
					Content:  doc.GetContent(),
					Version:  doc.GetVersion(),
					Time:     time.Now(),
				}
				if successBytes, err := json.Marshal(successMsg); err == nil {
					c.Send <- successBytes
					log.Printf("Sent document_updated response")
				}
			} else {
				log.Printf("Empty content in document_update message")
			}
		}
	}
}