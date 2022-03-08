package main

import (
	"github_spiders/bootstrap"
	"github_spiders/spiders/types"
	"time"
)

func init() {
	bootstrap.Setup()
}

// main 程序入口.
func main() {
	reposByUserC, reposByUserQ := bootstrap.NewCollector()
	usersByRepoC, usersByRepoQ := bootstrap.NewCollector()

	collector := types.GitHubCollector{
		ReposByUserC: reposByUserC,
		UsersByRepoC: usersByRepoC,
	}
	bootstrap.SetupCallback(collector)
	usersByRepoQ.AddURL("https://api.github.com/repos/1-bytes/GoBlog/stargazers")
	reposByUserQ.AddURL("https://api.github.com/users/1-bytes/starred")
	for {
		usersByRepoQ.Run(usersByRepoC)
		reposByUserQ.Run(reposByUserC)
		bootstrap.Waits(reposByUserC, usersByRepoC)

		time.Sleep(2 * time.Second)
		if usersByRepoQ.IsEmpty() && reposByUserQ.IsEmpty() {
			break
		}
	}

}
