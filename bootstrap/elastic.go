package bootstrap

import (
	"github.com/olivere/elastic/v7"
	configs "github_spiders/config"
	"github_spiders/pkg/config"
	"time"
)

// SetupElastic 初始化 Elastic.
func SetupElastic() {
	configs.ElasticOptions = []elastic.ClientOptionFunc{
		elastic.SetURL(config.GetString("elastic.github.host")),
		elastic.SetBasicAuth(
			config.GetString("elastic.github.username"),
			config.GetString("elastic.github.password"),
		),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(5 * time.Second),
	}
}
