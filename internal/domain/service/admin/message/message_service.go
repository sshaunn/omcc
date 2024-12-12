package message

import (
	"context"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/config"
	"ohmycontrolcenter.tech/omcc/pkg/exception"
	"ohmycontrolcenter.tech/omcc/pkg/logger"
	"sync"
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
	// 创建用于收集错误的 channel
	errChan := make(chan int64, len(chatIDs))

	// 创建 WaitGroup 来等待所有 goroutine 完成
	var wg sync.WaitGroup

	// 并发发送消息
	for _, chatID := range chatIDs {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()

			// 发送消息
			_, err := m.bot.Send(&tele.Chat{ID: id}, message)
			if err != nil {
				m.log.Info(fmt.Sprintf("Failed to send to this telegram user with userId=%d", id))
				errChan <- id
			}
		}(chatID)
	}

	// 在新的 goroutine 中等待所有发送完成并关闭 channel
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// 收集失败的 ID
	var errorIds []int64
	for id := range errChan {
		errorIds = append(errorIds, id)
	}

	// 返回结果
	if len(errorIds) > 0 {
		return errorIds, exception.ErrSendingMessage
	}
	return nil, nil
}

func (m *SendingMessageService) SendMessage2Groups(ctx context.Context, groupIDs []int64, message string) ([]int64, error) {
	errChan := make(chan int64, len(groupIDs))
	var wg sync.WaitGroup
	for _, groupID := range groupIDs {
		wg.Add(1)
		go func(id int64) {
			_, err := m.bot.Send(&tele.Chat{ID: groupID}, message)
			if err != nil {
				m.log.Info(fmt.Sprintf("Failed to send to this telegram group with groupId=%d", groupID))
			}
		}(groupID)

	}
	go func() {
		wg.Wait()
		close(errChan)
	}()
	return nil, nil
}
