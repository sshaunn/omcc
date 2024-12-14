package client

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"ohmycontrolcenter.tech/omcc/internal/infrastructure/config"
	"ohmycontrolcenter.tech/omcc/pkg/logger"
	"strconv"
	"time"
)

type BitgetClient struct {
	client      *fasthttp.Client
	log         logger.Logger
	cfg         *config.BitgetConfig
	rateLimiter *config.RateLimiter
}

type Options func(*BitgetClient)

func NewBitgetClient(cfg *config.BitgetConfig, log logger.Logger) *BitgetClient {
	return &BitgetClient{
		client: &fasthttp.Client{
			MaxConnsPerHost:     100,
			MaxIdleConnDuration: 30 * time.Second,
			ReadTimeout:         5 * time.Second,
			WriteTimeout:        5 * time.Second,
		},
		log:         log,
		cfg:         cfg,
		rateLimiter: config.GetBitgetRateLimiter(log),
	}
}

func (b *BitgetClient) Post(ctx context.Context, path string, body interface{}) ([]byte, error) {
	if err := b.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(b.cfg.BaseUrl + path)
	req.Header.SetMethod("POST")

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %w", err)
	}
	req.SetBody(jsonBody)
	startTime := time.Now()
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	signature := b.generateSignature(timestamp, "POST", path, "", string(jsonBody))

	b.setHeaders(req, timestamp, signature)
	err = b.client.Do(req, resp)
	if err != nil {
		b.log.Error("Error occurred while invoking bitget api",
			logger.String("path", path),
			logger.String("jsonBody", string(jsonBody)))
		return nil, err
	}

	b.log.Info("Successfully invoked bitget api",
		logger.String("path", path),
		logger.Int("status", resp.StatusCode()),
		logger.String("jsonBody", string(jsonBody)),
		logger.Duration("elapsedTime", time.Since(startTime)),
	)
	return resp.Body(), nil
}

func (b *BitgetClient) generateSignature(timestamp, method, path, query, body string) string {
	payload := fmt.Sprintf("%s%s%s%s%s",
		timestamp,
		method,
		path,
		query,
		body,
	)

	h := hmac.New(sha256.New, []byte(b.cfg.SecretKey))
	h.Write([]byte(payload))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (b *BitgetClient) setHeaders(req *fasthttp.Request, timestamp, signature string) {
	req.Header.Set("ACCESS-KEY", b.cfg.ApiKey)
	req.Header.Set("ACCESS-SIGN", signature)
	req.Header.Set("ACCESS-TIMESTAMP", timestamp)
	req.Header.Set("ACCESS-PASSPHRASE", b.cfg.Passphrase)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "locale=en-US")
	req.Header.Set("locale", "en-US")
}
