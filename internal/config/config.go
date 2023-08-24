package config

import (
	"fmt"

	"github.com/vrischmann/envconfig"
)

type Config struct {
	App struct {
		Port        int
		ServiceName string
	}

	DB struct {
		Username string
		Password string
		Port     string
		Host     string
		DBname   string
	}
}

func InitConfig(prefix string) (*Config, error) {
	conf := &Config{}

	if len(prefix) > 0 {
		if err := envconfig.InitWithPrefix(conf, prefix); err != nil {
			return nil, fmt.Errorf("init config error: %w", err)
		}
	} else {
		if err := envconfig.Init(conf); err != nil {
			return nil, fmt.Errorf("init config error: %w", err)
		}
	}

	return conf, nil
}
