package app

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"ohmycontrolcenter.tech/omcc/internal/domain/bot"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/config"
	"ohmycontrolcenter.tech/omcc/internal/middleware"
	"ohmycontrolcenter.tech/omcc/internal/server"
	"ohmycontrolcenter.tech/omcc/pkg/logger"
)

type App struct {
	cfg        *config.Config
	log        logger.Logger
	bot        *bot.TelegramBot
	httpServer *server.HTTPServer
	ctx        context.Context
	cancel     context.CancelFunc
	db         *gorm.DB
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

	httpServer := server.NewHTTPServer(cfg, log)

	return &App{
		cfg:        cfg,
		log:        log,
		bot:        b,
		httpServer: httpServer,
		ctx:        ctx,
		cancel:     cancel,
	}, nil
}

func (a *App) Start() error {
	a.log.Info("starting application with telebot and httpserver")

	// start bot in goroutine
	go func() {
		err := a.bot.Start(a.ctx)
		if err != nil {
			a.log.Info("starting application failed")
		}
	}()

	go func() {
		a.log.Info(fmt.Sprintf("Started admin server listening on port=%s", a.cfg.Server.Port))
		err := a.httpServer.Start()
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				a.log.Error("http server failed to start",
					logger.Error(err),
				)
			}
		}
	}()

	a.log.Info("application started successfully")
	return nil
}

func (a *App) Stop() error {
	a.log.Info("stopping application")

	// cancel ctx
	a.cancel()

	// stop bot
	a.bot.Stop()
	if err := a.httpServer.Stop(a.ctx); err != nil {
		a.log.Error("failed to stop http server",
			logger.Error(err),
		)
	}

	a.log.Info("application stopped successfully")
	return nil
}
