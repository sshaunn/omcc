package middleware

import (
	"errors"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"ohmycontrolcenter.tech/omcc/internal/common"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/logger"
	"ohmycontrolcenter.tech/omcc/pkg/exception"
	"time"
)

type Manager struct {
	log logger.Logger
}

type Handler struct {
	PrivateHandler    tele.HandlerFunc
	SuperGroupHandler tele.HandlerFunc
	DefaultHandler    tele.HandlerFunc
}

func NewManager(log logger.Logger) *Manager {
	return &Manager{
		log: log,
	}
}

type MessageInfo struct {
	fields   []logger.Field
	chatType tele.ChatType
}

// createTelegramMessageInfo 创建消息日志信息
func createTelegramMessageInfo(c tele.Context) MessageInfo {
	return MessageInfo{
		fields: []logger.Field{
			logger.Int64("chat_id", c.Chat().ID),
			logger.String("chat_type", string(c.Chat().Type)),
			logger.String("text", c.Text()),
			logger.String("username", c.Sender().Username),
			logger.String("first_name", c.Sender().FirstName),
			logger.String("last_name", c.Sender().LastName),
			logger.Int64("user_id", c.Sender().ID),
			logger.Any("message_id", c.Message().ID),
		},
		chatType: c.Chat().Type,
	}
}

func (m *Manager) getHandlerForChatType(handlers Handler, chatType tele.ChatType) tele.HandlerFunc {
	m.log.Info("getting handler for chat type",
		logger.String("chat_type", string(chatType)),
		logger.Any("has_private_handler", handlers.PrivateHandler != nil),
		logger.Any("has_supergroup_handler", handlers.SuperGroupHandler != nil),
	)
	switch chatType {
	case tele.ChatPrivate:
		return handlers.PrivateHandler
	case tele.ChatSuperGroup:
		return handlers.SuperGroupHandler
	default:
		return handlers.DefaultHandler
	}
}

func (m *Manager) logReceived(msgInfo MessageInfo) {
	var msgType string
	switch msgInfo.chatType {
	case tele.ChatPrivate:
		msgType = "private"
	case tele.ChatSuperGroup:
		msgType = "supergroup"
	default:
		msgType = string(msgInfo.chatType)
	}

	m.log.Info(fmt.Sprintf("received telegram %s message", msgType), msgInfo.fields...)
}

// handleError 处理错误并记录日志
func (m *Manager) handleError(err error, c tele.Context, msgInfo MessageInfo, duration time.Duration) error {
	if err == nil {
		m.logSuccess(msgInfo, duration)
		return nil
	}

	var cmdErr *exception.CommandError
	if errors.As(err, &cmdErr) {
		return m.handleCommandError(cmdErr, c, msgInfo, duration)
	}

	m.logError(err, msgInfo, duration)
	return err
}

// handleCommandError 处理命令错误
func (m *Manager) handleCommandError(cmdErr *exception.CommandError, c tele.Context,
	msgInfo MessageInfo, duration time.Duration) error {

	m.log.Info("command error",
		append(msgInfo.fields,
			logger.String("error_type", fmt.Sprintf("%d", cmdErr.Type)),
			logger.String("error_message", cmdErr.Message),
			logger.Duration("duration", duration),
		)...,
	)
	return c.Send(cmdErr.Message)
}

// logSuccess 记录成功日志
func (m *Manager) logSuccess(msgInfo MessageInfo, duration time.Duration) {
	m.log.Info("Completed handled telegram message success",
		append(msgInfo.fields,
			logger.Duration("duration", duration),
		)...,
	)
}

// logError 记录错误日志
func (m *Manager) logError(err error, msgInfo MessageInfo, duration time.Duration) {
	m.log.Error("failed to handle telegram message",
		append(msgInfo.fields,
			logger.Error(err),
			logger.Duration("duration", duration),
		)...,
	)
}

// logPanic 记录 panic 日志
func (m *Manager) logPanic(r interface{}, msgInfo MessageInfo, duration time.Duration) {
	m.log.Error("recovered from panic",
		append(msgInfo.fields,
			logger.Any("panic", r),
			logger.Duration("duration", duration),
		)...,
	)
}

// TelegramMiddleware 统一的 Telegram 中间件
func (m *Manager) TelegramMiddleware(handlers Handler) tele.HandlerFunc {
	return func(c tele.Context) error {
		msgInfo := createTelegramMessageInfo(c)

		// 记录收到的消息
		m.logReceived(msgInfo)

		// 获取对应聊天类型的处理器
		handler := m.getHandlerForChatType(handlers, c.Chat().Type)
		if handler == nil {
			return nil
		}

		start := time.Now()
		var err error

		// panic 恢复
		func() {
			defer func() {
				if r := recover(); r != nil {
					m.logPanic(r, msgInfo, time.Since(start))
					err = c.Send(common.InternalServerErrorMessage)
				}
			}()
			err = handler(c)
		}()

		return m.handleError(err, c, msgInfo, time.Since(start))
	}
}
