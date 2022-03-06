package callbacks

import (
	"fmt"
	"github.com/gocolly/colly"
	"github_spiders/pkg/utils"
	"github_spiders/spiders/github_com/user"
	"log"
)

// UsersByRepo 列出已为存储库加注星标的人员.
// GitHub API docs url:
// https://docs.github.com/cn/rest/reference/activity#list-stargazers
type UsersByRepo struct {
	Colly *colly.Collector
	index int
}

// Callbacks 爬虫回调函数.
func (ur *UsersByRepo) Callbacks() {
	ur.index = 0
	auth := user.NewAuth()
	ur.Colly.OnRequest(func(r *colly.Request) {
		// GitHub's docs:
		// By default, all requests to https://api.github.com receive the v3 version of the REST API.
		// We encourage you to explicitly request this version via the Accept header.
		r.Headers.Add("Accept", "application/vnd.github.v3+json")
		r.Headers = auth.AddToken(r.Headers, ur.index)
	})

	ur.Colly.OnResponse(func(resp *colly.Response) {
		log.Printf("%d, %s", resp.StatusCode, resp.Body)
		body, err := utils.JsonUnmarshalBody(resp.Body)
		if err != nil {
			log.Printf("err: Failed to unmarshal the json: %s", err)
			return
		}

		for k, v := range body {
			fmt.Println(k, v)
		}
	})
}
