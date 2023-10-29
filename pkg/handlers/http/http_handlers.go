package http

import (
	"ChatServer/pkg/service"
	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
	"time"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	c.String(429, "Too many requests. Try again in "+time.Until(info.ResetTime).String())
}

func (h *Handler) SetupRoutes(router *gin.Engine) {
	store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  time.Second,
		Limit: 80,
	})

	rateLimiter := ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})

	chat := router.Group("/chats")
	{
		chat.POST("/", rateLimiter, h.createChat)
		chat.DELETE("/:chat_id", rateLimiter, h.deleteChat)
	}

	user := router.Group("/user")
	{
		user.POST("/register", rateLimiter, h.registerUser)
	}

	router.POST("/get_users_by_prefix", rateLimiter, h.GetUsersByPrefix)
}
