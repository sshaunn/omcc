package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ohmycontrolcenter.tech/omcc/internal/domain/service/admin/customer"
	"ohmycontrolcenter.tech/omcc/pkg/logger"
	"strconv"
)

type CustomerHandler struct {
	customerService customer.CustomerServiceInterface
	log             logger.Logger
}

func NewCustomerHandler(customerService customer.CustomerServiceInterface, log logger.Logger) *CustomerHandler {
	return &CustomerHandler{
		customerService: customerService,
		log:             log,
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