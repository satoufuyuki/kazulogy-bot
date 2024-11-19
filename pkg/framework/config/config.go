package config

type Config struct {
	DatabaseLogLevel string `env:"DATABASE_LOG_LEVEL" envDefault:"DEBUG"`
	ClientLogLevel   string `env:"CLIENT_LOG_LEVEL" envDefault:"DEBUG"`
}
