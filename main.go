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
	user := collectors.GetInstance(types.TagsUser)
	collectors.CloneToTag(types.TagsUser, types.TagsRepo)
	repo := collectors.GetInstance(types.TagsRepo)
	bootstrap.SetupCallback()

	err := user.Visit("https://api.github.com/repos/1-bytes/GoBlog/stargazers")
	if err != nil {
		panic(err)
	}
	user.Wait()
	repo.Wait()
}
