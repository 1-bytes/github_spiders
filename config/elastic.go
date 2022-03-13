package config

import (
	"github.com/olivere/elastic/v7"
	"github_spiders/pkg/config"
)

var ElasticOptions []elastic.ClientOptionFunc

func init() {
	config.Add("elastic", config.StrMap{
		"github": map[string]interface{}{
			"host":     config.Env("GITHUB_ELASTIC_HOST", "127.0.0.1:9200"),
			"username": config.Env("GITHUB_ELASTIC_USERNAME", ""),
			"password": config.Env("GITHUB_ELASTIC_PASSWORD", ""),
		},
	})
}
