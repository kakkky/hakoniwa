package config

import "github.com/caarlos0/env/v11"

type Config struct {
	XdgDataHome  string `env:"XDG_DATA_HOME"`
	XdgStateHome string `env:"XDG_STATE_HOME"`
	GeminiAPIKey string `env:"GEMINI_API_KEY,required"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
