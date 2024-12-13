package customer

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"ohmycontrolcenter.tech/omcc/internal/common"
	"ohmycontrolcenter.tech/omcc/internal/domain/model"
	"ohmycontrolcenter.tech/omcc/pkg/logger"
	"testing"
	"time"
)

type MockCustomerRepository struct {
	mock.Mock
}

func (m *MockCustomerRepository) DeleteCustomer(ctx context.Context, tx *gorm.DB, ids []string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockCustomerRepository) Create(ctx context.Context, tx *gorm.DB, customer *model.Customer) (*model.Customer, error) {
	args := m.Called(ctx, tx, customer)
	return args.Get(0).(*model.Customer), args.Error(1)
}

func (m *MockCustomerRepository) FindById(ctx context.Context, tx *gorm.DB, id string) (*model.Customer, error) {
	args := m.Called(ctx, tx, id)
	return args.Get(0).(*model.Customer), args.Error(1)
}

func (m *MockCustomerRepository) FindAllCustomers(ctx context.Context, tx *gorm.DB, page, limit int) ([]*model.CustomerWithBindings, int64, error) {
	args := m.Called(ctx, tx, page, limit)
	return args.Get(0).([]*model.CustomerWithBindings), args.Get(1).(int64), args.Error(2)
}

type MockCustomerTradingBindingRepository struct {
	mock.Mock
}

type MockCustomerSocialBindingRepository struct {
	mock.Mock
}

type MockTradingPlatformRepository struct {
	mock.Mock
}

type MockSocialPlatformRepository struct {
	mock.Mock
}

func (m *MockCustomerTradingBindingRepository) Create(ctx context.Context, tx *gorm.DB, binding *model.CustomerTradingBinding) (*model.CustomerTradingBinding, error) {
	args := m.Called(ctx, tx, binding)
	return args.Get(0).(*model.CustomerTradingBinding), args.Error(1)
}

func (m *MockCustomerTradingBindingRepository) CheckMemberStatus(ctx context.Context, tx *gorm.DB, uid string) (common.MemberStatus, error) {
	args := m.Called(ctx, tx, uid)
	return args.Get(0).(common.MemberStatus), args.Error(1)
}

func (m *MockCustomerTradingBindingRepository) FindTradingBindingByUid(ctx context.Context, tx *gorm.DB, uid string) (*model.CustomerInfoResponse, error) {
	args := m.Called(ctx, tx, uid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CustomerInfoResponse), args.Error(1)
}

