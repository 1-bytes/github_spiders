package bootstrap

import (
	"github.com/olivere/elastic/v7"
	"github_spiders/pkg/config"
	pkgelastic "github_spiders/pkg/elastic"
	"time"
)

// SetupElastic 初始化 Elastic.
func SetupElastic() {
	pkgelastic.Options = []elastic.ClientOptionFunc{
		elastic.SetURL(config.GetString("elastic.github.host")),
		elastic.SetBasicAuth(
			config.GetString("elastic.github.username"),
			config.GetString("elastic.github.password"),
		),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(5 * time.Second),
	}
}
