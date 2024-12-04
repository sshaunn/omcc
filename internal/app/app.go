package app

import (
	"context"
	"gorm.io/gorm"
	"ohmycontrolcenter.tech/omcc/internal/api/middleware"
	"ohmycontrolcenter.tech/omcc/internal/domain/bot"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/config"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/logger"
)

type App struct {
	cfg    *config.Config
	log    logger.Logger
	bot    *bot.TelegramBot
	ctx    context.Context
	cancel context.CancelFunc
	db     *gorm.DB
}

func NewApp(ctx context.Context, cfg *config.Config, log logger.Logger) (*App, error) {
	ctx, cancel := context.WithCancel(ctx)

	// init middleware
	middlewareManager := middleware.NewManager(log)

	// init telebot
	b, err := bot.NewTelegramBot(cfg, log, middlewareManager)
	if err != nil {
		cancel()
		return nil, err
	}

	return &App{
		cfg:    cfg,
		log:    log,
		bot:    b,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (a *App) Start() error {
	a.log.Info("starting application")

	// 启动 bot
	if err := a.bot.Start(a.ctx); err != nil {
		return err
	}

	a.log.Info("application started successfully")
	return nil
}

func (a *App) Stop() error {
	a.log.Info("stopping application")

	// cancel ctx
	a.cancel()

	// stop bot
	a.bot.Stop()

	a.log.Info("application stopped successfully")
	return nil
}
