package config

import (
	"os"
)

type Config struct {
	Server       ServerConfig
	Database     DatabaseConfig
	Solana       SolanaConfig
	Notification NotificationConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Path string
}

type SolanaConfig struct {
	RPCURL string
	APIKey string
}

type NotificationConfig struct {
	TelegramBotToken string
	SMTP             SMTPConfig
	ProxyURL         string
}

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},
		Database: DatabaseConfig{
			Path: getEnv("DATABASE_PATH", "./data/monitor.db"),
		},
		Solana: SolanaConfig{
			RPCURL: getEnv("RPC_URL", "https://api.mainnet-beta.solana.com"),
			APIKey: getEnv("RPC_API_KEY", ""),
		},
		Notification: NotificationConfig{
			TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
			SMTP: SMTPConfig{
				Host:     getEnv("SMTP_HOST", "smtp.gmail.com"),
				Port:     getEnv("SMTP_PORT", "587"),
				Username: getEnv("SMTP_USER", ""),
				Password: getEnv("SMTP_PASS", ""),
				From:     getEnv("SMTP_FROM", "noreply@example.com"),
			},
			ProxyURL: getEnv("HTTP_PROXY", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}