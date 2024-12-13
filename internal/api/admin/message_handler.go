package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ohmycontrolcenter.tech/omcc/internal/domain/model"
	"ohmycontrolcenter.tech/omcc/internal/domain/service/admin/message"
	"ohmycontrolcenter.tech/omcc/pkg/logger"
)

type MessageHandler struct {
	MessageService *message.SendingMessageService
	log            logger.Logger
}

func NewMessageHandler(messageService *message.SendingMessageService, log logger.Logger) *MessageHandler {
	return &MessageHandler{
		MessageService: messageService,
		log:            log,
	}
}

func (m *MessageHandler) SendMessage2MultipleCustomers(c *gin.Context) {
	var messagingRequest *model.AdminMessagingRequest
	if err := c.ShouldBindJSON(&messagingRequest); err != nil {
		m.log.Error("invalid request parameters",
			logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Failure",
			"error":  "Invalid request parameters",
		})
		return
	}
	ids, err := m.MessageService.SendMessage2MultipleCustomers(c.Request.Context(), messagingRequest.Ids, messagingRequest.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"id_list": ids,
		"message": messagingRequest.Message,
	})
}

func (m *MessageHandler) SendMessage2Groups(c *gin.Context) {
	var messagingRequest *model.AdminMessagingRequest
	if err := c.ShouldBindJSON(&messagingRequest); err != nil {
		m.log.Error("invalid request parameters",
			logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Failure",
			"error":  "Invalid request parameters",
		})
		return
	}
	ids, err := m.MessageService.SendMessage2Groups(c.Request.Context(), messagingRequest.Ids, messagingRequest.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"id_list": ids,
		"message": messagingRequest.Message,
	})
}
