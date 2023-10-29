package websocket

import "log"

type ServerHub struct {
	Clients    map[*Client]bool
	ClientIds  map[int]*Client
	Broadcast  chan WsMessage
	Register   chan *Client
	Unregister chan *Client
}

func NewServerHub() *ServerHub {
	return &ServerHub{
		Broadcast:  make(chan WsMessage),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		ClientIds:  make(map[int]*Client),
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
				close(client.Message)
			}
		case message := <-h.Broadcast:
			if client, found := h.ClientIds[message.UserTo]; found {
				select {
				case client.Message <- message:
				default:
					close(client.Message)
					delete(h.Clients, client)
				}
			} else {
				log.Printf("No Client: %v", message.Content)
			}
		}
	}
}
