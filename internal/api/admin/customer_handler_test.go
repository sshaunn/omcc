package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"ohmycontrolcenter.tech/omcc/internal/domain/model"
	"ohmycontrolcenter.tech/omcc/pkg/logger"
	"testing"
	"time"
)

type MockCustomerService struct {
	mock.Mock
}

// 确保 MockCustomerService 实现了所有接口方法
func (m *MockCustomerService) GetCustomerInfoByUid(ctx context.Context, uid string) (*model.CustomerInfoResponse, error) {
	args := m.Called(ctx, uid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CustomerInfoResponse), args.Error(1)
}

func (m *MockCustomerService) GetAllCustomers(ctx context.Context, page, limit int) (*model.PaginatedResponse[*model.CustomerInfoResponse], error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PaginatedResponse[*model.CustomerInfoResponse]), args.Error(1)
}

func setupTestRouter(mockService *MockCustomerService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := NewCustomerHandler(mockService, logger.NewLogger())

	// 匹配实际的路由结构
	v1 := r.Group("/v1")
	{
		admin := v1.Group("/admin")
		{
			admin.GET("/customer", handler.SearchByUID)
			admin.GET("/customers", handler.GetAllCustomers)
		}
	}

	return r
}

func TestCustomerHandler_SearchByUID(t *testing.T) {
	// Create test data
	now := time.Now()
	testCustomer := &model.CustomerInfoResponse{
		Customer: model.CustomerInfo{
			ID:        "1",
			Username:  "test1",
			CreatedAt: now,
		},
		SocialAccountInfo: model.CustomerSocialInfo{
			UserID:    "social1",
			Username:  "social_test1",
			CreatedAt: now,
		},
		TradingAccountInfo: model.CustomerTradingInfo{
			UID:          "trading1",
			RegisterTime: now.Format(time.RFC3339),
			CreatedAt:    now,
		},
	}

	tests := []struct {
		name           string
		uid            string
		setupMock      func(*MockCustomerService)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "missing uid",
			uid:  "",
			setupMock: func(m *MockCustomerService) {
				// No mock setup needed
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: gin.H{
				"error": "uid is required",
			},
		},
		{
			name: "successful search",
			uid:  "test-uid",
			setupMock: func(m *MockCustomerService) {
				m.On("GetCustomerInfoByUid", mock.Anything, "test-uid").
					Return(testCustomer, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   testCustomer,
		},
		{
			name: "service error",
			uid:  "error-uid",
			setupMock: func(m *MockCustomerService) {
				m.On("GetCustomerInfoByUid", mock.Anything, "error-uid").
					Return(nil, fmt.Errorf("internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: gin.H{
				"error": "internal error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCustomerService)
			tt.setupMock(mockService)

			router := setupTestRouter(mockService)

			url := "/v1/admin/customer"
			if tt.uid != "" {
				url = fmt.Sprintf("%s?uid=%s", url, tt.uid)
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", url, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Compare response with expected body
			expectedJSON, _ := json.Marshal(tt.expectedBody)
			actualJSON, _ := json.Marshal(response)
			assert.JSONEq(t, string(expectedJSON), string(actualJSON))

			mockService.AssertExpectations(t)
		})
	}
}

func TestCustomerHandler_GetAllCustomers(t *testing.T) {
	// Create test data
	now := time.Now()
	testCustomers := &model.PaginatedResponse[*model.CustomerInfoResponse]{
		Data: []*model.CustomerInfoResponse{
			{
				Customer: model.CustomerInfo{
					ID:        "1",
					Username:  "test1",
					CreatedAt: now,
				},
			},
		},
		Total: 1,
		Page:  1,
		Limit: 10,
	}

	tests := []struct {
		name           string
		page           string
		limit          string
		setupMock      func(*MockCustomerService)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:  "default parameters",
			page:  "",
			limit: "",
			setupMock: func(m *MockCustomerService) {
				m.On("GetAllCustomers", mock.Anything, 1, 10).
					Return(testCustomers, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   testCustomers,
		},
		{
			name:  "custom parameters",
			page:  "2",
			limit: "20",
			setupMock: func(m *MockCustomerService) {
				m.On("GetAllCustomers", mock.Anything, 2, 20).
					Return(testCustomers, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   testCustomers,
		},
		{
			name:  "invalid page",
			page:  "invalid",
			limit: "10",
			setupMock: func(m *MockCustomerService) {
				// No mock setup needed
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: gin.H{
				"error": "page is required",
			},
		},
		{
			name:  "invalid limit",
			page:  "1",
			limit: "invalid",
			setupMock: func(m *MockCustomerService) {
				// No mock setup needed
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: gin.H{
				"error": "limit is required",
			},
		},
		{
			name:  "service error",
			page:  "1",
			limit: "10",
			setupMock: func(m *MockCustomerService) {
				m.On("GetAllCustomers", mock.Anything, 1, 10).
					Return(nil, fmt.Errorf("internal error"))
			},
			expectedStatus: http.StatusInternalServerError, // Your current implementation always returns 200
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockCustomerService)
			tt.setupMock(mockService)

			router := setupTestRouter(mockService)

			url := "/v1/admin/customers"
			if tt.page != "" || tt.limit != "" {
				url = fmt.Sprintf("%s?page=%s&limit=%s", url, tt.page, tt.limit)
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", url, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var response interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				expectedJSON, _ := json.Marshal(tt.expectedBody)
				actualJSON, _ := json.Marshal(response)
				assert.JSONEq(t, string(expectedJSON), string(actualJSON))
			}

			mockService.AssertExpectations(t)
		})
	}
}
