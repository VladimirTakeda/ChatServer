package websocket

import (
	"ChatServer/pkg/service"
	"log"
	"time"
)

type ServerHub struct {
	Clients    map[*Client]bool
	ClientIds  map[int]*Client
	Broadcast  chan WsMessage
	Register   chan *Client
	Unregister chan *Client
	services   *service.Service
}

func NewServerHub(services *service.Service) *ServerHub {
	return &ServerHub{
		Broadcast:  make(chan WsMessage),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		ClientIds:  make(map[int]*Client),
		services:   services,
	}
}

func (h *ServerHub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				delete(h.ClientIds, client.ID)
				close(client.Message)
			}
		case message := <-h.Broadcast:
			//TODO Message ordering
			if userFrom, found := h.ClientIds[message.UserFrom]; found {
				select {
				case userFrom.Message <- WsMessageOut{message, true, time.Now().UTC()}:
				default:
					close(userFrom.Message)
					delete(h.Clients, userFrom)
				}
			} else {
				log.Printf("No Client: %v", message.Content)
			}
			if userTo, found := h.ClientIds[message.UserTo]; found {
				select {
				case userTo.Message <- WsMessageOut{message, false, time.Now().UTC()}:
				default:
					close(userTo.Message)
					delete(h.Clients, userTo)
				}
			} else {
				log.Printf("No Client: %v", message.Content)
			}
		}
	}
}
