package config

import (
	"context"
	"golang.org/x/time/rate"
	"ohmycontrolcenter.tech/omcc/pkg/logger"
	"sync"
	"sync/atomic"
	"time"
)

type RateLimiter struct {
	limiter *rate.Limiter
	mu      sync.RWMutex
	stats   struct {
		totalRequests uint64
		errors        uint64
		lastMinute    struct {
			requests uint64
			errors   uint64
			time     time.Time
		}
	}
	log logger.Logger
}

var (
	globalBitgetLimiter *RateLimiter
	once                sync.Once
)

func GetBitgetRateLimiter(log logger.Logger) *RateLimiter {
	once.Do(func() {
		globalBitgetLimiter = &RateLimiter{
			limiter: rate.NewLimiter(rate.Limit(10), 10),
			log:     log,
		}
		globalBitgetLimiter.stats.lastMinute.time = time.Now()

		// 启动统计重置器
		go globalBitgetLimiter.resetMinuteStats()
	})
	return globalBitgetLimiter
}

func (r *RateLimiter) Do(ctx context.Context, fn func() error) error {
	r.mu.Lock()
	r.stats.totalRequests++
	r.mu.Unlock()
	err := r.limiter.Wait(ctx)
	if err != nil {
		r.mu.Lock()
		r.stats.errors++
		r.mu.Unlock()
		return err
	}

	if err := fn(); err != nil {
		r.mu.Lock()
		r.stats.errors++
		r.mu.Unlock()
		return err
	}
	return fn()
}

func (r *RateLimiter) GetStats() (total, errors, lastMinuteReq, lastMinuteErr uint64) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return atomic.LoadUint64(&r.stats.totalRequests),
		atomic.LoadUint64(&r.stats.errors),
		atomic.LoadUint64(&r.stats.lastMinute.requests),
		atomic.LoadUint64(&r.stats.lastMinute.errors)
}

func (r *RateLimiter) Wait(ctx context.Context) error {
	atomic.AddUint64(&r.stats.totalRequests, 1)
	atomic.AddUint64(&r.stats.lastMinute.requests, 1)

	err := r.limiter.Wait(ctx)
	if err != nil {
		atomic.AddUint64(&r.stats.errors, 1)
		atomic.AddUint64(&r.stats.lastMinute.errors, 1)
		return err
	}
	return nil
}

func (r *RateLimiter) resetMinuteStats() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		r.mu.Lock()
		atomic.StoreUint64(&r.stats.lastMinute.requests, 0)
		atomic.StoreUint64(&r.stats.lastMinute.errors, 0)
		r.stats.lastMinute.time = time.Now()
		r.mu.Unlock()

		// 记录每分钟的统计信息
		total, errors, lastMinReq, lastMinErr := r.GetStats()
		r.log.Info("Bitget API rate limiter stats",
			logger.UInt64("total_requests", total),
			logger.UInt64("total_errors", errors),
			logger.UInt64("last_minute_requests", lastMinReq),
			logger.UInt64("last_minute_errors", lastMinErr))
	}
}
