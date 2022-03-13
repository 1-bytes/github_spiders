package config

import (
	"github_spiders/pkg/config"
)

func init() {
	config.Add("redis", config.StrMap{
		"github": map[string]interface{}{
			"host":     config.Env("GITHUB_REDIS_HOST", "127.0.0.1"),
			"port":     config.Env("GITHUB_REDIS_PORT", "6379"),
			"password": config.Env("GITHUB_REDIS_PASSWORD", ""),
			"db":       config.Env("GITHUB_REDIS_DB", 0),
			"prefix":   config.Env("GITHUB_REDIS_PREFIX", ""),
		},
	})
}
