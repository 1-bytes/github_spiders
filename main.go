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
	repo := collectors.GetInstance(types.TagsRepo)
	collectors.CloneToTag(types.TagsRepo, types.TagsUser)
	user := collectors.GetInstance(types.TagsUser)
	bootstrap.SetupCallback()

	// err := c.Visit("https://api.github.com/repos/1-bytes/GoBlog/stargazers?per_page=100&page=1")
	err := repo.Visit("https://api.github.com/users/1-bytes/starred")
	if err != nil {
		panic(err)
	}
	repo.Wait()
	user.Wait()
}
