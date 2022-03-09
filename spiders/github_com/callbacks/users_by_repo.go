package callbacks

import (
	"github.com/gocolly/colly/v2"
	"github_spiders/pkg/collectors"
	"github_spiders/pkg/utils"
	"github_spiders/spiders/github_com/common"
	"github_spiders/spiders/github_com/user"
	"github_spiders/spiders/types"
	"log"
)

// UsersByRepo 列出已为存储库加注星标的人员.
// GitHub API docs url:
// https://docs.github.com/cn/rest/reference/activity#list-stargazers
type UsersByRepo struct {
	index int
}

// Callbacks 爬虫回调函数.
func (ur *UsersByRepo) Callbacks() {
	ur.index = 0
	auth := user.NewAuth()
	collector := collectors.GetInstance(types.TagsUser)
	collector.OnRequest(func(r *colly.Request) {
		// GitHub's docs:
		// By default, all requests to https://api.github.com receive the v3 version of the REST API.
		// We encourage you to explicitly request this version via the Accept header.
		r.Headers.Add("Accept", "application/vnd.github.v3+json")
		r.Headers = auth.AddToken(r.Headers, ur.index)
	})

	collector.OnResponse(func(resp *colly.Response) {
		users, err := utils.JsonUnmarshalBody(resp.Body)
		if err != nil {
			log.Printf("err: Failed to unmarshal the json: %s", err)
			return
		}

		userLen := len(users)
		if userLen == 0 {
			return
		}

		for _, u := range users {
			userName := u["login"]
			starViewUrl := u["url"].(string) + "/starred"
			log.Printf("【New User】 Name:%s, URL:%s",
				userName, starViewUrl)
			_ = collectors.GetInstance(types.TagsRepo).Visit(starViewUrl)
		}

		url := common.GetNextPageUrl(resp.Request, userLen)
		if url == "" {
			return
		}
		_ = collector.Visit(url)
	})
}
