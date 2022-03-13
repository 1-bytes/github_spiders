package callbacks

import (
	"github.com/gocolly/colly/v2"
	"github_spiders/pkg/collectors"
	"github_spiders/pkg/utils"
	"github_spiders/spiders/github_com"
	"github_spiders/spiders/github_com/common"
	"github_spiders/spiders/github_com/user"
	"github_spiders/spiders/types"
	"log"
)

// UsersByRepo 列出已为存储库加注星标的人员.
// GitHub API docs url:
// https://docs.github.com/cn/rest/reference/activity#list-stargazers
type UsersByRepo struct {
}

// Callbacks 爬虫回调函数.
func (ur *UsersByRepo) Callbacks() {
	// 初始化变量
	auth := user.NewAuth()
	collector := collectors.GetInstance(types.TagsUser)

	// 对要提交的请求进行处理
	collector.OnRequest(func(r *colly.Request) {
		// GitHub's docs:
		// By default, all requests to https://api.github.com receive the v3 version of the REST API.
		// We encourage you to explicitly request this version via the Accept header.
		r.Headers.Add("Accept", "application/vnd.github.v3+json")
		r.Headers = auth.AddToken(r.Headers)
	})

	//  返回数据处理
	collector.OnResponse(func(resp *colly.Response) {
		users, err := utils.JsonUnmarshalBody(resp.Body)
		if err != nil {
			log.Printf("err: Failed to unmarshal the json: %s", err)
			return
		}

		// 检查该用户有没有 star 的存储库，或者是 star 的存储库列表到了最终页
		userLen := len(users)
		if userLen == 0 {
			return
		}

		var userName, starViewUrl string
		for _, u := range users {
			// 有关用户的一些信息，想要什么值可以自己取
			userName = u["login"].(string)
			starViewUrl = u["url"].(string) + "/starred"
			log.Printf("【New User】 Name:%s, URL:%s", userName, starViewUrl)
			starViewUrl = common.CheckUrl(starViewUrl)
			_ = collectors.GetInstance(types.TagsRepo).Visit(starViewUrl)
		}

		// 下一页
		url := common.GetNextPageUrl(resp.Request, userLen)
		if url == "" {
			return
		}
		_ = resp.Request.Visit(url)
	})

	// 错误处理
	collector.OnError(func(resp *colly.Response, err error) {
		github_com.ErrLock.Lock()
		defer github_com.ErrLock.Unlock()
		validity, t := auth.CheckTokenValidity(resp)
		if !validity {
			log.Printf("Toktn or IP is temporarily blocked, "+
				"unblock time is: %s Trying to change Token", t)
			auth.NextToken()
			auth.DelToken(resp.Request.Headers)
			_ = resp.Request.Retry()
		}
	})
}
