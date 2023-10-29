package ws

import (
	websocket2 "ChatServer/pkg/connection/websocket"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) WebSocketCreate(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	userTo, err := strconv.Atoi(c.Request.URL.Query().Get("user_id"))
	if err != nil {
		log.Printf("Error parsing userTo: %v", err)
		return
	}

	log.Printf("New User Connection with Id: %v", userTo)

	client := websocket2.NewClient(conn, userTo)
	h.server.ClientIds[userTo] = client

	h.server.Register <- client

	go client.WriteMessage()
	client.ReadMessage(h.server)
}
