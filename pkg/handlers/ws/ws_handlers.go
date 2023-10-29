package ws

import (
	"ChatServer/pkg/connection/websocket"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	server *websocket.ServerHub
}

func NewHandler(s *websocket.ServerHub) *Handler {
	return &Handler{
		server: s,
	}
}

func (h *Handler) SetupRoutes(router *gin.Engine) {
	router.GET("/create", h.WebSocketCreate)
}
