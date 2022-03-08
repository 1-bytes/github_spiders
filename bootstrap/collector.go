package bootstrap

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	queued "github_spiders/pkg/colly/queue"
	"github_spiders/pkg/colly/redis"
	"github_spiders/pkg/config"
	"github_spiders/spiders/github_com"
)

var (
	// collector *colly.Collector
	cfg github_com.Spider
)

// SetupCollector 初始化爬虫收集器.
func SetupCollector() {
	var (
		domain      = config.GetString("spiders.github.domain", "api.github.com")
		userAgent   = config.GetString("spiders.github.user_agent", "")
		parallelism = config.GetInt("spiders.github.parallelism", 3)
		socks5      = config.GetString("spiders.github.socks5", nil)
		cacheDir    = config.GetString("spiders.github.cache_dir", "./runtime/cache")
	)

	cfg = github_com.Spider{
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
	}
}

// NewCollector 创建一个新的 Collector 并将其 Queue 一起返回.
func NewCollector() (*colly.Collector, *queue.Queue) {
	threads := config.GetInt("spiders.github.parallelism", 3)
	collector := cfg.Create()
	id := int(collector.ID)
	storage := redis.GetInstance(id).NewRedisStorage(collector, GetCollyRedisConfig())
	queued.GetInstance(id).NewQueue(threads, storage)
	q := queued.GetInstance(id).GetQueue()
	return collector, q
}

func Waits(c ...*colly.Collector) {
	for _, collector := range c {
		collector.Wait()
	}
}
