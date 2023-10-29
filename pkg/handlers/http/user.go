package http

import (
	"ChatServer/pkg/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) registerUser(c *gin.Context) {

	var registerUser types.RegisterUser

	if err := c.BindJSON(&registerUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	userId, err := h.services.Register(c, registerUser.Nickname)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{"userId": userId})
}
