package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
	"time"
)

type Config struct {
	App        AppConfig      `mapstructure:"app"`
	Telegram   TelegramConfig `mapstructure:"telegram"`
	Server     ServerConfig   `mapstructure:"server"`
	Bitget     BitgetConfig   `mapstructure:"exchange"`
	Database   DatabaseConfig `mapstructure:"database"`
	TimeFormat TimeFormatConfig
	// TODO redis
}

type AppConfig struct {
	Name        string `mapstructure:"name"`
	Environment string `mapstructure:"environment"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type TelegramConfig struct {
	Token      string        `mapstructure:"token" env:"TELEGRAM_BOT_TOKEN"`
	WebhookURL string        `mapstructure:"webhookUrl" env:"TELEGRAM_WEBHOOK_URL"`
	Timeout    time.Duration `mapstructure:"timeout"`
}

type BitgetConfig struct {
	ApiKey              string `mapstructure:"apiKey" env:"BITGET_API_KEY"`
	SecretKey           string `mapstructure:"secretKey" env:"BITGET_SECRET_KEY"`
	Passphrase          string `mapstructure:"passphrase" env:"BITGET_PASSPHRASE"`
	CustomerList        string `mapstructure:"customer_list"`
	CustomerTradeVolume string `mapstructure:"customer_trade_volume"`
}

type DatabaseConfig struct {
	Host               string        `mapstructure:"host"`
	Port               int           `mapstructure:"port"`
	User               string        `mapstructure:"user"`
	Password           string        `mapstructure:"password"`
	Database           string        `mapstructure:"database"`
	MaxIdleConnections int           `mapstructure:"max_idle_connections"`
	MaxOpenConnections int           `mapstructure:"max_open_connections"`
	MaxLifetime        time.Duration `mapstructure:"max_lifetime"`
}

type TimeFormatConfig struct {
	TimeFormat   string
	DateFormat   string
	TimeLocation *time.Location
}

func NewConfig(configPath string) (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}
	err := godotenv.Load(fmt.Sprintf("%s/.env.%s", configPath, env))
	if err != nil {
		return nil, fmt.Errorf("failed to load .env.%s", env)
	}
	viper.SetConfigFile(fmt.Sprintf("%s/config.%s.yaml", configPath, env))
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	loadSensitiveConfig()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// loadSensitiveConfig env variables for sensitive data
func loadSensitiveConfig() {
	// Telegram config
	viper.Set("telegram.token", os.Getenv("TELEGRAM_BOT_TOKEN"))
	viper.Set("telegram.webhookUrl", os.Getenv("TELEGRAM_WEBHOOK_URL"))
	viper.Set("telegram.botName", os.Getenv("TELEGRAM_BOT_NAME"))

	// Bitget config
	viper.Set("exchange.apiKey", os.Getenv("BITGET_API_KEY"))
	viper.Set("exchange.secretKey", os.Getenv("BITGET_SECRET_KEY"))
	viper.Set("exchange.passphrase", os.Getenv("BITGET_PASSPHRASE"))
}
