package config

import (
	"github_spiders/pkg/config"
)

func init() {
	config.Add("redis", config.StrMap{
		"github": map[string]interface{}{
			"address":  config.Env("GITHUB_REDIS_ADDRESS", "127.0.0.1"),
			"password": config.Env("GITHUB_REDIS_PASSWORD", ""),
			"db":       config.Env("GITHUB_REDIS_DB", 0),
			"prefix":   config.Env("GITHUB_REDIS_PREFIX", ""),
		},
	})
}
