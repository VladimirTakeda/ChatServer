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

func (h *Handler) WebSocketCreate(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	userTo, err := strconv.Atoi(ctx.Request.URL.Query().Get("user_id"))
	if err != nil {
		log.Printf("Error parsing userTo: %v", err)
		return
	}

	deviceId := ctx.Request.URL.Query().Get("device_id")

	log.Printf("New User Connection with Id: %v and device %s", userTo, deviceId)

	client := websocket2.NewClient(conn, userTo, deviceId, h.services)

	h.server.ClientIds[userTo] = client

	h.server.Register <- client

	// TODO: now user has subscription only for his private messages
	// TODO: in future we need to add subscription for every group chat
	subscriber, err := h.server.Services.CreateSubscription(ctx, "user:"+strconv.Itoa(userTo))
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	log.Printf("Subscription %s created", "user:"+strconv.Itoa(userTo))

	go func() {
		messages, err := h.server.Services.GetMissedMessages(ctx, userTo, deviceId)
		if err != nil {
		}
		client.WriteUpdates(messages)
		go client.ConsumeMessagesFromPubSub(ctx, subscriber)
	}()

	// send messages from websocket to PubSub
	client.ReadMessage(h.server)
}
