package websocket

import (
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	Conn    *websocket.Conn
	Message chan WsMessage
	ID      int
}

func NewClient(conn *websocket.Conn, id int) *Client {
	return &Client{
		Conn:    conn,
		Message: make(chan WsMessage, 10),
		ID:      id,
	}
}

func (c *Client) WriteMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Message
		log.Printf("Write Message: %v", message.Content)
		if !ok {
			return
		}

		c.Conn.WriteJSON(message)
	}
}

func (c *Client) ReadMessage(hub *ServerHub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		var message WsMessage
		err := c.Conn.ReadJSON(&message)
		if err != nil {
			log.Printf("error: %v", err)
		}

		log.Printf("New Incoming Message: %v", message.Content)

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		hub.Broadcast <- message
	}
}

type WsMessage struct {
	Content  string `json:"content"`
	UserFrom int    `json:"user_from_id"`
	UserTo   int    `json:"user_to_id"`
}
