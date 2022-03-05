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
	err := c.Visit("https://api.github.com/user/starred")
	// err := c.Visit("https://api.github.com/user")
	if err != nil {
		panic(err)
	}
	c.Wait()
}
