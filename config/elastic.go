package config

import (
	"github_spiders/pkg/config"
)

func init() {
	config.Add("elastic", config.StrMap{
		"github": map[string]interface{}{
			"host":     config.Env("GITHUB_ELASTIC_HOST", "127.0.0.1:9200"),
			"username": config.Env("GITHUB_ELASTIC_USERNAME", ""),
			"password": config.Env("GITHUB_ELASTIC_PASSWORD", ""),
		},
	})
}
