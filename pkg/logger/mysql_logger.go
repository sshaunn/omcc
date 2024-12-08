package logger

import (
	"context"
	"fmt"
	gormLog "gorm.io/gorm/logger"
	"time"
)

// GormLogger implements GORM's logger interface
type GormLogger struct {
	gormLog.Interface
	log     Logger
	SlowSQL time.Duration
}

// NewGormLogger creates a new GORM logger instance
func NewGormLogger(log Logger) gormLog.Interface {
	return &GormLogger{
		log:     log,
		SlowSQL: time.Second, // Threshold for slow SQL logging
	}
}

// LogMode implements logger.Interface
func (l *GormLogger) LogMode(level gormLog.LogLevel) gormLog.Interface {
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
			String("sql", sql),
			Int64("rows", rows),
			Duration("elapsed", elapsed),
		)
	}

	// Log SQL errors
	if err != nil {
		l.log.Error("sql error",
			String("sql", sql),
			Error(err),
			Duration("elapsed", elapsed),
		)
		return
	}

	// Debug level for normal queries
	l.log.Debug("sql trace",
		String("sql", sql),
		Int64("rows", rows),
		Duration("elapsed", elapsed),
	)
}
