package bootstrap

import (
	"github_spiders/spiders/github_com/callbacks"
	"github_spiders/spiders/types"
)

// SetupCallback 初始化各个爬虫 Collector 的回调.
func SetupCallback(c types.GitHubCollector) {
	reposByUser := callbacks.ReposByUser{
		Colly: c,
	}
	reposByUser.Callbacks()

	usersByRepo := callbacks.UsersByRepo{
		Colly: c,
	}
	usersByRepo.Callbacks()
}
