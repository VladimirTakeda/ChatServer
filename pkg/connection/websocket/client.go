package websocket

import (
	"ChatServer/pkg/service"
	"ChatServer/pkg/types"
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type Client struct {
	Conn     *websocket.Conn
	userId   int
	deviceId string
	service  *service.Service
}

func (c *Client) WriteUpdates(messages []types.WsMessageWithTime) {
	for _, message := range messages {
		c.Conn.WriteJSON(message)
	}
}

func NewClient(conn *websocket.Conn, id int, deviceId string, service *service.Service) *Client {
	return &Client{
		Conn:     conn,
		userId:   id,
		deviceId: deviceId,
		service:  service,
	}
}

// Takes a message from PubSub subscruptions and pass it to user via WebSocket
func (c *Client) ConsumeMessagesFromPubSub(ctx context.Context, subscriber *redis.PubSub) {
	log.Printf("Start consuming messages from PubSub")
	defer func() {
		c.Conn.Close()
	}()

	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			log.Printf("error: %v", err)
			return
		}

		log.Printf("Got message: %s", msg.Payload)

		var message types.WsMessage

		if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
			log.Printf("error: %v", err)
			return
		}

		c.Conn.WriteJSON(message)
	}
}

// Consumes message from WebSocket connection and process it
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

		members, err := c.service.GetChatMembers(context.Background(), message.ChatTo)
		if err != nil {
			return
		}

		err = c.service.AddMessage(context.Background(),
			types.WsMessageWithTime{WsMessage: message, Time: time.Now().UTC()},
			members)
		if err != nil {
			log.Printf("error: %v", err)
		}
	}
}
