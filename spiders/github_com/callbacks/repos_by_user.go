package callbacks

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github_spiders/pkg/collectors"
	"github_spiders/pkg/utils"
	"github_spiders/spiders/github_com/common"
	"github_spiders/spiders/github_com/user"
	"github_spiders/spiders/types"
	"log"
)

// ReposByUser 列出用户已加星标的存储库.
// GitHub API docs url:
// https://docs.github.com/cn/rest/reference/activity#list-repositories-starred-by-a-user
type ReposByUser struct {
	index int
}

// Callbacks 爬虫回调函数.
func (ru *ReposByUser) Callbacks() {
	ru.index = 0
	auth := user.NewAuth()
	collector := collectors.GetInstance(types.TagsRepo)
	collector.OnRequest(func(r *colly.Request) {
		// GitHub's docs:
		// By default, all requests to https://api.github.com receive the v3 version of the REST API.
		// We encourage you to explicitly request this version via the Accept header.
		r.Headers.Add("Accept", "application/vnd.github.v3+json")
		r.Headers = auth.AddToken(r.Headers, ru.index)
	})

	collector.OnResponse(func(resp *colly.Response) {
		// if resp.StatusCode != http.StatusOK {
		// TODO:// 回头再处理这个问题，状态码为非200状态时，有可能是需要更换帐号的 token 了 ..
		// }
		repos, err := utils.JsonUnmarshalBody(resp.Body)
		if err != nil {
			log.Printf("err: Failed to unmarshal the json: %s", err)
			return
		}
		// 检查该用户有没有 star 的存储库，或者是 star 的存储库列表到了最终页
		repoLen := len(repos)
		if repoLen == 0 {
			return
		}

		for _, repo := range repos {
			repoName := repo["full_name"]
			starViewUrl := repo["stargazers_url"]
			starCount := repo["stargazers_count"]
			log.Printf("【New Repo】 Name:%s, StarCount:%v, URL:%s",
				repoName, starCount, starViewUrl)
			_ = collectors.GetInstance(types.TagsUser).Visit(starViewUrl.(string))
		}

		// 下一页
		url := common.GetNextPageUrl(resp.Request, repoLen)
		if url == "" {
			return
		}
		_ = resp.Request.Visit(url)
	})

	collector.OnError(func(resp *colly.Response, err error) {
		fmt.Println("rbu_errors:::", resp.StatusCode, resp.Body, err)
	})
}
