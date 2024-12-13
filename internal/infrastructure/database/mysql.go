package database

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/config"
	log "ohmycontrolcenter.tech/omcc/pkg/logger"

	"time"
)

// GormLogger implements GORM's logger interface
type GormLogger struct {
	logger  logger.Interface
	log     log.Logger
	SlowSQL time.Duration
}

// NewGormLogger creates a new GORM logger instance
func NewGormLogger(log log.Logger) logger.Interface {
	return &GormLogger{
		log:     log,
		SlowSQL: time.Second, // Threshold for slow SQL logging
	}
}

// NewPostgresClient Temporary using this for migrating db records from postgres to mysql
// TODO after success migrating then cleanup this
func NewPostgresClient(log log.Logger) (*gorm.DB, error) {
	// this is in .env.dev file temporarily
	dsn := ""
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres database: %w", err)
	}
	return db, nil
}

func NewMySqlClient(cfg *config.DatabaseConfig, log log.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: NewGormLogger(log),
		// Disable default transaction
		//SkipDefaultTransaction: true,

		// Prepare statement
		PrepareStmt: true,

		//// Name strategy
		//NamingStrategy: schema.NamingStrategy{
		//	SingularTable: true, // Use singular table names
		//},
		//
		//// Disable foreign key constraint when migrating
		//DisableForeignKeyConstraintWhenMigrating: true,

		// Query timeout
		QueryFields: true,

		// Create/Update timestamp tracking
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(cfg.MaxLifetime)

	return db, nil
}

func WithTransaction(db *gorm.DB, fn func(*gorm.DB) error) error {
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// LogMode implements logger.Interface
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

// Info implements logger.Interface
func (l *GormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	l.log.Info(fmt.Sprintf(msg, args...))
}

// Warn implements logger.Interface
func (l *GormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	l.log.Warn(fmt.Sprintf(msg, args...))
}

// Error implements logger.Interface
func (l *GormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	l.log.Error(fmt.Sprintf(msg, args...))
}

// Trace implements logger.Interface
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	// Log slow SQL queries
	if elapsed > l.SlowSQL {
		l.log.Warn("slow sql",
			log.String("sql", sql),
			log.Int64("rows", rows),
			log.Duration("elapsed", elapsed),
		)
	}

	// Log SQL errors
	if err != nil {
		l.log.Error("sql error",
			log.String("sql", sql),
			log.Error(err),
			log.Duration("elapsed", elapsed),
		)
		return
	}

	// Debug level for normal queries
	l.log.Debug("sql trace",
		log.String("sql", sql),
		log.Int64("rows", rows),
		log.Duration("elapsed", elapsed),
	)
}
