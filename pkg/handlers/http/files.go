package http

import (
	"ChatServer/pkg/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (h *Handler) uploadFile(ctx *gin.Context) {
	userId, _ := strconv.Atoi(ctx.Query("userId"))
	chatId, _ := strconv.Atoi(ctx.Query("chatId"))
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	logrus.Debugln(header.Size)
	logrus.Debugln(header.Filename)
	logrus.Debugln(file)
	fileId, err := h.services.SaveFile(ctx, userId, chatId, file, header.Size, header.Filename)
	if err != nil {
		logrus.Debugln("Can't save file : " + err.Error())
		return
	}

	// Возвращаем ответ клиенту
	ctx.JSON(http.StatusOK, gin.H{"file_id": fileId})
}

func (h *Handler) downloadFile(ctx *gin.Context) {
	var getFileReq types.GetFileReq

	if err := ctx.BindJSON(&getFileReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	Object, err := h.services.LoadFile(ctx, getFileReq.FileId)
	if err != nil {
		logrus.Debugln("Can't load file : " + err.Error())
		return
	}

	objInfo, err := Object.Stat()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении информации об объекте"})
		return
	}

	// Устанавливаем заголовки ответа
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", objInfo.Key))
	ctx.Header("Content-Type", objInfo.ContentType)
	ctx.Header("Content-Length", fmt.Sprintf("%d", objInfo.Size))

	// Копируем содержимое объекта в ответ
	ctx.DataFromReader(http.StatusOK, objInfo.Size, objInfo.ContentType, Object, nil)
}
