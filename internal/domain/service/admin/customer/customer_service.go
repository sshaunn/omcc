package customer

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"ohmycontrolcenter.tech/omcc/internal/common"
	"ohmycontrolcenter.tech/omcc/internal/domain/model"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/repository"
	"ohmycontrolcenter.tech/omcc/pkg/logger"
)

type CustomerServiceInterface interface {
	GetCustomerInfoByUid(ctx context.Context, uid string) (*model.CustomerInfoResponse, error)
	GetAllCustomers(ctx context.Context, page, limit int) (*model.PaginatedResponse[*model.CustomerInfoResponse], error)
	UpdateCustomerStatus(ctx context.Context, req *model.UpdateCustomerStatusRequest) error
	DeleteCustomer(ctx context.Context, req *model.DeleteCustomerRequest) ([]string, error)
}

// CustomerService struct
type CustomerService struct {
	customerRepo       repository.CustomerRepository
	socialBindingRepo  repository.CustomerSocialBindingRepository
	tradingBindingRepo repository.CustomerTradingBindingRepository
	tradingPlatform    repository.TradingPlatformRepository
	socialPlatform     repository.SocialPlatformRepository
	db                 *gorm.DB
	Log                logger.Logger
}

func NewCustomerService(db *gorm.DB, log logger.Logger) *CustomerService {
	return &CustomerService{
		customerRepo:       repository.NewCustomerRepository(db, log),
		socialBindingRepo:  repository.NewCustomerSocialRepository(db, log),
		tradingBindingRepo: repository.NewCustomerTradingRepository(db, log),
		tradingPlatform:    repository.NewTradingPlatformRepository(db, log),
		socialPlatform:     repository.NewSocialPlatformRepository(db, log),
		db:                 db,
		Log:                log,
	}
}

func (c *CustomerService) GetCustomerInfoByUid(ctx context.Context, uid string) (*model.CustomerInfoResponse, error) {
	customerInfo, err := c.tradingBindingRepo.FindTradingBindingByUid(ctx, c.db, uid)

	if err != nil {
		c.Log.Error("failed to get customer info",
			logger.String("uid", uid),
			logger.Error(err))
		return nil, err
	}
	return customerInfo, nil
}

func (c *CustomerService) GetAllCustomers(ctx context.Context, page, limit int) (*model.PaginatedResponse[*model.CustomerInfoResponse], error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	results, total, err := c.customerRepo.FindAllCustomers(ctx, nil, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get customers: %w", err)
	}
	var customerResponses []*model.CustomerInfoResponse
	for _, r := range results {
		customerResponses = append(customerResponses, &model.CustomerInfoResponse{
			Customer: model.CustomerInfo{
				ID:        r.CustomerId,
				Username:  r.CustomerUsername,
				CreatedAt: r.CustomerCreatedAt,
			},
			SocialAccountInfo: model.CustomerSocialInfo{
				UserID:       r.SocialUserId,
				Username:     r.SocialUsername,
				Firstname:    r.SocialFirstname,
				Lastname:     r.SocialLastname,
				IsActive:     r.SocialIsActive,
				Status:       r.SocialStatus,
				MemberStatus: r.SocialMemberStatus,
				CreatedAt:    r.SocialCreatedAt,
			},
			TradingAccountInfo: model.CustomerTradingInfo{
				UID:          r.TradingUid,
				RegisterTime: r.TradingRegisterTime,
				CreatedAt:    r.TradingCreatedAt,
			},
		})
	}
	return model.NewPaginatedResponse(customerResponses, total, page, limit), nil
}

func (c *CustomerService) UpdateCustomerStatus(ctx context.Context, req *model.UpdateCustomerStatusRequest) error {
	customer, err := c.customerRepo.FindById(ctx, c.db, req.CustomerId)
	if err != nil {
		return fmt.Errorf("customer not found: %w", err)
	}
	if customer == nil {
		return fmt.Errorf("customer not found")
	}
	err = c.socialBindingRepo.UpdateCustomerStatus(ctx, c.db, req.CustomerId, req.SocialId, *req.Status, common.GetMemberStatusFromString(*req.MemberStatus))
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	return nil
}

func (c *CustomerService) DeleteCustomer(ctx context.Context, req *model.DeleteCustomerRequest) ([]string, error) {
	ids, err := c.customerRepo.DeleteCustomer(ctx, c.db, req.IdList)
	if err != nil {
		return ids, fmt.Errorf("failed to delete customer ids=%v, error=%w", ids, err)
	}
	return ids, nil
}
