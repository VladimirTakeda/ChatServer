package websocket

import (
	"ChatServer/pkg/service"
	"context"
	"time"
)

type ServerHub struct {
	Clients    map[*Client]bool
	ClientIds  map[int]*Client
	Register   chan *Client
	Unregister chan *Client
	Services   *service.Service
}

func NewServerHub(services *service.Service) *ServerHub {
	return &ServerHub{
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
				err := h.Services.SaveLastActiveTime(context.Background(), client.userId, client.deviceId, time.Now())
				if err != nil {
					return
				}
				delete(h.Clients, client)
				delete(h.ClientIds, client.userId)
			}
		}
	}
}
