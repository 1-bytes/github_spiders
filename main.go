package main

import (
	"github_spiders/bootstrap"
	"github_spiders/pkg/colly/queue"
	"github_spiders/spiders/types"
)

func init() {
	bootstrap.Setup()
}

// main 程序入口.
func main() {
	reposByUserC := bootstrap.GetCollector()
	collector := types.GitHubCollector{
		ReposByUserC: reposByUserC,
		UsersByRepoC: reposByUserC.Clone(),
	}
	bootstrap.SetupEngine(collector)

	rbu := queue.GetInstance(int(collector.ReposByUserC.ID)).GetQueue()
	ubr := queue.GetInstance(int(collector.UsersByRepoC.ID)).GetQueue()
	rbu.AddURL("https://api.github.com/users/1-bytes/starred?per_page=100&page=1")
	ubr.AddURL("https://api.github.com/repos/1-bytes/GoBlog/stargazers?per_page=100&page=1")

	rbu.Run(collector.ReposByUserC)
	ubr.Run(collector.UsersByRepoC)
	reposByUserC.Wait()

	// err := c2.Visit("https://api.github.com/repos/1-bytes/GoBlog/stargazers?per_page=100&page=1")
	// err := reposByUserC.Visit("https://api.github.com/users/1-bytes/starred?per_page=100&page=1")
	// if err != nil {
	// 	panic(err)
	// }
	// reposByUserC.Wait()

	// time.Sleep(10 * time.Second)
}
