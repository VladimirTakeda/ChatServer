package websocket

import (
	"ChatServer/pkg/service"
	"ChatServer/pkg/types"
	"context"
	"fmt"
	"log"
	"time"
)

type ServerHub struct {
	Clients    map[*Client]bool
	ClientIds  map[int]*Client
	Broadcast  chan types.WsMessage
	Register   chan *Client
	Unregister chan *Client
	Services   *service.Service
}

func NewServerHub(services *service.Service) *ServerHub {
	return &ServerHub{
		Broadcast:  make(chan types.WsMessage, 10),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		ClientIds:  make(map[int]*Client),
		Services:   services,
	}
}

func (h *ServerHub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				h.Services.SaveLastActiveTime(context.Background(), client.userId, client.deviceId, time.Now())
				delete(h.Clients, client)
				delete(h.ClientIds, client.userId)
				close(client.Message)
			}
		case message := <-h.Broadcast:
			log.Printf("Current channel size: %d", len(h.Broadcast))
			log.Printf("Current message: %s", message.Content)
			ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

			err := h.Services.AddMessage(ctx, message.UserFrom, message.ChatTo, message.Content)
			if err != nil {
				log.Printf("Failed save message: %s", err.Error())
			}

			members, err := h.Services.Chat.GetChatMembers(ctx, message.ChatTo)
			if err != nil {
				log.Printf("Failed to get chat members: %s", err.Error())
			}

			log.Printf("Содержимое среза: %v", members)

			for _, member := range members {
				if client, found := h.ClientIds[member]; found {
					select {
					case client.Message <- types.WsMessageOut{WsMessage: message, Time: time.Now().UTC()}:
					default:
						fmt.Printf("Member closed %d: ", member)
						close(client.Message)
						delete(h.Clients, client)
					}
				} else {
					log.Printf("No Client: %v", message.Content)
				}
			}
		}
	}
}
