package websocket

import (
	"ChatServer/pkg/types"
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan types.WsMessageOut
	userId   int
	deviceId string
}

func (c *Client) WriteUpdates(messages []types.WsMessageOut) {
	for _, message := range messages {
		c.Conn.WriteJSON(message)
	}
}

func NewClient(conn *websocket.Conn, id int, deviceId string) *Client {
	return &Client{
		Conn:     conn,
		Message:  make(chan types.WsMessageOut, 10),
		userId:   id,
		deviceId: deviceId,
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
		log.Printf("unregister")
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		var message types.WsMessage
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
