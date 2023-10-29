package http

import (
	"ChatServer/pkg/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) createChat(c *gin.Context) {

	var chatInfo types.ChatInfo

	if err := c.BindJSON(&chatInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	chatId, err := h.services.CreateChat(c, []int{chatInfo.FromUser, chatInfo.ToUser})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{"chatId": chatId})
}

func (h *Handler) deleteChat(c *gin.Context) {

	var deleteChatInfo types.DeleteChatInfo

	if err := c.BindJSON(&deleteChatInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	err := h.services.DeleteChat(c, deleteChatInfo.ChatId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
