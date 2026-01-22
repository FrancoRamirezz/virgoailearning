package utils

import "backend/config"

var cfg *config.Config

func GetConfig() *config.Config {
	if cfg == nil {
		cfg = config.Load()
	}
	return cfg
}


