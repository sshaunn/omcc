package group

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
	"ohmycontrolcenter.tech/omcc/internal/common"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/config"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/logger"
	"regexp"
	"time"
)

type MessageHandler struct {
	bot             *tele.Bot
	log             logger.Logger
	cfg             *config.TelegramConfig
	commandPatterns []*regexp.Regexp
}

func NewGroupMessageHandler(cfg *config.TelegramConfig, bot *tele.Bot, log logger.Logger) *MessageHandler {
	var patterns []*regexp.Regexp
	for _, pattern := range cfg.CommandPatterns {
		patterns = append(patterns, regexp.MustCompile(pattern))
	}
	return &MessageHandler{
		bot:             bot,
		log:             log,
		cfg:             cfg,
		commandPatterns: patterns,
	}
}

func (h *MessageHandler) Handle(c tele.Context) error {
	message := c.Message()
	logFields := []logger.Field{
		logger.String("text", c.Text()),
		logger.String("chat_type", string(c.Chat().Type)),
		logger.Int64("chat_id", c.Chat().ID),
	}

	if message.ThreadID != 0 {
		logFields = append(logFields, logger.String("thread_id", fmt.Sprintf("%d", message.ThreadID)))
	}

	h.log.Info("Handling group message", logFields...)

	if !h.shouldProcessMessage(c) {
		return nil
	}

	if h.shouldDeleteMessage(message, c) {
		return h.deleteMessage(c)
	}
	return nil
}

func (h *MessageHandler) shouldProcessMessage(c tele.Context) bool {
	if !h.isTargetGroup(c.Chat().ID) {
		return false
	}

	threadId := c.Message().ThreadID
	if len(h.cfg.MonitoredTopics) > 0 && threadId != 0 {
		return h.isTargetTopics(threadId)
	}
	return true
}

func (h *MessageHandler) shouldDeleteMessage(msg *tele.Message, c tele.Context) bool {
	// no behaviour if test is empty
	if msg.Text == "" {
		return false
	}

	chatMember, _ := h.bot.ChatMemberOf(&tele.Chat{ID: c.Message().Chat.ID}, &tele.User{ID: msg.Sender.ID})
	if chatMember.Role == tele.Administrator || chatMember.Role == tele.Creator {
		return false
	}
	// check if its forbidden message by regex
	for _, pattern := range h.commandPatterns {
		if pattern.MatchString(msg.Text) {
			h.log.Info("detected command pattern in group",
				logger.String("text", msg.Text),
				logger.Int64("chat_id", msg.Chat.ID),
				logger.String("username", msg.Sender.Username),
			)
			return true
		}
	}

	return false
}

func (h *MessageHandler) deleteMessage(c tele.Context) error {
	msg := c.Message()

	err := h.bot.Delete(msg)
	if err != nil {
		h.log.Error("failed to delete message",
			logger.Error(err),
			logger.Int64("chat_id", msg.Chat.ID),
			logger.Int("message_id", msg.ID),
			logger.String("text", msg.Text),
		)
		return err
	}

	h.log.Info("message deleted",
		logger.Int64("chat_id", msg.Chat.ID),
		logger.Int("message_id", msg.ID),
		logger.String("text", msg.Text),
		logger.String("username", msg.Sender.Username),
	)

	//Optional: sending warning message for forbidden messages
	warning := fmt.Sprintf(common.UserWarningMessage, msg.Sender.Username)
	warningMsg, err := c.Bot().Send(msg.Chat, warning, &tele.SendOptions{
		ThreadID: msg.ThreadID,
	})
	if err != nil {
		h.log.Error("failed to send warning message", logger.Error(err))
		return err
	}

	// delete the warning message after 10 secs
	time.AfterFunc(10*time.Second, func() {
		_ = c.Bot().Delete(warningMsg)
	})

	return nil
}

func (h *MessageHandler) isCommandFormat(text string) bool {
	for _, pattern := range h.commandPatterns {
		if pattern.MatchString(text) {
			return true
		}
	}
	return false
}

func (h *MessageHandler) isTargetGroup(chatId int64) bool {
	for _, id := range h.cfg.MonitoredGroups {
		if id == chatId {
			return true
		}
	}
	return false
}

func (h *MessageHandler) isTargetTopics(threadId int) bool {
	for _, id := range h.cfg.MonitoredTopics {
		if id == threadId {
			return true
		}
	}
	return false
}
