package bootstrap

import (
	"github.com/gocolly/colly/v2"
	"github_spiders/pkg/collectors"
	"github_spiders/pkg/config"
	"github_spiders/spiders/github_com"
)

// SetupCollector 初始化 Collector.
func SetupCollector() {
	var (
		domain      = config.GetString("spiders.github.domain", "api.github.com")
		userAgent   = config.GetString("spiders.github.user_agent", "")
		parallelism = config.GetInt("spiders.github.parallelism", 3)
		socks5      = config.GetString("spiders.github.socks5", nil)
		cacheDir    = config.GetString("spiders.github.cache_dir", "./runtime/cache")
	)

	// 这里只初始化配置 实例惰性加载
	collectors.SetConfig(&github_com.Spider{
		Debug:     config.GetBool("app.debug", true),
		Async:     config.GetBool("app.async", true),
		Domain:    domain,
		UserAgent: userAgent,
		LimitRule: colly.LimitRule{
			DomainRegexp: domain,
			Parallelism:  parallelism,
		},
		Socks5:   socks5,
		CacheDir: cacheDir,
	})
}
