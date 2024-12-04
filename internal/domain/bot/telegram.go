package bot

import (
	"context"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"log"
	"ohmycontrolcenter.tech/omcc/internal/api/middleware"
	"ohmycontrolcenter.tech/omcc/internal/domain/bot/handler"
	"ohmycontrolcenter.tech/omcc/internal/domain/service"
	"ohmycontrolcenter.tech/omcc/internal/domain/service/exchange"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/config"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/logger"
)

type TelegramBot struct {
	bot        *tele.Bot
	cfg        *config.Config
	log        logger.Logger
	middleware *middleware.Manager
}

func NewTelegramBot(cfg *config.Config, log logger.Logger, middleware *middleware.Manager) (*TelegramBot, error) {
	log.Info("initializing telegram bot",
		logger.String("webhook_url", cfg.Telegram.WebhookURL),
	)

	webhook := &tele.Webhook{
		Listen: ":8080",
		Endpoint: &tele.WebhookEndpoint{
			PublicURL: cfg.Telegram.WebhookURL,
		},
		AllowedUpdates: []string{"message", "callback_query"},
		MaxConnections: 40,
	}

	settings := tele.Settings{
		Token:  cfg.Telegram.Token,
		Poller: webhook,
		//Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(settings)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %w", err)
	}

	tb := &TelegramBot{
		bot:        b,
		cfg:        cfg,
		log:        log,
		middleware: middleware,
	}

	// 注册命令处理器
	tb.registerHandlers()

	return tb, nil
}

func (t *TelegramBot) checkWebhookStatus() error {
	webhook, err := t.bot.Webhook()
	if err != nil {
		return fmt.Errorf("failed to get webhook info: %w", err)
	}

	log.Printf("Webhook status - URL: %s, Pending updates: %d",
		webhook.Endpoint.PublicURL, webhook.PendingUpdates)

	return nil
}

func (t *TelegramBot) Start(ctx context.Context) error {
	go func() {
		t.bot.Start()
	}()

	return nil
}

// Stop the telegram bot
func (t *TelegramBot) Stop() {
	t.bot.Stop()
}

// registerHandlers 注册命令处理器
func (t *TelegramBot) registerHandlers() {
	middlewareHandler := t.middleware.TelegramMiddleware

	bitgetClient := exchange.NewBitgetClient(&t.cfg.Bitget, t.log)
	verifyService := service.NewVerifyService(t.cfg, bitgetClient, t.log)

	verifyCommand := handler.NewVerifyCommand(t.log, *verifyService)
	startCommand := handler.NewStartCommand(t.log)
	helpCommand := handler.NewHelpCommand(t.log)
	onTextCommand := handler.NewOnTextCommand(t.log)

	// register /start command
	t.bot.Handle("/start", middlewareHandler(startCommand.Handle))
	// register /help command
	t.bot.Handle("/help", middlewareHandler(helpCommand.Handle))
	// register /verify command
	t.bot.Handle("/verify", middlewareHandler(verifyCommand.Handle))
	// processing non-command text message
	t.bot.Handle(tele.OnText, middlewareHandler(onTextCommand.Handle))
}

// SendMessage 发送消息
func (t *TelegramBot) SendMessage(ctx context.Context, chatID int64, message string) error {
	_, err := t.bot.Send(&tele.Chat{ID: chatID}, message)
	return err
}
