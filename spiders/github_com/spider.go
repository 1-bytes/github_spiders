package github_com

import (
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/proxy"
	"github_spiders/pkg/config"
	"github_spiders/spiders/types"
	"log"
	"os"
)

// Spider github 爬虫.
type Spider struct {
	Debug     bool
	Domain    string
	UserAgent string
	Async     bool
	Users     types.GitHubUser
	LimitRule colly.LimitRule
	Socks5    string
	CacheDir  string
}

// Create 用于初始化 colly 框架的对象.
func (s Spider) Create() *colly.Collector {
	// 初始化爬虫框架
	cfg := []colly.CollectorOption{
		colly.AllowedDomains(s.Domain),
		colly.UserAgent(s.UserAgent),
		colly.Async(s.Async),
		colly.DetectCharset(),
		colly.CacheDir(s.CacheDir),
	}

	if s.Debug {
		output, _ := os.OpenFile(
			config.GetString("app.logger")+"/debug.log",
			os.O_CREATE|os.O_WRONLY,
			0666,
		)
		cfg = append(cfg, colly.Debugger(&debug.LogDebugger{
			Output: output,
			Prefix: "GitHub",
			Flag:   0,
		}))
	}
	c := colly.NewCollector(cfg...)
	_ = c.Limit(&s.LimitRule)

	// 设置代理
	if s.Socks5 != "" {
		rp, err := proxy.RoundRobinProxySwitcher(s.Socks5)
		if err != nil {
			log.Println("attempt to use Socks5 proxy failed.")
			panic(err)
		}
		c.SetProxyFunc(rp)
	}
	return c
}
