package http

import (
	"ChatServer/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) GetUsersByPrefix(c *gin.Context) {
	logrus.Infoln("Got search user request")

	var getUsersByPrefixReq types.GetUsersByPrefixReq

	if err := c.BindJSON(&getUsersByPrefixReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	logrus.Infoln("With prefix: ", getUsersByPrefixReq.Prefix)

	users, err := h.services.GetUsersByPrefix(c, getUsersByPrefixReq.Prefix)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	logrus.Infoln("Answer: ", users)

	c.JSON(http.StatusOK, users)
}
