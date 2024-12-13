package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ohmycontrolcenter.tech/omcc/internal/domain/model"
	"ohmycontrolcenter.tech/omcc/internal/domain/service/admin/customer"
	"ohmycontrolcenter.tech/omcc/internal/domain/service/tasks"
	"ohmycontrolcenter.tech/omcc/pkg/logger"
	"strconv"
)

type CustomerHandler struct {
	customerService customer.CustomerServiceInterface
	migrateService  tasks.MigrateService
	log             logger.Logger
}

func NewCustomerHandler(customerService customer.CustomerServiceInterface, log logger.Logger) *CustomerHandler {
	return &CustomerHandler{
		customerService: customerService,
		log:             log,
	}
}

func NewMigrateCustomerHandler(migrateService *tasks.MigrateService, log logger.Logger) *CustomerHandler {
	return &CustomerHandler{
		migrateService: *migrateService,
		log:            log,
	}
}

func (h *CustomerHandler) SearchByUID(c *gin.Context) {
	uid := c.Query("uid")
	if uid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uid is required"})
		return
	}

	customerInfo, err := h.customerService.GetCustomerInfoByUid(c.Request.Context(), uid)
	if err != nil {
		h.log.Error("failed to get user info",
			logger.String("uid", uid),
			logger.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, customerInfo)
}

func (h *CustomerHandler) GetAllCustomers(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page is required"})
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit is required"})
		return
	}
	customersInfo, err := h.customerService.GetAllCustomers(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, customersInfo)
}

func (h *CustomerHandler) UpdateCustomerStatus(c *gin.Context) {
	var req model.UpdateCustomerStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("invalid request parameters",
			logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request parameters",
		})
		return
	}

	err := h.customerService.UpdateCustomerStatus(c.Request.Context(), &req)
	if err != nil {
		h.log.Error("failed to update customer status",
			logger.String("customer_id", req.CustomerId),
			logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "Success",
		"message": "Customer status updated successfully",
	})
}

func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	var req model.DeleteCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("invalid request parameters")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request parameters",
		})
		return
	}
	ids, err := h.customerService.DeleteCustomer(c.Request.Context(), &req)
	if err != nil {
		h.log.Error("failed to delete customer",
			logger.Any("failed_ids", ids))
		c.JSON(http.StatusInternalServerError, gin.H{
			"id_list": ids,
			"code":    "Failure",
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"id_list": ids,
		"code":    "Success",
		"message": "Customer deleted successfully",
	})
}

func (h *CustomerHandler) Migrate(c *gin.Context) {
	err := h.migrateService.Migrate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "Failure",
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    "Success",
		"message": "migrate successfully",
	})
}
