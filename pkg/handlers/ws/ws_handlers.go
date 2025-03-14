package ws

import (
	"ChatServer/pkg/connection/websocket"
	"ChatServer/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	server   *websocket.ServerHub
	services *service.Service
}

func NewHandler(s *websocket.ServerHub, services *service.Service) *Handler {
	return &Handler{
		server:   s,
		services: services,
	}
}

func (h *Handler) SetupRoutes(router *gin.Engine) {
	router.GET("/create", h.WebSocketCreate)
}
