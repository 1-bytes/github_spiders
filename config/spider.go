package config

import (
	"github_spiders/pkg/config"
)

func init() {
	config.Add("spiders", config.StrMap{
		"github": map[string]interface{}{
			"domain":      config.Env("GITHUB_SPIDER_DOMAIN", ""),
			"async":       config.Env("GITHUB_SPIDER_ASYNC", false),
			"user_agent":  config.Env("GITHUB_SPIDER_USER_AGENT", ""),
			"parallelism": config.Env("GITHUB_SPIDER_PARALLELISM", 1),
			"socks5":      config.Env("GITHUB_SPIDER_SOCKS5", ""),
			"users":       config.Env("GITHUB_SPIDER_USERS", "aaa,bbb,ccc"),
		},
	})
}