func TestCustomerService_GetAllCustomers(t *testing.T) {
	// Create test time
	now := time.Now()

	// Create test data
	testCustomers := []*model.CustomerWithBindings{
		{
			CustomerId:          "1",
			CustomerUsername:    "test1",
			CustomerCreatedAt:   now,
			SocialUserId:        "social1",
			SocialUsername:      "social_test1",
			SocialFirstname:     "John",
			SocialLastname:      "Doe",
			SocialIsActive:      true,
			SocialStatus:        "normal",
			SocialMemberStatus:  "member",
			SocialCreatedAt:     now,
			TradingUid:          "trading1",
			TradingRegisterTime: now.Format(time.RFC3339),
			TradingCreatedAt:    now,
		},
		// Add more test data as needed
	}

	tests := []struct {
		name          string
		page          int
		limit         int
		setupMocks    func(*MockCustomerRepository, *MockCustomerSocialBindingRepository, *MockCustomerTradingBindingRepository, *MockTradingPlatformRepository, *MockSocialPlatformRepository)
		expectedTotal int64
		expectError   bool
		expectedCount int
	}{
		{
			name:  "success case",
			page:  1,
			limit: 10,
			setupMocks: func(customerRepo *MockCustomerRepository, socialRepo *MockCustomerSocialBindingRepository, tradingRepo *MockCustomerTradingBindingRepository, tradingPlatformRepo *MockTradingPlatformRepository, socialPlatformRepo *MockSocialPlatformRepository) {
				customerRepo.On("FindAllCustomers", mock.Anything, mock.AnythingOfType("*gorm.DB"), 1, 10).
					Return(testCustomers, int64(1), nil)
			},
			expectedTotal: 1,
			expectError:   false,
			expectedCount: 1,
		},
		{
			name:  "invalid page number",
			page:  -1,
			limit: 10,
			setupMocks: func(customerRepo *MockCustomerRepository, socialRepo *MockCustomerSocialBindingRepository, tradingRepo *MockCustomerTradingBindingRepository, tradingPlatformRepo *MockTradingPlatformRepository, socialPlatformRepo *MockSocialPlatformRepository) {
				customerRepo.On("FindAllCustomers", mock.Anything, mock.AnythingOfType("*gorm.DB"), 1, 10).
					Return(testCustomers, int64(1), nil)
			},
			expectedTotal: 1,
			expectError:   false,
			expectedCount: 1,
		},
		{
			name:  "invalid limit",
			page:  1,
			limit: 200,
			setupMocks: func(customerRepo *MockCustomerRepository, socialRepo *MockCustomerSocialBindingRepository, tradingRepo *MockCustomerTradingBindingRepository, tradingPlatformRepo *MockTradingPlatformRepository, socialPlatformRepo *MockSocialPlatformRepository) {
				customerRepo.On("FindAllCustomers", mock.Anything, mock.AnythingOfType("*gorm.DB"), 1, 10).
					Return(testCustomers, int64(1), nil)
			},
			expectedTotal: 1,
			expectError:   false,
			expectedCount: 1,
		},
		{
			name:  "repository error",
			page:  1,
			limit: 10,
			setupMocks: func(customerRepo *MockCustomerRepository, socialRepo *MockCustomerSocialBindingRepository, tradingRepo *MockCustomerTradingBindingRepository, tradingPlatformRepo *MockTradingPlatformRepository, socialPlatformRepo *MockSocialPlatformRepository) {
				customerRepo.On("FindAllCustomers", mock.Anything, mock.AnythingOfType("*gorm.DB"), 1, 10).
					Return([]*model.CustomerWithBindings(nil), int64(0), fmt.Errorf("db error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockCustomerRepo := new(MockCustomerRepository)
			mockCustomerSocialRepo := new(MockCustomerSocialBindingRepository)
			mockCustomerTradingRepo := new(MockCustomerTradingBindingRepository)
			mockTradingPlatformRepo := new(MockTradingPlatformRepository)
			mockSocialPlatformRepo := new(MockSocialPlatformRepository)
			tt.setupMocks(mockCustomerRepo, mockCustomerSocialRepo, mockCustomerTradingRepo, mockTradingPlatformRepo, mockSocialPlatformRepo)

			// Create service
			service := &CustomerServices{
				customerRepo: mockCustomerRepo,
				Log:          logger.NewLogger(),
			}

			// Execute test
			result, err := service.GetAllCustomers(context.Background(), tt.page, tt.limit)

			// Verify results
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedTotal, result.Total)
				assert.Equal(t, tt.expectedCount, len(result.Data))

				if tt.expectedCount > 0 {
					customer := result.Data[0]
					assert.Equal(t, testCustomers[0].CustomerId, customer.Customer.ID)
					assert.Equal(t, testCustomers[0].CustomerUsername, customer.Customer.Username)
					assert.Equal(t, testCustomers[0].SocialUsername, customer.SocialAccountInfo.Username)
					assert.Equal(t, testCustomers[0].TradingUid, customer.TradingAccountInfo.UID)
				}
			}

			// Verify mock expectations
			mockCustomerRepo.AssertExpectations(t)
		})
	}
}

func TestCustomerService_GetCustomerInfoByUid(t *testing.T) {
	tests := []struct {
		name        string
		uid         string
		setupMocks  func(*MockCustomerTradingBindingRepository)
		expectError bool
	}{
		{
			name: "success case",
			uid:  "test-uid",
			setupMocks: func(r *MockCustomerTradingBindingRepository) {
				r.On("FindTradingBindingByUid", mock.Anything, mock.Anything, "test-uid").
					Return(&model.CustomerInfoResponse{
						Customer: model.CustomerInfo{ID: "1"},
					}, nil)
			},
			expectError: false,
		},
		{
			name: "not found",
			uid:  "invalid-uid",
			setupMocks: func(r *MockCustomerTradingBindingRepository) {
				r.On("FindTradingBindingByUid", mock.Anything, mock.Anything, "invalid-uid").
					Return(nil, fmt.Errorf("not found"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTradingRepo := new(MockCustomerTradingBindingRepository)
			tt.setupMocks(mockTradingRepo)

			service := &CustomerServices{
				tradingBindingRepo: mockTradingRepo,
				Log:                logger.NewLogger(),
			}

			result, err := service.GetCustomerInfoByUid(context.Background(), tt.uid)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "1", result.Customer.ID)
			}

			mockTradingRepo.AssertExpectations(t)
		})
	}
}
