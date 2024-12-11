package repository

import (
	"context"
	"fmt"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"ohmycontrolcenter.tech/omcc/internal/common"
	"ohmycontrolcenter.tech/omcc/internal/domain/model"
	"ohmycontrolcenter.tech/omcc/pkg/logger"
	"os"
	"testing"
	"time"
)

var (
	testDB   *gorm.DB
	pool     *dockertest.Pool
	resource *dockertest.Resource
)

func TestMain(m *testing.M) {
	// 创建 docker pool
	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// 设置超时时间
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// 启动 MySQL 容器
	resource, err = pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mysql",
		Tag:        "8.0",
		Env: []string{
			"MYSQL_ROOT_PASSWORD=secret",
			"MYSQL_DATABASE=testdb",
			"MYSQL_USER=test",
			"MYSQL_PASSWORD=test",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// 获取容器端口
	port := resource.GetPort("3306/tcp")

	// 重试直到数据库准备就绪
	if err := pool.Retry(func() error {
		var err error
		dsn := fmt.Sprintf("root:secret@(localhost:%s)/testdb?charset=utf8mb4&parseTime=True&loc=Local", port)
		testDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return err
		}

		sqlDB, err := testDB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Ping()
	}); err != nil {
		log.Printf("waiting for mysql database to be available...")
	}

	// 运行迁移
	err = testDB.AutoMigrate(
		&model.Customer{},
		&model.CustomerSocialBinding{},
		&model.CustomerTradingBinding{},
	)
	if err != nil {
		log.Fatalf("Could not migrate database: %s", err)
	}

	// 运行测试
	code := m.Run()

	// 清理
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestCustomerRepository_FindAllCustomers(t *testing.T) {
	tests := []struct {
		name        string
		page        int
		limit       int
		setupData   func(t *testing.T, db *gorm.DB) []string
		expectTotal int64
		expectError bool
		validate    func(t *testing.T, results []*model.CustomerWithBindings, customerIDs []string)
	}{
		{
			name:  "success with data",
			page:  1,
			limit: 10,
			setupData: func(t *testing.T, db *gorm.DB) []string {
				// 1. 创建平台
				socialPlatform := &model.SocialPlatform{
					Id:       1,
					Name:     "TestSocialPlatform",
					IsActive: true,
				}
				err := db.Create(socialPlatform).Error
				require.NoError(t, err)

				tradingPlatform := &model.TradingPlatform{
					Id:       1,
					Name:     "TestTradingPlatform",
					IsActive: true,
				}
				err = db.Create(tradingPlatform).Error
				require.NoError(t, err)

				// 2. 创建 customer
				customer := &model.Customer{
					Username: "test_user",
				}
				err = db.Create(customer).Error
				require.NoError(t, err)
				customerID := customer.Id // 获取生成的 UUID

				// 3. 创建 social binding
				socialBinding := &model.CustomerSocialBinding{
					CustomerID:   customerID,
					SocialID:     socialPlatform.Id,
					UserID:       "social1",
					Username:     "social_test1",
					MemberStatus: common.Member,
					Status:       "normal",
					IsActive:     true,
				}
				err = db.Create(socialBinding).Error
				require.NoError(t, err)

				// 4. 创建 trading binding
				tradingBinding := &model.CustomerTradingBinding{
					CustomerID:   customerID,
					TradingID:    tradingPlatform.Id,
					UID:          "trading1",
					RegisterTime: time.Now().Format(time.RFC3339),
				}
				err = db.Create(tradingBinding).Error
				require.NoError(t, err)

				return []string{customerID}
			},
			expectTotal: 1,
			expectError: false,
			validate: func(t *testing.T, results []*model.CustomerWithBindings, customerIDs []string) {
				require.Len(t, results, 1)
				assert.Equal(t, customerIDs[0], results[0].CustomerId)
				assert.Equal(t, "test_user", results[0].CustomerUsername)
				assert.Equal(t, "social1", results[0].SocialUserId)
				assert.Equal(t, "social_test1", results[0].SocialUsername)
				assert.Equal(t, string(common.Member), results[0].SocialMemberStatus)
				assert.Equal(t, "trading1", results[0].TradingUid)
			},
		},
		{
			name:  "empty result",
			page:  1,
			limit: 10,
			setupData: func(t *testing.T, db *gorm.DB) []string {
				// 返回空的 ID 列表
				return []string{}
			},
			expectTotal: 0,
			expectError: false,
			validate: func(t *testing.T, results []*model.CustomerWithBindings, customerIDs []string) {
				assert.Empty(t, results)
			},
		},
		{
			name:  "pagination",
			page:  2,
			limit: 1,
			setupData: func(t *testing.T, db *gorm.DB) []string {
				// 1. 创建所需的平台
				// 1.1 创建 social platform
				socialPlatform := &model.SocialPlatform{
					Id:       1,
					Name:     "TestSocialPlatform",
					IsActive: true,
				}
				err := db.Create(socialPlatform).Error
				require.NoError(t, err)

				// 1.2 创建 trading platform
				tradingPlatform := &model.TradingPlatform{
					Id:       1,
					Name:     "TestTradingPlatform",
					IsActive: true,
				}
				err = db.Create(tradingPlatform).Error
				require.NoError(t, err)

				// 2. 创建 customers 并收集 ID
				var customerIDs []string
				usernames := []string{"user1", "user2"}

				for _, username := range usernames {
					customer := &model.Customer{
						Username: username,
					}
					err := db.Create(customer).Error
					require.NoError(t, err)
					customerIDs = append(customerIDs, customer.Id)
					t.Logf("Created customer with ID: %s, Username: %s", customer.Id, customer.Username)
				}

				// 3. 创建 social bindings
				for i, customerID := range customerIDs {
					socialBinding := &model.CustomerSocialBinding{
						CustomerID:   customerID,
						SocialID:     socialPlatform.Id,
						UserID:       "social-" + customerID,
						Username:     "social_" + usernames[i],
						MemberStatus: "member",
						Status:       "normal",
						IsActive:     true,
					}
					err := db.Create(socialBinding).Error
					require.NoError(t, err)
				}

				// 4. 创建 trading bindings
				for _, customerID := range customerIDs {
					tradingBinding := &model.CustomerTradingBinding{
						CustomerID:   customerID,
						TradingID:    tradingPlatform.Id,
						UID:          "trading-" + customerID,
						RegisterTime: time.Now().Format(time.RFC3339),
					}
					err := db.Create(tradingBinding).Error
					require.NoError(t, err)
				}

				return customerIDs
			},
			expectTotal: 2,
			expectError: false,
			validate: func(t *testing.T, results []*model.CustomerWithBindings, customerIDs []string) {
				require.Len(t, results, 1)
				// 第二页，每页一条记录，按创建时间倒序，应该是第一条记录的ID
				assert.Equal(t, customerIDs[0], results[0].CustomerId,
					"Expected first customer ID (due to DESC order)")
				assert.Equal(t, "user1", results[0].CustomerUsername,
					"Expected first customer username")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 清理数据库
			cleanup(t, testDB)

			// 设置测试数据并获取生成的 IDs
			customerIDs := tt.setupData(t, testDB)

			// 创建 repository
			repo := NewCustomerRepository(testDB, logger.NewLogger())

			// 执行测试
			results, total, err := repo.FindAllCustomers(context.Background(), nil, tt.page, tt.limit)

			// 记录一下返回的结果用于调试
			if results != nil && len(results) > 0 {
				t.Logf("Results: ID=%s, Username=%s",
					results[0].CustomerId,
					results[0].CustomerUsername)
			}

			// 验证结果
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, results)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectTotal, total)
				tt.validate(t, results, customerIDs)
			}
		})
	}
}

func cleanup(t *testing.T, db *gorm.DB) {
	// 禁用外键检查
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")

	// 清理表（注意顺序）
	tables := []string{
		"customer_social_bindings",
		"customer_trading_bindings",
		"customers",
		"social_platforms",
		"trading_platforms",
	}
	for _, table := range tables {
		err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table)).Error
		require.NoError(t, err)
	}

	// 重新启用外键检查
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")
}
