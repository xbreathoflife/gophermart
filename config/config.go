package config

import (
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	Address        string `env:"RUN_ADDRESS"`
	ConnString     string `env:"DATABASE_URI"`
	ServiceAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func Init() Config {
	cfg := Config{
		Address:        "localhost:8080",
		ConnString:     "",
		ServiceAddress: "",
	}
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}
