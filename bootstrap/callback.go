package bootstrap

import "github_spiders/spiders/github_com/callbacks"

// SetupCallback 初始化各个爬虫 Collector 的回调.
func SetupCallback() {
	reposByUser := callbacks.ReposByUser{}
	reposByUser.Callbacks()
	usersByRepo := callbacks.UsersByRepo{}
	usersByRepo.Callbacks()
}
