package config

import (
	"testValidate/internal/erro"
)

type Config struct {
	Port int
}

type ConfigOption func(*Config) error

func WithPort(port int) ConfigOption {
	return func(cfg *Config) error {
		if port < 0 {
			return erro.ErrorPort
		}
		cfg.Port = port
		return nil
	}
}

func NewConfig(options ...ConfigOption) (*Config, error) {
	cfg := &Config{
		Port: 8080, // Default port
	}

	for _, option := range options {
		err := option(cfg)
		if err != nil {
			return nil, err
		}
	}

	return cfg, nil
}
