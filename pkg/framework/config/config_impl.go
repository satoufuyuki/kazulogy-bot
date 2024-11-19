package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

func NewConfig() (config Config, err error) {
	config = Config{}
	godotenv.Load()

	if err = env.Parse(&config); err != nil {
		return
	}

	return
}
