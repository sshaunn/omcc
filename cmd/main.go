package main

import (
	"context"
	"github.com/spf13/viper"
	"ohmycontrolcenter.tech/omcc/internal/app"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/config"
	"ohmycontrolcenter.tech/omcc/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 初始化日志
	log := logger.NewLogger()
	defer func(log logger.Logger) {
		err := log.Sync()
		if err != nil {

		}
	}(log)

	cfg, err := config.NewConfig("configs")
	log.Info("starting application",
		logger.String("config_file", viper.ConfigFileUsed()),
	)
	if err != nil {
		log.Fatal("failed to load config",
			logger.Error(err),
			logger.String("config_path", "configs/config.dev.yaml"),
		)
	}

	// create ctx
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create app instance
	application, err := app.NewApp(ctx, cfg, log)
	if err != nil {
		log.Fatal("failed to create application",
			logger.Error(err),
		)
	}

	// start app
	if err := application.Start(); err != nil {
		log.Fatal("failed to start application",
			logger.Error(err),
		)
	}

	// stop app
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Info("shutting down server...")

	if err := application.Stop(); err != nil {
		log.Error("exception during shutdown",
			logger.Error(err),
		)
	}
}
