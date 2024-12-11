package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/config"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/database"
	"ohmycontrolcenter.tech/omcc/internal/middleware"
	"ohmycontrolcenter.tech/omcc/pkg/logger"
	"time"
)

type HTTPServer struct {
	engine *gin.Engine
	db     *gorm.DB
	cfg    *config.Config
	log    logger.Logger
	srv    *http.Server
}

func NewHTTPServer(cfg *config.Config, log logger.Logger) *HTTPServer {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	db, _ := database.NewMySqlClient(&cfg.Database, log)
	// using custom middleware
	engine.Use(gin.Recovery(), middleware.LoggerMiddleware(log))

	server := &HTTPServer{
		engine: engine,
		db:     db,
		cfg:    cfg,
		log:    log,
	}

	// route register
	server.registerRoutes()

	return server
}

func (s *HTTPServer) Start() error {
	s.srv = &http.Server{
		Addr:    ":" + s.cfg.Server.Port,
		Handler: s.engine,
	}

	return s.srv.ListenAndServe()
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *HTTPServer) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Format(time.DateTime),
	})
}
