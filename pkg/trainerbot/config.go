package trainerbot

import "english_trainer/internal/store"

type Config struct {
	LogLevel string
	Store    *store.Config
}

func NewConfig() *Config {
	return &Config{
		LogLevel: "debug",
		Store:    store.NewConfig(),
	}
}
