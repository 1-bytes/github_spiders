package main

import (
	"github_spiders/bootstrap"
	"github_spiders/spiders/types"
)

func init() {
	bootstrap.Setup()
}

// main 程序入口.
func main() {
	c := bootstrap.SetupCollector()
	bootstrap.SetupCallback(types.GitHubCollector{
		ReposByUserC: c,
		UsersByRepoC: c.Clone(),
	})
	err := c.Visit("https://api.github.com/users/1-bytes/starred?per_page=100&page=1")
	// err := c.Visit("https://api.github.com/users/1-bytes/starred?per_page=100&page=1")
	// err := c.Visit("https://api.github.com/user")
	if err != nil {
		panic(err)
	}
	c.Wait()
}
