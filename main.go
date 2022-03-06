package main

import (
	"github_spiders/bootstrap"
)

func init() {
	bootstrap.Setup()
}

// main 程序入口.
func main() {
	c := bootstrap.SetupCollector()
	bootstrap.SetupCallback(c)
	err := c.Visit("https://api.github.com/users/1-bytes/starred?per_page=100&page=1")
	// err := c.Visit("https://api.github.com/user")
	if err != nil {
		panic(err)
	}
	c.Wait()
}
