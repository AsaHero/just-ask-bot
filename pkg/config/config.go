package config

import "os"

const (
	ProductionEnvironment  = "production"
	DevelopmentEnvironment = "development"
	LocalEnvironment       = "local"
)

// Config содержит конфигурационные параметры приложения
type Config struct {
	APP         string
	Environment string
	LogLevel    string

	Context struct {
		Timeout string
	}

	Server struct {
		Host         string
		Port         string
		ReadTimeout  string
		WriteTimeout string
		IdleTimeout  string
	}

	DB struct {
		Host     string
		Port     string
		Name     string
		User     string
		Password string
		SSLMode  string
	}

	Telegram struct {
		BotToken   string
		WebhookURL string
	}

	OpenAI struct {
		SecretKey string
	}
}

// New создает и возвращает новый экземпляр Config
func New() *Config {
	var config Config

	config.APP = getEnv("APP", "esoterica_api_server")
	config.Environment = getEnv("ENVIRONMENT", "develop")
	config.LogLevel = getEnv("LOG_LEVEL", "debug")
	config.Context.Timeout = getEnv("CONTEXT_TIMEOUT", "10m")

	// server configuration
	config.Server.Host = getEnv("SERVER_HOST", "localhost")
	config.Server.Port = getEnv("SERVER_PORT", ":8000")
	config.Server.ReadTimeout = getEnv("SERVER_READ_TIMEOUT", "10s")
	config.Server.WriteTimeout = getEnv("SERVER_WRITE_TIMEOUT", "10s")
	config.Server.IdleTimeout = getEnv("SERVER_IDLE_TIMEOUT", "120s")

	// initialization db
	config.DB.Host = getEnv("POSTGRES_HOST", "localhost")
	config.DB.Port = getEnv("POSTGRES_PORT", "5432")
	config.DB.Name = getEnv("POSTGRES_DATABASE", "just_ask")
	config.DB.User = getEnv("POSTGRES_USER", "postgres")
	config.DB.Password = getEnv("POSTGRES_PASSWORD", "postgres")
	config.DB.SSLMode = getEnv("POSTGRES_SSLMODE", "disable")

	config.Telegram.BotToken = getEnv("TELEGRAM_BOT_TOKEN", "7103040220:AAHLl21r2EzsR2hFltrI740Ix9JGOcAnzzs")
	config.Telegram.WebhookURL = getEnv("TELEGRAM_WEBHOOK_URL", "https://92f8ca69fcdf-3450089824392956192.ngrok-free.app/webhooks/telegram")
	config.OpenAI.SecretKey = getEnv("OPENAI_KEY", "")

	return &config
}

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultValue
}
