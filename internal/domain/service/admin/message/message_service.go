package message

import (
	"context"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/config"
	"ohmycontrolcenter.tech/omcc/pkg/exception"
	"ohmycontrolcenter.tech/omcc/pkg/logger"
)

type SendingMessageService struct {
	bot *tele.Bot
	cfg *config.Config
	log logger.Logger
}

func NewSendingMessageService(cfg *config.Config, bot *tele.Bot, log logger.Logger) *SendingMessageService {
	return &SendingMessageService{
		bot: bot,
		cfg: cfg,
		log: log,
	}
}

func (m *SendingMessageService) SendMessage(ctx context.Context, chatID int64, message string) error {
	_, err := m.bot.Send(&tele.Chat{ID: chatID}, message)
	if err != nil {
		return err
	}
	return nil
}

func (m *SendingMessageService) SendMessage2MultipleCustomers(ctx context.Context, chatIDs []int64, message string) ([]int64, error) {
	var errorIds []int64
	for _, chatID := range chatIDs {
		_, err := m.bot.Send(&tele.Chat{ID: chatID}, message)
		if err != nil {
			errorIds = append(errorIds, chatID)
			m.log.Info(fmt.Sprintf("Failed to send to this telegram user with userId=%d", chatID))
		}
	}
	if len(errorIds) > 0 {
		return errorIds, exception.ErrSendingMessage
	}
	return nil, nil
}

func (m *SendingMessageService) SendMessage2Groups(ctx context.Context, groupIDs []int64, message string) ([]int64, error) {
	for _, groupID := range groupIDs {
		_, err := m.bot.Send(&tele.Chat{ID: groupID}, message)
		if err != nil {
			m.log.Info(fmt.Sprintf("Failed to send to this telegram group with groupId=%d", groupID))
		}
	}
	return nil, nil
}
