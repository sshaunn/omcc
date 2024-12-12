package server

import (
	handler "ohmycontrolcenter.tech/omcc/internal/api/admin"
	"ohmycontrolcenter.tech/omcc/internal/domain/service/admin/customer"
)

func (s *HTTPServer) registerRoutes() {

	customerService := customer.NewCustomerService(s.db, s.log)
	customerHandler := handler.NewCustomerHandler(customerService, s.log)

	// API version
	v1 := s.engine.Group("/v1")
	{
		ad := v1.Group("/admin")
		{
			ad.GET("/customer", customerHandler.SearchByUID)
			ad.GET("/customers", customerHandler.GetAllCustomers)
			ad.PUT("/customer/update", customerHandler.UpdateCustomerStatus)
		}
	}

	{
		// 健康检查
		v1.GET("/health", s.healthCheck)

		// TODO: add other routes
	}
}
