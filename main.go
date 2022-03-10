package main

import (
	"github_spiders/bootstrap"
	"github_spiders/pkg/collectors"
	"github_spiders/spiders/github_com/common"
	"github_spiders/spiders/types"
)

func init() {
	bootstrap.Setup()
}

// main 程序入口.
func main() {
	repo := collectors.GetInstance(types.TagsRepo)
	collectors.CloneToTag(types.TagsRepo, types.TagsUser)
	// user := collectors.GetInstance(types.TagsUser)
	bootstrap.SetupCallback()

	err := repo.Visit(common.CheckUrl("https://api.github.com/users/asciimoo/starred"))
	// err := user.Visit(common.CheckUrl("https://api.github.com/repos/go101/go101/stargazers?per_page=100&page=1"))
	if err != nil {
		panic(err)
	}
	repo.Wait()
	// user.Wait()
}
