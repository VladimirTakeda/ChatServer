package http

import (
	"ChatServer/pkg/types"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getHash(info types.DeviceInfo) string {
	infoJSON, err := json.Marshal(info)
	if err != nil {
		panic(err)
	}

	hash := sha256.Sum256(infoJSON)

	hashString := hex.EncodeToString(hash[:])

	return hashString
}

func (h *Handler) registerUser(c *gin.Context) {

	var registerUser types.RegisterUser

	if err := c.BindJSON(&registerUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	userId, deviceId, err := h.services.RegisterUser(c, registerUser.Nickname, getHash(registerUser.DeviceInfo))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{"userId": userId, "deviceId": deviceId})
}
