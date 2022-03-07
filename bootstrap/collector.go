package bootstrap

import (
	"github.com/gocolly/colly"
	"github_spiders/pkg/config"
	"github_spiders/spiders/github_com"
)

var collector *colly.Collector

// SetupCollector 初始化爬虫收集器.
func SetupCollector() {
	if collector != nil {
		return
	}
	var (
		domain      = config.GetString("spiders.github.domain", "api.github.com")
		userAgent   = config.GetString("spiders.github.user_agent", "")
		parallelism = config.GetInt("parallelism", 3)
		socks5      = config.GetString("spiders.github.socks5", nil)
	)

	cfg := github_com.Spider{
		Debug:     config.GetBool("app.debug", true),
		Async:     config.GetBool("app.async", true),
		Domain:    domain,
		UserAgent: userAgent,
		LimitRule: colly.LimitRule{
			DomainRegexp: domain,
			Parallelism:  parallelism,
		},
		Socks5: socks5,
	}
	collector = cfg.Create()
}

// GetCollector 获取 colly.Collector.
func GetCollector() *colly.Collector {
	if collector == nil {
		SetupCollector()
	}
	return collector
}
