package bot

import (
	"context"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"log"
	"ohmycontrolcenter.tech/omcc/internal/common"
	"ohmycontrolcenter.tech/omcc/internal/domain/bot/handler/group"
	"ohmycontrolcenter.tech/omcc/internal/domain/bot/handler/private"
	"ohmycontrolcenter.tech/omcc/internal/domain/service"
	"ohmycontrolcenter.tech/omcc/internal/domain/service/exchange"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/config"
	"ohmycontrolcenter.tech/omcc/internal/middleware"
	"ohmycontrolcenter.tech/omcc/pkg/logger"
)

type TelegramBot struct {
	Bot        *tele.Bot
	cfg        *config.Config
	log        logger.Logger
	middleware *middleware.Manager
}

func NewTelegramBot(cfg *config.Config, log logger.Logger, middleware *middleware.Manager) (*TelegramBot, error) {
	log.Info("initializing telegram Bot",
		logger.String("webhook_url", cfg.Telegram.WebhookURL),
	)

	webhook := &tele.Webhook{
		Listen: cfg.Telegram.Port,
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
		//Poller: &tele.LongPoller{
		//	Timeout:        10 * time.Second,
		//	AllowedUpdates: []string{"message", "callback_query"},
		//},
		//URL: "https://api.telegram.org",
	}

	b, err := tele.NewBot(settings)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram Bot: %w", err)
	}

	tb := &TelegramBot{
		Bot:        b,
		cfg:        cfg,
		log:        log,
		middleware: middleware,
	}

	// 注册命令处理器
	tb.registerHandlers()

	return tb, nil
}

func (t *TelegramBot) checkWebhookStatus() error {
	webhook, err := t.Bot.Webhook()
	if err != nil {
		return fmt.Errorf("failed to get webhook info: %w", err)
	}

	log.Printf("Webhook status - URL: %s, Pending updates: %d",
		webhook.Endpoint.PublicURL, webhook.PendingUpdates)

	return nil
}

func (t *TelegramBot) Start(ctx context.Context) error {
	go func() {
		t.Bot.Start()
	}()

	return nil
}

// Stop the telegram Bot
func (t *TelegramBot) Stop() {
	t.Bot.Stop()
}

// registerHandlers 注册命令处理器
func (t *TelegramBot) registerHandlers() {
	middlewareHandler := t.middleware.TelegramMiddleware

	groupHandler := group.NewGroupMessageHandler(&t.cfg.Telegram, t.Bot, t.log)

	bitgetClient := exchange.NewBitgetClient(&t.cfg.Exchange.BitgetConfig, t.log)
	verifyService := service.NewVerifyService(t.cfg, bitgetClient, t.log)
	volumeService := service.NewVolumeService(&t.cfg.Database, bitgetClient, t.log)
	checkService := service.NewStatusService(t.cfg, t.log)
	accountService := service.NewAccountService(t.cfg, t.log)

	verifyCommand := private.NewVerifyCommand(t.Bot, t.log, *verifyService)
	volumeCommand := private.NewVolumeCommand(t.log, *volumeService)
	startCommand := private.NewStartCommand(t.log)
	checkCommand := private.NewCheckCommand(t.log, *checkService)
	helpCommand := private.NewHelpCommand(t.log)
	accountCommand := private.NewAccountCommand(t.Bot, t.log, *accountService)
	onTextCommand := private.NewOnTextCommand(t.log)

	// processing non-command text message
	t.Bot.Handle(tele.OnText, middlewareHandler(handlerType(onTextCommand.Handle, groupHandler.Handle)))

	// register /start command
	t.Bot.Handle(common.StartCommandName, middlewareHandler(handlerType(startCommand.Handle, groupHandler.Handle)))
	// register /help command
	t.Bot.Handle(common.HelpCommandName, middlewareHandler(handlerType(helpCommand.Handle, groupHandler.Handle)))
	// register /verify command
	t.Bot.Handle(common.VerifyCommandName, middlewareHandler(handlerType(verifyCommand.Handle, groupHandler.Handle)))
	// register /volume command
	t.Bot.Handle(common.VolumeCommandName, middlewareHandler(handlerType(volumeCommand.Handle, groupHandler.Handle)))
	// register /check command
	t.Bot.Handle(common.StatusCommandName, middlewareHandler(handlerType(checkCommand.Handle, groupHandler.Handle)))
	// register /account command
	t.Bot.Handle(common.AccountCommandName, middlewareHandler(handlerType(accountCommand.Handle, groupHandler.Handle)))

}

func handlerType(handlerFunc ...tele.HandlerFunc) middleware.Handler {
	if len(handlerFunc) < 2 {
		return middleware.Handler{
			PrivateHandler: handlerFunc[0],
		}
	}
	return middleware.Handler{
		PrivateHandler:    handlerFunc[0],
		SuperGroupHandler: handlerFunc[1],
		DefaultHandler:    handlerFunc[1], // TODO added if needed
	}

}

// SendMessage 发送消息
func (t *TelegramBot) SendMessage(ctx context.Context, chatID int64, message string) error {
	_, err := t.Bot.Send(&tele.Chat{ID: chatID}, message)
	return err
}
