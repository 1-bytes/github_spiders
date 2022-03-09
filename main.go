package main

import (
	"github_spiders/bootstrap"
	"github_spiders/pkg/collectors"
	"github_spiders/spiders/types"
)

func init() {
	bootstrap.Setup()
}

// main 程序入口.
func main() {
	reposByUserC := collectors.GetInstance(types.TagsRepo)
	collectors.CloneToTag(types.TagsRepo, types.TagsUser)
	usersByRepoC := collectors.GetInstance(types.TagsUser)

	err := reposByUserC.Visit("https://api.github.com/users/1-bytes/starred?per_page=100&page=1")
	// _ = usersByRepoC.Visit("https://api.github.com/repos/1-bytes/GoBlog/stargazers?per_page=100&page=1")
	if err != nil {
		panic(err)
	}
	reposByUserC.Wait()
	usersByRepoC.Wait()
}
