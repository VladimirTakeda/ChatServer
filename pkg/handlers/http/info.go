package http

import (
	"ChatServer/pkg/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetUsersByPrefix(c *gin.Context) {

	var getUsersByPrefixReq types.GetUsersByPrefixReq

	if err := c.BindJSON(&getUsersByPrefixReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	users, err := h.services.GetUsersByPrefix(c, getUsersByPrefixReq.Prefix)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	c.JSON(http.StatusOK, users)
}
