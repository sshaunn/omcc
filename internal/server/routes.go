package server

import (
	handler "ohmycontrolcenter.tech/omcc/internal/api/admin"
	"ohmycontrolcenter.tech/omcc/internal/domain/bot"
	"ohmycontrolcenter.tech/omcc/internal/domain/service/admin/customer"
	"ohmycontrolcenter.tech/omcc/internal/domain/service/admin/message"
	"ohmycontrolcenter.tech/omcc/internal/domain/service/tasks"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/config"
)

func (s *HTTPServer) registerRoutes(cfg *config.Config, bot *bot.TelegramBot) {

	customerService := customer.NewCustomerService(s.db, s.log)
	customerHandler := handler.NewCustomerHandler(customerService, s.log)
	messagingService := message.NewSendingMessageService(cfg, bot.Bot, s.log)
	messagingHandler := handler.NewMessageHandler(messagingService, s.log)
	migratorService := tasks.NewMigrateService(cfg, s.log)
	migratorHandler := handler.NewMigrateCustomerHandler(migratorService, s.log)

	// API version
	v1 := s.engine.Group("/v1")
	{
		ad := v1.Group("/admin")
		{
			// admin management webui
			ad.GET("/customer", customerHandler.SearchByUID)
			ad.GET("/customers", customerHandler.GetAllCustomers)
			ad.PUT("/customer/update", customerHandler.UpdateCustomerStatus)
			ad.DELETE("/customer/delete", customerHandler.DeleteCustomer)

			// for admin to send message
			ad.POST("/customers/messaging", messagingHandler.SendMessage2MultipleCustomers)

			// for temp migrate
			ad.POST("/customers/migrate", migratorHandler.Migrate)
		}
	}

	{
		// 健康检查
		v1.GET("/health", s.healthCheck)

		// TODO: add other routes
	}
}
