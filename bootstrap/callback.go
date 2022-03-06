package bootstrap

import (
	"github.com/gocolly/colly"
	"github_spiders/spiders/github_com/callbacks"
)

// SetupCallback 初始化各个爬虫 Collector 的回调.
func SetupCallback(c *colly.Collector) {
	reposByUser := callbacks.ReposByUser{
		Colly: c,
	}
	reposByUser.Callbacks()

	// usersByRepo := callbacks.UsersByRepo{
	// 	Colly: c,
	// }
	// usersByRepo.Callbacks()
}
