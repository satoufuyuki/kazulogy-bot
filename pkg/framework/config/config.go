package config

type Config struct {
	DatabaseLogLevel string `env:"DATABASE_LOG_LEVEL" envDefault:"INFO"`
	ClientLogLevel   string `env:"CLIENT_LOG_LEVEL" envDefault:"INFO"`
	BotPrefix        string `env:"BOT_PREFIX"`
}
