package main

import (
	"github.com/gocolly/colly/v2/queue"
	"github_spiders/bootstrap"
	"github_spiders/pkg/collectors"
	"github_spiders/pkg/queued"
	"github_spiders/spiders/github_com/callbacks"
	"sync"
)

func init() {
	bootstrap.Setup()
}

func main() {
	var wg sync.WaitGroup

	urls := map[string][]string{
		callbacks.TagRepo: {
			"https://api.github.com/users/asciimoo/starred",
			"https://api.github.com/users/LinuxMercedes/starred",
			"https://api.github.com/users/brson/starred",
			"https://api.github.com/users/avelino/starred",
			"https://api.github.com/users/Kikobeats/starred",
		},
		callbacks.TagUser: {
			"https://api.github.com/repos/golang/go/stargazers",
			"https://api.github.com/repos/laravel/laravel/stargazers",
			"https://api.github.com/repos/torvalds/linux/stargazers",
			"https://api.github.com/repos/atom/atom/stargazers",
			"https://api.github.com/repos/browserless/chrome/stargazers",
		},
	}

	bootstrap.SetupCallback()
	queues := make(map[string]*queue.Queue)
	for _, tag := range callbacks.Tags {
		queues[tag] = queued.GetInstance(tag)
		for _, url := range urls[tag] {
			_ = queues[tag].AddURL(url)
		}
	}

	for tag, q := range queues {
		wg.Add(1)
		go func(tag string, q *queue.Queue) {
			instance := collectors.GetInstance(tag)
			_ = q.Run(instance)
			wg.Done()
		}(tag, q)
	}
	wg.Wait()
}
